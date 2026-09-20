package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const apiBaseURL = "https://api.telegram.org"

// Bot es el cliente HTTP de la Telegram Bot API, sin dependencias
// externas. Long-polling, no webhook (ver .claude/rules/telegram-rules.md).
type Bot struct {
	token      string
	httpClient *http.Client
}

func NewBot(token string) *Bot {
	return &Bot{
		token: token,
		// Mayor al timeout de long-poll (30s) que pedimos en GetUpdates,
		// si no el cliente corta la conexión antes de que Telegram responda.
		httpClient: &http.Client{Timeout: 40 * time.Second},
	}
}

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ID int64 `json:"id"`
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

// GetUpdates hace long-polling nativo de Telegram (timeout=30s del lado
// del servidor, no sondeo agresivo del cliente).
func (b *Bot) GetUpdates(ctx context.Context, offset int) ([]Update, error) {
	params := url.Values{}
	params.Set("offset", strconv.Itoa(offset))
	params.Set("timeout", "30")

	var updates []Update
	if err := b.call(ctx, "getUpdates", params, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}

// SendMessage manda texto plano a un chat.
func (b *Bot) SendMessage(ctx context.Context, chatID int64, text string) error {
	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(chatID, 10))
	params.Set("text", text)
	return b.call(ctx, "sendMessage", params, nil)
}

func (b *Bot) call(ctx context.Context, method string, params url.Values, out any) error {
	endpoint := fmt.Sprintf("%s/bot%s/%s?%s", apiBaseURL, b.token, method, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("telegram: build request: %w", err)
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("telegram: read response: %w", err)
	}

	var parsed apiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("telegram: decode response: %w", err)
	}
	if !parsed.OK {
		return fmt.Errorf("telegram: api error: %s", parsed.Description)
	}

	if out != nil {
		if err := json.Unmarshal(parsed.Result, out); err != nil {
			return fmt.Errorf("telegram: decode result: %w", err)
		}
	}
	return nil
}
