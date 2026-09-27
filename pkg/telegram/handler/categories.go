package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/keyboard"
	"finance-tracker/pkg/telegram/model/update"
)

const (
	catCallbackPrefix = "cat:"
	catCallbackList   = "cat:listar"
	catCallbackAdd    = "cat:agregar"
)

// handleCategoriasCommand manda los botones Listar/Agregar de /categorias.
func (h *Handler) handleCategoriasCommand(ctx context.Context, msg *update.Message) {
	buttons := []keyboard.InlineButton{
		{Text: "Listar categorías", CallbackData: catCallbackList},
		{Text: "Agregar categoría", CallbackData: catCallbackAdd},
	}
	if _, err := h.Messenger.SendKeyboard(ctx, msg.Chat.ID, "¿Qué querés hacer?", buttons); err != nil {
		log.Printf("error mandando botones de categorías a chat %d: %v", msg.Chat.ID, err)
	}
}

// handleCategoriaCallback procesa la respuesta a los botones de
// /categorias: "listar" edita el mensaje con las categorías reales,
// "agregar" edita el mensaje pidiendo el nombre y deja al chat esperándolo
// (lo que hubiera pendiente se descarta).
func (h *Handler) handleCategoriaCallback(ctx context.Context, cq *update.CallbackQuery) {
	chatID := cq.Message.Chat.ID
	messageID := cq.Message.MessageID

	switch cq.Data {
	case catCallbackList:
		names, err := h.Categories.ListNames()
		if err != nil {
			log.Printf("error listando categorías para chat %d: %v", chatID, err)
			editMsg(ctx, h.Messenger, chatID, messageID, "No se pudo consultar las categorías. Intenta de nuevo.")
			return
		}
		text := "Sin categorías todavía."
		if len(names) > 0 {
			text = "Categorías:\n" + strings.Join(names, "\n")
		}
		editMsg(ctx, h.Messenger, chatID, messageID, text)

	case catCallbackAdd:
		h.setState(chatID, chatstate.NewCategoryPending(h.now()))
		editMsg(ctx, h.Messenger, chatID, messageID, "¿Cómo se llama la nueva categoría?")

	default:
		log.Printf("callback_data no reconocido: %q", cq.Data)
	}
}

// handleCategoryNameReply toma el texto como nombre de la categoría nueva
// y la crea en Notion. Sea cual sea el resultado, el chat deja de esperar
// el nombre.
func (h *Handler) handleCategoryNameReply(ctx context.Context, msg *update.Message) chatstate.ChatState {
	chatID := msg.Chat.ID

	name := strings.TrimSpace(msg.Text)
	if name == "" {
		send(ctx, h.Messenger, chatID, "Nombre vacío, no se creó ninguna categoría.")
		return chatstate.ChatState{}
	}

	if _, err := h.Categories.Create(name); err != nil {
		log.Printf("error creando categoría %q para chat %d: %v", name, chatID, err)
		send(ctx, h.Messenger, chatID, "No se pudo crear la categoría. Intenta de nuevo.")
		return chatstate.ChatState{}
	}

	send(ctx, h.Messenger, chatID, fmt.Sprintf("Categoría %q creada.", name))
	return chatstate.ChatState{}
}
