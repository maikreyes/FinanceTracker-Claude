package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"finance-tracker/pkg/config"
	"finance-tracker/pkg/ports"
)

const defaultAPIBaseURL = "https://api.telegram.org"

var _ ports.Messenger = (*Services)(nil)

// Services es el cliente HTTP de la Telegram Bot API, sin dependencias
// externas.
type Services struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

func NewServices(cfg *config.Config) *Services {
	return newServices(cfg.TelegramToken, defaultAPIBaseURL)
}

func newServices(token, baseURL string) *Services {
	return &Services{
		token:   token,
		baseURL: baseURL,
		// Mayor al timeout de long-poll (30s) que pedimos en GetUpdates,
		// si no el cliente corta la conexión antes de que Telegram responda.
		httpClient: &http.Client{Timeout: 40 * time.Second},
	}
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

func (s *Services) call(ctx context.Context, method string, params url.Values, out any) error {
	endpoint := fmt.Sprintf("%s/bot%s/%s?%s", s.baseURL, s.token, method, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("telegram: build request: %w", withoutURL(err))
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", withoutURL(err))
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

// withoutURL quita la URL del *url.Error que devuelve net/http: la API de
// Telegram lleva el token del bot en el path, y el error se loguea tal cual.
func withoutURL(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return urlErr.Err
	}
	return err
}
