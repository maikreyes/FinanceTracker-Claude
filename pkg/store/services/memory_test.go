package services

import (
	"testing"
	"time"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/transaction/model/transaction"
)

func TestMemory_GetSetDelete(t *testing.T) {
	s := NewMemory()

	got, err := s.Get(1)
	if err != nil {
		t.Fatalf("Get() error inesperado: %v", err)
	}
	if got != (chatstate.ChatState{}) {
		t.Errorf("Get(chat sin estado) = %+v, want ChatState{}", got)
	}

	if err := s.Set(1, chatstate.ChatState{Kind: chatstate.PendingPhotoType, PhotoFileID: "file-1"}); err != nil {
		t.Fatalf("Set() error inesperado: %v", err)
	}
	got, _ = s.Get(1)
	if got.PhotoFileID != "file-1" {
		t.Errorf("Get() = %+v, want PhotoFileID file-1", got)
	}

	if err := s.Delete(1); err != nil {
		t.Fatalf("Delete() error inesperado: %v", err)
	}
	got, _ = s.Get(1)
	if got != (chatstate.ChatState{}) {
		t.Errorf("Get() después de Delete = %+v, want ChatState{}", got)
	}
}

func TestMemory_TakeConsumeElEstado(t *testing.T) {
	s := NewMemory()
	s.Set(1, chatstate.ChatState{Kind: chatstate.PendingPhotoType, PhotoFileID: "file-1"})

	got, err := s.Take(1)
	if err != nil {
		t.Fatalf("Take() error inesperado: %v", err)
	}
	if got.PhotoFileID != "file-1" {
		t.Errorf("Take() = %+v, want PhotoFileID file-1", got)
	}

	again, _ := s.Take(1)
	if again != (chatstate.ChatState{}) {
		t.Errorf("segundo Take() = %+v, want ChatState{}", again)
	}
}

func TestMemory_GuardaCopias(t *testing.T) {
	s := NewMemory()
	w := &chatstate.Wizard{TxType: transaction.TypeExpense, Amount: 100}
	s.Set(1, chatstate.ChatState{Kind: chatstate.PendingWizard, Wizard: w})

	w.Amount = 999 // modificar lo que se pasó a Set no debe alterar lo guardado
	got, _ := s.Get(1)
	if got.Wizard.Amount != 100 {
		t.Errorf("Amount guardado = %v, want 100", got.Wizard.Amount)
	}

	got.Wizard.Amount = 555 // ni modificar lo que devuelve Get
	again, _ := s.Get(1)
	if again.Wizard.Amount != 100 {
		t.Errorf("Amount guardado = %v, want 100", again.Wizard.Amount)
	}
}

func TestMemory_MarkUpdateSeen(t *testing.T) {
	m := NewMemory()
	now := time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC)
	m.now = func() time.Time { return now }

	first, _ := m.MarkUpdateSeen(10)
	again, _ := m.MarkUpdateSeen(10)
	other, _ := m.MarkUpdateSeen(11)
	if !first || again || !other {
		t.Errorf("first=%v again=%v other=%v, want true false true", first, again, other)
	}

	now = now.Add(seenUpdateTTL + time.Minute)
	if expired, _ := m.MarkUpdateSeen(10); !expired {
		t.Error("pasado el TTL el update_id debería poder verse de nuevo")
	}
}
