package handler

import (
	"context"
	"log"
	"strings"

	"finance-tracker/pkg/telegram/model/update"
	"finance-tracker/pkg/transaction/model/transaction"
)

// normalizeCommand deja un comando en su forma comparable: sin espacios,
// sin el sufijo "@BotUsername" que Telegram agrega en algunos contextos, y
// en minúsculas.
func normalizeCommand(text string) string {
	text = strings.TrimSpace(text)
	if at := strings.Index(text, "@"); at != -1 {
		text = text[:at]
	}
	return strings.ToLower(text)
}

func matchesCommand(text, command string) bool {
	return normalizeCommand(text) == command
}

// parseStartCommand reconoce "/egreso" o "/ingreso".
func parseStartCommand(text string) (transaction.Type, bool) {
	switch normalizeCommand(text) {
	case "/egreso":
		return transaction.TypeExpense, true
	case "/ingreso":
		return transaction.TypeIncome, true
	default:
		return "", false
	}
}

// handleMessage registra el movimiento y responde con la confirmación de
// lo guardado o con un mensaje de error específico y accionable por tipo de
// falla — nunca el error crudo de Notion.
func (h *Handler) handleMessage(ctx context.Context, msg *update.Message) {
	result, err := h.Transactions.Register(msg.Text)
	if err != nil {
		log.Printf("error registrando mensaje de chat %d: %v", msg.Chat.ID, err)
		send(ctx, h.Messenger, msg.Chat.ID, errorReply(err))
		return
	}
	send(ctx, h.Messenger, msg.Chat.ID, successReply(result))
}

// handleResumen responde con el balance del mes calendario actual, los
// totales de Egreso/Ingreso, y cuántos movimientos son de cada tipo.
func (h *Handler) handleResumen(ctx context.Context, msg *update.Message) {
	total, err := h.Summarizer.SumThisMonth()
	if err != nil {
		log.Printf("error calculando resumen para chat %d: %v", msg.Chat.ID, err)
		send(ctx, h.Messenger, msg.Chat.ID, "No se pudo calcular el resumen. Intenta de nuevo en un momento.")
		return
	}
	send(ctx, h.Messenger, msg.Chat.ID, resumenReply(total))
}

func (h *Handler) handleAyuda(ctx context.Context, msg *update.Message) {
	send(ctx, h.Messenger, msg.Chat.ID, ayudaText)
}

// handleUltimo responde con el movimiento registrado más recientemente.
func (h *Handler) handleUltimo(ctx context.Context, msg *update.Message) {
	lt, found, err := h.Finder.FindLatest()
	if err != nil {
		log.Printf("error buscando el último movimiento para chat %d: %v", msg.Chat.ID, err)
		send(ctx, h.Messenger, msg.Chat.ID, "No se pudo consultar el último movimiento. Intenta de nuevo en un momento.")
		return
	}
	if !found {
		send(ctx, h.Messenger, msg.Chat.ID, "Todavía no hay movimientos registrados.")
		return
	}
	send(ctx, h.Messenger, msg.Chat.ID, ultimoReply(lt))
}
