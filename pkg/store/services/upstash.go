package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"finance-tracker/pkg/config"
	"finance-tracker/pkg/ports"
	"finance-tracker/pkg/telegram/model/chatstate"
)

var _ ports.ChatStore = (*Upstash)(nil)

// Upstash guarda el estado en Redis a través de la REST API de Upstash: un
// POST por comando, sin SDK ni conexión persistente, que es lo que cabe en
// una función serverless.
type Upstash struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewUpstash(cfg *config.Config) *Upstash {
	return newUpstash(cfg.UpstashURL, cfg.UpstashToken)
}

func newUpstash(baseURL, token string) *Upstash {
	return &Upstash{baseURL: baseURL, token: token, httpClient: &http.Client{Timeout: 5 * time.Second}}
}

func chatKey(chatID int64) string { return "finance-tracker:chat:" + strconv.FormatInt(chatID, 10) }
func updateKey(updateID int) string {
	return "finance-tracker:update:" + strconv.Itoa(updateID)
}

func (u *Upstash) Get(chatID int64) (chatstate.ChatState, error) {
	raw, err := u.command("GET", chatKey(chatID))
	if err != nil {
		return chatstate.ChatState{}, err
	}
	return decodeState(raw)
}

// Set guarda el estado con un TTL igual al tiempo que le queda antes de
// vencer, para que una conversación abandonada no quede colgada en Redis.
func (u *Upstash) Set(chatID int64, state chatstate.ChatState) error {
	body, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("upstash: marshal chat state: %w", err)
	}
	ttl := chatstate.PendingTTL
	if !state.ExpiresAt.IsZero() {
		ttl = time.Until(state.ExpiresAt)
	}
	seconds := int(math.Ceil(ttl.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	_, err = u.command("SET", chatKey(chatID), string(body), "EX", seconds)
	return err
}

func (u *Upstash) Delete(chatID int64) error {
	_, err := u.command("DEL", chatKey(chatID))
	return err
}

// Take usa GETDEL, que lee y borra en un único comando: dos updates
// concurrentes no pueden consumir el mismo estado.
func (u *Upstash) Take(chatID int64) (chatstate.ChatState, error) {
	raw, err := u.command("GETDEL", chatKey(chatID))
	if err != nil {
		return chatstate.ChatState{}, err
	}
	return decodeState(raw)
}

// MarkUpdateSeen usa SET NX: Redis responde OK solo a quien crea la clave.
func (u *Upstash) MarkUpdateSeen(updateID int) (bool, error) {
	raw, err := u.command("SET", updateKey(updateID), "1", "NX", "EX", int(seenUpdateTTL.Seconds()))
	if err != nil {
		return false, err
	}
	return string(raw) != "null", nil
}

func decodeState(raw json.RawMessage) (chatstate.ChatState, error) {
	var text *string
	if err := json.Unmarshal(raw, &text); err != nil {
		return chatstate.ChatState{}, fmt.Errorf("upstash: decode result: %w", err)
	}
	if text == nil {
		return chatstate.ChatState{}, nil
	}
	var state chatstate.ChatState
	if err := json.Unmarshal([]byte(*text), &state); err != nil {
		return chatstate.ChatState{}, fmt.Errorf("upstash: decode chat state: %w", err)
	}
	return state, nil
}

type commandResponse struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

func (u *Upstash) command(args ...any) (json.RawMessage, error) {
	body, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("upstash: marshal command: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, u.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("upstash: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+u.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstash: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("upstash: read response: %w", err)
	}

	var parsed commandResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, fmt.Errorf("upstash: decode response (status %d): %w", resp.StatusCode, err)
	}
	if parsed.Error != "" {
		return nil, fmt.Errorf("upstash: %s", parsed.Error)
	}
	return parsed.Result, nil
}
