package services

import (
	"net/http"
	"time"

	"finance-tracker/pkg/config"
	"finance-tracker/pkg/ports"
)

// apiBaseURL es var, no const, para poder apuntarla a un httptest.Server
// en los tests de este paquete.
var apiBaseURL = "https://api.groq.com/openai/v1"

// defaultModel: los modelos de Groq con soporte de visión cambian o se
// retiran seguido (confirmado en vivo el 2026-09-20: el modelo original,
// "llama-3.2-11b-vision-preview", ya estaba dado de baja). Si deja de
// existir, `GET /v1/models` filtrando por `input_modalities` con "image"
// da la lista vigente (ver .claude/rules/groq-rules.md).
const defaultModel = "qwen/qwen3.8-27b"

// Services interpreta fotos de recibos con Groq (API compatible con
// OpenAI), sin SDK externo.
type Services struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewServices(cfg *config.Config) *Services {
	return &Services{
		apiKey:     cfg.GroqAPIKey,
		model:      defaultModel,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

var _ ports.ReceiptInterpreter = (*Services)(nil)
