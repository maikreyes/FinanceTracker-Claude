package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"finance-tracker/pkg/telegram/model/keyboard"
)

// SendMessage manda texto plano a un chat.
func (s *Services) SendMessage(ctx context.Context, chatID int64, text string) error {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("text", text)
	return s.call(ctx, "sendMessage", params, nil)
}

// SendKeyboard manda un mensaje con botones inline en una sola fila y
// devuelve el message_id.
func (s *Services) SendKeyboard(ctx context.Context, chatID int64, text string, buttons []keyboard.InlineButton) (int, error) {
	row := make([]map[string]any, 0, len(buttons))
	for _, btn := range buttons {
		row = append(row, map[string]any{"text": btn.Text, "callback_data": btn.CallbackData})
	}
	markup, err := json.Marshal(map[string]any{"inline_keyboard": [][]map[string]any{row}})
	if err != nil {
		return 0, fmt.Errorf("telegram: marshal reply_markup: %w", err)
	}

	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("text", text)
	params.Set("reply_markup", string(markup))

	var sent struct {
		MessageID int `json:"message_id"`
	}
	if err := s.call(ctx, "sendMessage", params, &sent); err != nil {
		return 0, err
	}
	return sent.MessageID, nil
}
