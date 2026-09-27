package handler

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/keyboard"
	"finance-tracker/pkg/telegram/model/update"
	"finance-tracker/pkg/transaction/model/receipt"
	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
	"finance-tracker/pkg/transaction/services"
)

// fakeMessenger registra todo lo que el dispatcher le manda al canal.
type fakeMessenger struct {
	texts       []string // SendMessage y SendKeyboard, en orden
	keyboards   [][]keyboard.InlineButton
	edits       []string
	answered    []string
	downloaded  []string
	nextMsgID   int
	downloadErr error
	keyboardErr error
}

func (m *fakeMessenger) SendMessage(_ context.Context, _ int64, text string) error {
	m.texts = append(m.texts, text)
	return nil
}

func (m *fakeMessenger) SendKeyboard(_ context.Context, _ int64, text string, buttons []keyboard.InlineButton) (int, error) {
	if m.keyboardErr != nil {
		return 0, m.keyboardErr
	}
	m.texts = append(m.texts, text)
	m.keyboards = append(m.keyboards, buttons)
	m.nextMsgID++
	return m.nextMsgID, nil
}

func (m *fakeMessenger) EditMessageText(_ context.Context, _ int64, _ int, text string) error {
	m.edits = append(m.edits, text)
	return nil
}

func (m *fakeMessenger) AnswerCallbackQuery(_ context.Context, id string) error {
	m.answered = append(m.answered, id)
	return nil
}

func (m *fakeMessenger) DownloadPhoto(_ context.Context, fileID string) ([]byte, error) {
	m.downloaded = append(m.downloaded, fileID)
	if m.downloadErr != nil {
		return nil, m.downloadErr
	}
	return []byte("fake-image-bytes"), nil
}

func (m *fakeMessenger) lastText() string {
	if len(m.texts) == 0 {
		return ""
	}
	return m.texts[len(m.texts)-1]
}

func (m *fakeMessenger) lastEdit() string {
	if len(m.edits) == 0 {
		return ""
	}
	return m.edits[len(m.edits)-1]
}

// jsonStore guarda el estado serializado a JSON, como lo haría un store
// externo (Redis): un Set olvidado o un puntero compartido que en memoria
// pasaría desapercibido acá se nota.
type jsonStore struct {
	mu     sync.Mutex
	data   map[int64][]byte
	sets   int
	getErr error
	seen   map[int]bool
}

func newJSONStore() *jsonStore {
	return &jsonStore{data: make(map[int64][]byte), seen: make(map[int]bool)}
}

func (s *jsonStore) Get(chatID int64) (chatstate.ChatState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.getErr != nil {
		return chatstate.ChatState{}, s.getErr
	}
	return s.decode(chatID), nil
}

func (s *jsonStore) Set(chatID int64, state chatstate.ChatState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, err := json.Marshal(state)
	if err != nil {
		return err
	}
	s.data[chatID] = raw
	s.sets++
	return nil
}

func (s *jsonStore) Delete(chatID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, chatID)
	return nil
}

func (s *jsonStore) Take(chatID int64) (chatstate.ChatState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state := s.decode(chatID)
	delete(s.data, chatID)
	return state, nil
}

func (s *jsonStore) MarkUpdateSeen(updateID int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen[updateID] {
		return false, nil
	}
	s.seen[updateID] = true
	return true, nil
}

func (s *jsonStore) decode(chatID int64) chatstate.ChatState {
	var state chatstate.ChatState
	if raw, ok := s.data[chatID]; ok {
		if err := json.Unmarshal(raw, &state); err != nil {
			panic(err)
		}
	}
	return state
}

type fakeCategories struct {
	ids     map[string]string
	names   []string
	created []string
}

func (f *fakeCategories) Resolve(name string) (string, bool, error) {
	id, ok := f.ids[name]
	return id, ok, nil
}
func (f *fakeCategories) ListNames() ([]string, error) { return f.names, nil }
func (f *fakeCategories) Create(name string) (string, error) {
	f.created = append(f.created, name)
	return "new-cat-id", nil
}

type fakeWriter struct {
	got     transaction.Transaction
	created int
	id      string
	err     error
}

func (f *fakeWriter) CreateExpense(tx transaction.Transaction) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.got = tx
	f.created++
	return f.id, nil
}

type fakeUploader struct{}

func (fakeUploader) Upload([]byte, string, string) (string, error) { return "upload-id", nil }

type fakeInterpreter struct {
	result receipt.Interpreted
	err    error
}

func (f *fakeInterpreter) Interpret([]byte, []string) (receipt.Interpreted, error) {
	return f.result, f.err
}

type fakeSummarizer struct {
	summary summary.Monthly
	panics  bool
}

func (f fakeSummarizer) SumThisMonth() (summary.Monthly, error) {
	if f.panics {
		panic("boom")
	}
	return f.summary, nil
}

type fakeFinder struct{}

func (fakeFinder) FindLatest() (summary.Latest, bool, error) {
	return summary.Latest{}, false, nil
}

// harness junta las fakes para que cada test pueda mirar lo que necesite.
type harness struct {
	handler     *Handler
	messenger   *fakeMessenger
	store       *jsonStore
	categories  *fakeCategories
	writer      *fakeWriter
	interpreter *fakeInterpreter
	clock       time.Time
}

const (
	testChat   = int64(1)
	testSecret = "s3cret"
)

func newHarness(t *testing.T) *harness {
	t.Helper()
	h := &harness{
		messenger:   &fakeMessenger{},
		store:       newJSONStore(),
		categories:  &fakeCategories{ids: map[string]string{"comida": "cat-123"}, names: []string{"comida"}},
		writer:      &fakeWriter{id: "page-1"},
		interpreter: &fakeInterpreter{},
		clock:       time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC),
	}
	h.handler = &Handler{
		Messenger:      h.messenger,
		Transactions:   services.NewServices(h.categories, h.writer, fakeUploader{}, h.interpreter),
		Summarizer:     fakeSummarizer{},
		Finder:         fakeFinder{},
		Categories:     h.categories,
		Store:          h.store,
		AllowedChatIDs: []int64{testChat},
		WebhookSecret:  testSecret,
		Now:            func() time.Time { return h.clock },
	}
	return h
}

func (h *harness) say(text string) {
	h.handler.Dispatch(context.Background(), update.Update{Message: &update.Message{Chat: update.Chat{ID: testChat}, Text: text}})
}

func (h *harness) sendPhoto(fileID, caption string) {
	h.handler.Dispatch(context.Background(), update.Update{Message: &update.Message{
		Chat:    update.Chat{ID: testChat},
		Caption: caption,
		Photo:   []update.PhotoSize{{FileID: "small-" + fileID}, {FileID: fileID}},
	}})
}

func (h *harness) tap(messageID int, data string) {
	h.handler.Dispatch(context.Background(), update.Update{CallbackQuery: &update.CallbackQuery{
		ID:      "cq",
		Data:    data,
		Message: &update.Message{Chat: update.Chat{ID: testChat}, MessageID: messageID},
	}})
}

func (h *harness) state(t *testing.T) chatstate.ChatState {
	t.Helper()
	state, err := h.store.Get(testChat)
	if err != nil {
		t.Fatalf("Get() error inesperado: %v", err)
	}
	return state
}

var errBoom = errors.New("boom")
