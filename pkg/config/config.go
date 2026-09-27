package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

// Config agrupa las variables de entorno del bot. NewConfig no falla por
// una variable ausente: cada punto de entrada valida solo las que necesita
// (CheckTelegram, CheckWebhook).
type Config struct {
	TelegramToken  string
	AllowedChatIDs []int64
	WebhookSecret  string

	NotionAPIKey               string
	NotionDatabaseID           string // data source "Expenses"
	NotionCategoryDataSourceID string // data source "Category"

	GroqAPIKey string

	UpstashURL   string
	UpstashToken string
}

func NewConfig() *Config {
	return &Config{
		TelegramToken:              os.Getenv("TELEGRAM_BOT_TOKEN"),
		AllowedChatIDs:             ParseAllowedChatIDs(os.Getenv("TELEGRAM_ALLOWED_CHAT_IDS")),
		WebhookSecret:              os.Getenv("TELEGRAM_WEBHOOK_SECRET"),
		NotionAPIKey:               os.Getenv("NOTION_API_KEY"),
		NotionDatabaseID:           os.Getenv("NOTION_DATABASE_ID"),
		NotionCategoryDataSourceID: os.Getenv("NOTION_CATEGORY_DATA_SOURCE_ID"),
		GroqAPIKey:                 os.Getenv("GROQ_API_KEY"),
		UpstashURL:                 os.Getenv("UPSTASH_REDIS_REST_URL"),
		UpstashToken:               os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
	}
}

// CheckTelegram valida lo mínimo para hablar con Telegram.
func (c *Config) CheckTelegram() error {
	if c.TelegramToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN no está seteado")
	}
	return nil
}

// CheckWebhook valida lo que exige el modo webhook: sin secreto cualquiera
// que adivine la URL podría mandar updates falsos, y sin store externo un
// flujo pendiente se perdería entre invocaciones.
func (c *Config) CheckWebhook() error {
	if err := c.CheckTelegram(); err != nil {
		return err
	}
	switch {
	case c.WebhookSecret == "":
		return errors.New("TELEGRAM_WEBHOOK_SECRET no está seteado")
	case c.UpstashURL == "" || c.UpstashToken == "":
		return errors.New("UPSTASH_REDIS_REST_URL y UPSTASH_REDIS_REST_TOKEN no están seteados")
	}
	return nil
}

// ParseAllowedChatIDs parsea TELEGRAM_ALLOWED_CHAT_IDS (lista separada por
// comas). Un valor no numérico se ignora con una advertencia: no aborta el
// arranque por un typo en un solo ID de una lista con varios.
func ParseAllowedChatIDs(raw string) []int64 {
	var ids []int64
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			log.Printf("TELEGRAM_ALLOWED_CHAT_IDS: valor inválido ignorado: %q", part)
			continue
		}
		ids = append(ids, id)
	}
	return ids
}
