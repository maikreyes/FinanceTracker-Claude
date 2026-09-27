package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/update"
	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

// handleStartWizard arranca (o reinicia) el flujo guiado.
func (h *Handler) handleStartWizard(ctx context.Context, msg *update.Message, txType transaction.Type) chatstate.ChatState {
	send(ctx, h.Messenger, msg.Chat.ID, fmt.Sprintf("Registrando un %s.\n¿Cuál es el monto?", txType))
	return chatstate.NewWizardPending(chatstate.Wizard{TxType: txType, Step: chatstate.StepAmount}, h.now())
}

// handleWizardReply procesa la respuesta al paso actual del flujo guiado.
// Una respuesta inválida deja el flujo en el mismo paso sin perder lo ya
// contestado.
func (h *Handler) handleWizardReply(ctx context.Context, state chatstate.ChatState, msg *update.Message) chatstate.ChatState {
	w := *state.Wizard
	text := strings.TrimSpace(msg.Text)

	if w.Step == chatstate.StepCategory {
		return h.finishWizard(ctx, state, w, skipDash(text), msg.Chat.ID)
	}

	next, reply, advanced := advanceWizard(w, text)
	send(ctx, h.Messenger, msg.Chat.ID, reply)
	if !advanced {
		return state
	}
	return chatstate.NewWizardPending(next, h.now())
}

// advanceWizard aplica text al paso actual de w. advanced=false significa
// que la respuesta no sirvió y reply explica por qué.
func advanceWizard(w chatstate.Wizard, text string) (next chatstate.Wizard, reply string, advanced bool) {
	switch w.Step {
	case chatstate.StepAmount:
		amount, err := strconv.ParseFloat(text, 64)
		if err != nil || !transaction.ValidAmount(amount) {
			return w, "No entendí el monto. Mandá solo un número mayor a 0 (ej. 13000).", false
		}
		w.Amount = amount
		w.Step = chatstate.StepConcept
		return w, "¿Cuál es el concepto? (mandá \"-\" para omitirlo)", true

	case chatstate.StepConcept:
		w.Concept = skipDash(text)
		w.Step = chatstate.StepPaymentMethod
		return w, fmt.Sprintf("¿Cuál fue el medio de pago? Opciones: %s", strings.Join(transaction.ValidPaymentMethods(), ", ")), true

	case chatstate.StepPaymentMethod:
		method, ok := transaction.ParsePaymentMethod(text)
		if !ok {
			return w, fmt.Sprintf("No reconocí el medio de pago. Opciones: %s", strings.Join(transaction.ValidPaymentMethods(), ", ")), false
		}
		w.Method = method
		w.Step = chatstate.StepCategory
		return w, "¿Cuál es la categoría? (mandá \"-\" para omitirla)", true
	}
	return w, "", false
}

// finishWizard registra con lo recolectado. Solo una categoría no
// reconocida permite reintentar el último paso; cualquier otro resultado
// cierra el flujo.
func (h *Handler) finishWizard(ctx context.Context, state chatstate.ChatState, w chatstate.Wizard, categoryName string, chatID int64) chatstate.ChatState {
	result, err := h.Transactions.RegisterFromFields(register.ConversationInput{
		Type:          w.TxType,
		Concept:       w.Concept,
		Amount:        w.Amount,
		PaymentMethod: w.Method,
		CategoryName:  categoryName,
	})
	if err != nil {
		var errCategory failure.ErrUnrecognizedCategory
		if errors.As(err, &errCategory) {
			send(ctx, h.Messenger, chatID, errorReply(err)+"\nMandá \"-\" para omitirla.")
			return state
		}
		log.Printf("error registrando conversación de chat %d: %v", chatID, err)
		send(ctx, h.Messenger, chatID, errorReply(err))
		return chatstate.ChatState{}
	}

	send(ctx, h.Messenger, chatID, successReply(result))
	return chatstate.ChatState{}
}

// skipDash devuelve "" para la respuesta "-" (omitir un campo opcional).
func skipDash(text string) string {
	if text == "-" {
		return ""
	}
	return text
}
