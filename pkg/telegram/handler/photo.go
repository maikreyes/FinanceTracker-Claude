package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/keyboard"
	"finance-tracker/pkg/telegram/model/update"
	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

const (
	tipoCallbackPrefix = "tipo:"
	staleQuestionText  = "Esta pregunta ya no aplica. Mandá la foto de nuevo."
)

func tipoCallbackData(t transaction.Type) string { return tipoCallbackPrefix + string(t) }

func parseTipoCallback(data string) (transaction.Type, bool) {
	switch data {
	case tipoCallbackData(transaction.TypeIncome):
		return transaction.TypeIncome, true
	case tipoCallbackData(transaction.TypeExpense):
		return transaction.TypeExpense, true
	default:
		return "", false
	}
}

// handlePhotoMessage maneja un mensaje con foto. Con caption, el caption
// es el mismo formato de texto de siempre: se descarga la foto y se
// registra directo. Sin caption se pregunta Ingreso/Egreso con botones y
// se guarda solo el file_id — la descarga real y la interpretación por IA
// pasan cuando el usuario contesta (handleTipoCallback). Cualquier foto
// cancela lo que hubiera pendiente.
func (h *Handler) handlePhotoMessage(ctx context.Context, msg *update.Message) chatstate.ChatState {
	chatID := msg.Chat.ID
	caption := strings.TrimSpace(msg.Caption)
	fileID := msg.Photo[len(msg.Photo)-1].FileID // la última es la de mayor calidad

	if caption == "" {
		buttons := []keyboard.InlineButton{
			{Text: "Ingreso", CallbackData: tipoCallbackData(transaction.TypeIncome)},
			{Text: "Egreso", CallbackData: tipoCallbackData(transaction.TypeExpense)},
		}
		promptID, err := h.Messenger.SendKeyboard(ctx, chatID, "¿Es un ingreso o un egreso?", buttons)
		if err != nil {
			log.Printf("error mandando botones a chat %d: %v", chatID, err)
			return chatstate.ChatState{}
		}
		return chatstate.NewPhotoPending(fileID, promptID, h.now())
	}

	photo, err := h.Messenger.DownloadPhoto(ctx, fileID)
	if err != nil {
		log.Printf("error descargando foto para chat %d: %v", chatID, err)
		send(ctx, h.Messenger, chatID, "No se pudo descargar la foto. Intenta de nuevo.")
		return chatstate.ChatState{}
	}

	result, err := h.Transactions.RegisterWithPhoto(caption, photo, register.ReceiptFilename, register.ReceiptContentType)
	if err != nil {
		log.Printf("error registrando mensaje con foto de chat %d: %v", chatID, err)
		send(ctx, h.Messenger, chatID, errorReply(err))
		return chatstate.ChatState{}
	}
	send(ctx, h.Messenger, chatID, successReply(result))
	return chatstate.ChatState{}
}

// handleTipoCallback resuelve la pregunta Ingreso/Egreso de una foto: toma
// (de forma atómica) el file_id pendiente de ese chat, descarga la foto
// recién ahora, la interpreta con IA usando el tipo elegido, y registra.
//
// El estado se consume antes de registrar para que un doble toque o una
// reentrega de Telegram no cree dos páginas. Si el fallo es transitorio
// (Telegram, la IA o Notion), el estado se restaura: los botones siguen en
// pantalla y tocarlos de nuevo reintenta sin reenviar la foto.
func (h *Handler) handleTipoCallback(ctx context.Context, cq *update.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	promptID := cq.Message.MessageID

	txType, ok := parseTipoCallback(cq.Data)
	if !ok {
		log.Printf("callback_data no reconocido: %q", cq.Data)
		return
	}

	state, err := h.Store.Take(chatID)
	if err != nil {
		log.Printf("error tomando estado de chat %d: %v", chatID, err)
		send(ctx, h.Messenger, chatID, "Tuve un problema técnico. Intenta de nuevo en un momento.")
		return
	}

	live := state.Kind != chatstate.PendingNone && !state.Expired(h.now())
	if state.Kind != chatstate.PendingPhotoType || state.PromptMessageID != promptID || !live {
		// Botones de una foto anterior, o el chat espera otra cosa: lo que
		// esté pendiente no es de este toque y se devuelve intacto.
		if live {
			h.setState(chatID, state)
		}
		editMsg(ctx, h.Messenger, chatID, promptID, staleQuestionText)
		return
	}

	photo, err := h.Messenger.DownloadPhoto(ctx, state.PhotoFileID)
	if err != nil {
		log.Printf("error descargando foto para chat %d: %v", chatID, err)
		h.setState(chatID, state)
		send(ctx, h.Messenger, chatID, "No se pudo descargar la foto. Tocá el botón de nuevo para reintentar.")
		return
	}

	result, err := h.Transactions.RegisterFromPhotoAI(photo, register.ReceiptFilename, register.ReceiptContentType, txType)
	if err != nil {
		log.Printf("error registrando foto (botón) de chat %d: %v", chatID, err)
		var noAmount failure.ErrReceiptAmountNotFound
		if errors.As(err, &noAmount) {
			// Reintentar la misma foto daría el mismo resultado.
			editMsg(ctx, h.Messenger, chatID, promptID, fmt.Sprintf("Tipo elegido: %s", txType))
		} else {
			h.setState(chatID, state)
		}
		send(ctx, h.Messenger, chatID, errorReply(err))
		return
	}

	editMsg(ctx, h.Messenger, chatID, promptID, fmt.Sprintf("Tipo elegido: %s", txType))
	send(ctx, h.Messenger, chatID, successReply(result))
}
