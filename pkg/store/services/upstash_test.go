package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/transaction/model/transaction"
)

// fakeRedis atiende el subconjunto de comandos que usa Upstash (GET, SET
// con EX/NX, DEL, GETDEL), con la misma forma de respuesta que la REST
// API real.
type fakeRedis struct {
	data     map[string]string
	ttls     map[string]int
	commands [][]any
	failWith string
}

func newFakeRedis(t *testing.T) (*fakeRedis, *Upstash) {
	t.Helper()
	fr := &fakeRedis{data: map[string]string{}, ttls: map[string]int{}}
	server := httptest.NewServer(http.HandlerFunc(fr.serve))
	t.Cleanup(server.Close)
	return fr, newUpstash(server.URL, "tok")
}

func (fr *fakeRedis) serve(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer tok" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "unauthorized"})
		return
	}
	var args []any
	json.NewDecoder(r.Body).Decode(&args)
	fr.commands = append(fr.commands, args)

	if fr.failWith != "" {
		json.NewEncoder(w).Encode(map[string]any{"error": fr.failWith})
		return
	}

	key, _ := args[1].(string)
	var result any
	switch args[0] {
	case "GET":
		if v, ok := fr.data[key]; ok {
			result = v
		}
	case "GETDEL":
		if v, ok := fr.data[key]; ok {
			result = v
			delete(fr.data, key)
		}
	case "DEL":
		delete(fr.data, key)
		result = 1
	case "SET":
		nx := false
		ttl := 0
		for i, a := range args {
			switch a {
			case "NX":
				nx = true
			case "EX":
				ttl = int(args[i+1].(float64))
			}
		}
		if _, exists := fr.data[key]; nx && exists {
			break
		}
		fr.data[key] = args[2].(string)
		fr.ttls[key] = ttl
		result = "OK"
	}
	json.NewEncoder(w).Encode(map[string]any{"result": result})
}

func TestUpstash_SetGetConTTL(t *testing.T) {
	fr, u := newFakeRedis(t)
	expires := time.Now().Add(10 * time.Minute)
	state := chatstate.NewWizardPending(chatstate.Wizard{TxType: transaction.TypeExpense, Amount: 500}, time.Now())
	state.ExpiresAt = expires

	if err := u.Set(7, state); err != nil {
		t.Fatalf("Set() error inesperado: %v", err)
	}

	got, err := u.Get(7)
	if err != nil {
		t.Fatalf("Get() error inesperado: %v", err)
	}
	if got.Kind != chatstate.PendingWizard || got.Wizard == nil || got.Wizard.Amount != 500 {
		t.Errorf("Get() = %+v, want el wizard guardado", got)
	}
	if ttl := fr.ttls[chatKey(7)]; ttl < 590 || ttl > 600 {
		t.Errorf("TTL = %ds, want ~600 (lo que le queda al estado)", ttl)
	}
}

func TestUpstash_GetSinEstadoDevuelveCero(t *testing.T) {
	_, u := newFakeRedis(t)

	got, err := u.Get(99)
	if err != nil {
		t.Fatalf("Get() error inesperado: %v", err)
	}
	if got != (chatstate.ChatState{}) {
		t.Errorf("Get() = %+v, want ChatState{}", got)
	}
}

func TestUpstash_TakeConsumeElEstado(t *testing.T) {
	fr, u := newFakeRedis(t)
	u.Set(7, chatstate.NewPhotoPending("file-1", 3, time.Now()))

	got, err := u.Take(7)
	if err != nil {
		t.Fatalf("Take() error inesperado: %v", err)
	}
	if got.PhotoFileID != "file-1" {
		t.Errorf("Take() = %+v, want PhotoFileID file-1", got)
	}
	if again, _ := u.Take(7); again != (chatstate.ChatState{}) {
		t.Errorf("segundo Take() = %+v, want ChatState{}", again)
	}
	last := fr.commands[len(fr.commands)-1]
	if last[0] != "GETDEL" {
		t.Errorf("Take usó %v, want GETDEL (atómico)", last[0])
	}
}

func TestUpstash_Delete(t *testing.T) {
	_, u := newFakeRedis(t)
	u.Set(7, chatstate.NewCategoryPending(time.Now()))

	if err := u.Delete(7); err != nil {
		t.Fatalf("Delete() error inesperado: %v", err)
	}
	if got, _ := u.Get(7); got != (chatstate.ChatState{}) {
		t.Errorf("Get() después de Delete = %+v, want ChatState{}", got)
	}
}

func TestUpstash_MarkUpdateSeen(t *testing.T) {
	fr, u := newFakeRedis(t)

	first, err := u.MarkUpdateSeen(42)
	if err != nil {
		t.Fatalf("MarkUpdateSeen() error inesperado: %v", err)
	}
	again, _ := u.MarkUpdateSeen(42)
	if !first || again {
		t.Errorf("first=%v again=%v, want true false", first, again)
	}
	if ttl := fr.ttls[updateKey(42)]; ttl != int((24 * time.Hour).Seconds()) {
		t.Errorf("TTL = %s, want 24h", strconv.Itoa(ttl))
	}
}

func TestUpstash_ErrorDeLaAPI(t *testing.T) {
	fr, u := newFakeRedis(t)
	fr.failWith = "ERR boom"

	if _, err := u.Get(1); err == nil {
		t.Error("Get() con error de Redis debería fallar")
	}
	if err := u.Set(1, chatstate.NewCategoryPending(time.Now())); err == nil {
		t.Error("Set() con error de Redis debería fallar")
	}
}

func TestUpstash_TokenIncorrecto(t *testing.T) {
	fr, _ := newFakeRedis(t)
	server := httptest.NewServer(http.HandlerFunc(fr.serve))
	defer server.Close()

	if _, err := newUpstash(server.URL, "otro").Get(1); err == nil {
		t.Error("un token incorrecto debería fallar")
	}
}
