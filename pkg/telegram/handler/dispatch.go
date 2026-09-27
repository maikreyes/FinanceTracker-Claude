package handler

import (
	"context"
	"log"
	"reflect"
	"runtime/debug"
	"strings"

	"finance-tracker/pkg/ports"
	"finance-tracker/pkg/telegram/model/chatstate"
	"finance-tracker/pkg/telegram/model/update"
)

func isAllowed(chatID int64, allowed []int64) bool {
	for _, id := range allowed {
		if id == chatID {
			return true
		}
	}
	return false
}

// Dispatch procesa un único update de Telegram: autoriza el chat, decide
// qué handler corre y gestiona el estado pendiente del chat en Store. Un
// panic dentro de un handler se loguea y no tumba al proceso: Telegram
// reentregaría el mismo update y el bot entraría en un bucle de caídas.
func (h *Handler) Dispatch(ctx context.Context, u update.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic procesando update %d: %v\n%s", u.UpdateID, r, debug.Stack())
		}
	}()

	switch {
	case u.CallbackQuery != nil:
		cq := u.CallbackQuery
		if cq.Message == nil || !isAllowed(cq.Message.Chat.ID, h.AllowedChatIDs) {
			log.Printf("callback ignorado de chat no autorizado")
			return
		}
		h.handleCallbackQuery(ctx, cq)

	case u.Message != nil:
		msg := u.Message
		if !isAllowed(msg.Chat.ID, h.AllowedChatIDs) {
			log.Printf("mensaje ignorado de chat no autorizado: %d", msg.Chat.ID)
			return
		}
		h.dispatchMessage(ctx, msg)
	}
}

func (h *Handler) dispatchMessage(ctx context.Context, msg *update.Message) {
	// Stickers, audios, documentos sin foto: no traen nada que interpretar
	// y no deben contar como respuesta al paso de un flujo pendiente.
	if len(msg.Photo) == 0 && msg.Text == "" {
		return
	}

	chatID := msg.Chat.ID
	prev, ok := h.loadState(ctx, chatID)
	if !ok {
		return
	}
	h.saveState(chatID, prev, h.routeMessage(ctx, prev, msg))
}

// routeMessage elige el handler y devuelve el estado resultante.
//
//   - Una foto, /egreso e /ingreso reemplazan lo pendiente por lo suyo.
//   - Los demás comandos cancelan lo pendiente.
//   - Un mensaje de texto durante un flujo guiado o la espera de un nombre
//     de categoría es la respuesta a ese paso.
//   - Cualquier otro texto se registra como movimiento y deja intacto lo
//     pendiente (una foto esperando sus botones sigue esperándolos).
func (h *Handler) routeMessage(ctx context.Context, state chatstate.ChatState, msg *update.Message) chatstate.ChatState {
	startType, isStart := parseStartCommand(msg.Text)

	switch {
	case len(msg.Photo) > 0:
		return h.handlePhotoMessage(ctx, msg)

	case isStart:
		return h.handleStartWizard(ctx, msg, startType)

	case matchesCommand(msg.Text, "/ayuda"):
		h.handleAyuda(ctx, msg)
		return chatstate.ChatState{}

	case matchesCommand(msg.Text, "/ultimo"):
		h.handleUltimo(ctx, msg)
		return chatstate.ChatState{}

	case matchesCommand(msg.Text, "/categorias"):
		h.handleCategoriasCommand(ctx, msg)
		return chatstate.ChatState{}

	case matchesCommand(msg.Text, "/resumen"):
		h.handleResumen(ctx, msg)
		return chatstate.ChatState{}

	case state.Kind == chatstate.PendingCategoryName:
		return h.handleCategoryNameReply(ctx, msg)

	case state.Kind == chatstate.PendingWizard && state.Wizard != nil:
		return h.handleWizardReply(ctx, state, msg)

	default:
		h.handleMessage(ctx, msg)
		return state
	}
}

// handleCallbackQuery reconoce el toque (para que el botón deje de mostrar
// "cargando") y reparte por el prefijo de cq.Data.
func (h *Handler) handleCallbackQuery(ctx context.Context, cq *update.CallbackQuery) {
	if err := h.Messenger.AnswerCallbackQuery(ctx, cq.ID); err != nil {
		log.Printf("error respondiendo callback query: %v", err)
	}

	switch {
	case strings.HasPrefix(cq.Data, tipoCallbackPrefix):
		h.handleTipoCallback(ctx, cq)
	case strings.HasPrefix(cq.Data, catCallbackPrefix):
		h.handleCategoriaCallback(ctx, cq)
	default:
		log.Printf("callback_data no reconocido: %q", cq.Data)
	}
}

// loadState lee el estado del chat. Si la lectura falla devuelve ok=false
// y avisa al usuario: seguir con un estado vacío y luego guardarlo pisaría
// el estado real (un flujo guiado a medias) por un error transitorio del
// store.
func (h *Handler) loadState(ctx context.Context, chatID int64) (chatstate.ChatState, bool) {
	state, err := h.Store.Get(chatID)
	if err != nil {
		log.Printf("error leyendo estado de chat %d: %v", chatID, err)
		send(ctx, h.Messenger, chatID, "Tuve un problema técnico. Intenta de nuevo en un momento.")
		return chatstate.ChatState{}, false
	}
	if state.Expired(h.now()) {
		if err := h.Store.Delete(chatID); err != nil {
			log.Printf("error borrando estado vencido de chat %d: %v", chatID, err)
		}
		return chatstate.ChatState{}, true
	}
	return state, true
}

// saveState persiste next solo si cambió respecto de prev: un mensaje que
// no toca lo pendiente no cuesta una escritura al store.
func (h *Handler) saveState(chatID int64, prev, next chatstate.ChatState) {
	if reflect.DeepEqual(prev, next) {
		return
	}
	h.setState(chatID, next)
}

func (h *Handler) setState(chatID int64, state chatstate.ChatState) {
	var err error
	if state.Kind == chatstate.PendingNone {
		err = h.Store.Delete(chatID)
	} else {
		err = h.Store.Set(chatID, state)
	}
	if err != nil {
		log.Printf("error guardando estado de chat %d: %v", chatID, err)
	}
}

func send(ctx context.Context, m ports.Messenger, chatID int64, text string) {
	if err := m.SendMessage(ctx, chatID, text); err != nil {
		log.Printf("error respondiendo a chat %d: %v", chatID, err)
	}
}

func editMsg(ctx context.Context, m ports.Messenger, chatID int64, messageID int, text string) {
	if err := m.EditMessageText(ctx, chatID, messageID, text); err != nil {
		log.Printf("error editando mensaje %d de chat %d: %v", messageID, chatID, err)
	}
}
