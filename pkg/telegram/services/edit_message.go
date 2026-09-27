package services

import (
	"context"
	"net/url"
	"strconv"
)

// EditMessageText reemplaza el texto de un mensaje ya enviado.
func (s *Services) EditMessageText(ctx context.Context, chatID int64, messageID int, text string) error {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("message_id", strconv.Itoa(messageID))
	params.Set("text", text)
	return s.call(ctx, "editMessageText", params, nil)
}

// AnswerCallbackQuery le avisa a Telegram que ya se procesó el toque del
// botón: sin esto, el botón queda con el ícono de "cargando" en el
// cliente del usuario indefinidamente.
func (s *Services) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error {
	params := url.Values{}
	params.Set("callback_query_id", callbackQueryID)
	return s.call(ctx, "answerCallbackQuery", params, nil)
}
