package handler

import (
	"time"

	"finance-tracker/pkg/config"
	"finance-tracker/pkg/ports"
)

// Handler procesa los updates de Telegram. Cada punto de entrada
// (long-polling, webhook) arma uno con su propio ChatStore y llama a los
// mismos métodos.
type Handler struct {
	Messenger      ports.Messenger
	Transactions   ports.TransactionServices
	Summarizer     ports.MonthlySummarizer
	Finder         ports.LatestTransactionFinder
	Categories     ports.CategoryResolver
	Store          ports.ChatStore
	AllowedChatIDs []int64
	WebhookSecret  string

	// Now es el reloj de los vencimientos del estado pendiente; nil usa
	// time.Now.
	Now func() time.Time
}

func NewHandler(
	cfg *config.Config,
	messenger ports.Messenger,
	transactions ports.TransactionServices,
	summarizer ports.MonthlySummarizer,
	finder ports.LatestTransactionFinder,
	categories ports.CategoryResolver,
	store ports.ChatStore,
) *Handler {
	return &Handler{
		Messenger:      messenger,
		Transactions:   transactions,
		Summarizer:     summarizer,
		Finder:         finder,
		Categories:     categories,
		Store:          store,
		AllowedChatIDs: cfg.AllowedChatIDs,
		WebhookSecret:  cfg.WebhookSecret,
	}
}

func (h *Handler) now() time.Time {
	if h.Now != nil {
		return h.Now()
	}
	return time.Now()
}
