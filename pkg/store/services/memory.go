package services

import (
	"sync"
	"time"

	"finance-tracker/pkg/ports"
	"finance-tracker/pkg/telegram/model/chatstate"
)

const seenUpdateTTL = 24 * time.Hour

var _ ports.ChatStore = (*Memory)(nil)

// Memory guarda el estado en memoria de proceso, pensado para desarrollo
// local con long-polling. Guarda copias, como lo haría un store externo
// serializado: modificar un ChatState obtenido con Get no altera lo
// guardado.
type Memory struct {
	mu     sync.Mutex
	states map[int64]chatstate.ChatState
	seen   map[int]time.Time // update_id -> vencimiento
	now    func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		states: make(map[int64]chatstate.ChatState),
		seen:   make(map[int]time.Time),
		now:    time.Now,
	}
}

func (m *Memory) Get(chatID int64) (chatstate.ChatState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return clone(m.states[chatID]), nil
}

func (m *Memory) Set(chatID int64, state chatstate.ChatState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[chatID] = clone(state)
	return nil
}

func (m *Memory) Delete(chatID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, chatID)
	return nil
}

func (m *Memory) Take(chatID int64) (chatstate.ChatState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	state := m.states[chatID]
	delete(m.states, chatID)
	return state, nil
}

func (m *Memory) MarkUpdateSeen(updateID int) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	for id, expires := range m.seen {
		if now.After(expires) {
			delete(m.seen, id)
		}
	}
	if _, seen := m.seen[updateID]; seen {
		return false, nil
	}
	m.seen[updateID] = now.Add(seenUpdateTTL)
	return true, nil
}

func clone(state chatstate.ChatState) chatstate.ChatState {
	if state.Wizard != nil {
		w := *state.Wizard
		state.Wizard = &w
	}
	return state
}
