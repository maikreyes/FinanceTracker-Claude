package app

import (
	"testing"

	"finance-tracker/pkg/config"
	store "finance-tracker/pkg/store/services"
)

func TestNew_CablaTodosLosPuertos(t *testing.T) {
	a := New(&config.Config{TelegramToken: "t", AllowedChatIDs: []int64{1}}, store.NewMemory())

	h := a.Handler
	if h.Messenger == nil || h.Transactions == nil || h.Summarizer == nil || h.Finder == nil || h.Categories == nil || h.Store == nil {
		t.Errorf("handler con puertos sin cablear: %+v", h)
	}
	if len(h.AllowedChatIDs) != 1 || a.Telegram == nil {
		t.Error("New debería propagar la whitelist y exponer el cliente de Telegram")
	}
}

func TestCommands_IncluyenResumen(t *testing.T) {
	for _, c := range Commands {
		if c.Command == "resumen" {
			return
		}
	}
	t.Error("el menú de comandos debería incluir /resumen")
}
