// Comando de un solo uso para apuntar el webhook de Telegram a la URL
// desplegada en Vercel, o para borrarlo y volver a long-polling.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/joho/godotenv"

	"finance-tracker/pkg/config"
	telegram "finance-tracker/pkg/telegram/services"
)

func main() {
	webhookURL := flag.String("url", "", "URL pública HTTPS del webhook, ej. https://mi-proyecto.vercel.app/api/webhook")
	del := flag.Bool("delete", false, "borra el webhook en vez de configurarlo (para volver a long-polling)")
	dropPending := flag.Bool("drop-pending", false, "con -delete, descarta los updates que Telegram tenía encolados")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("no se cargó .env, se usa el entorno del proceso: %v", err)
	}

	cfg := config.NewConfig()
	if err := cfg.CheckTelegram(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tg := telegram.NewServices(cfg)

	if *del {
		if err := tg.DeleteWebhook(ctx, *dropPending); err != nil {
			log.Fatalf("no se pudo borrar el webhook: %v", err)
		}
		log.Println("webhook borrado: ya se puede usar long-polling (go run ./cmd)")
		return
	}

	if cfg.WebhookSecret == "" {
		log.Fatal("TELEGRAM_WEBHOOK_SECRET no está seteado: sin secreto cualquiera podría mandar updates falsos")
	}
	if err := validateWebhookURL(*webhookURL); err != nil {
		log.Fatal(err)
	}

	if err := tg.SetWebhook(ctx, *webhookURL, cfg.WebhookSecret); err != nil {
		log.Fatalf("no se pudo configurar el webhook: %v", err)
	}
	log.Printf("webhook configurado en %s (long-polling deja de funcionar hasta borrarlo con -delete)", *webhookURL)
}

var errInvalidURL = errors.New("-url es obligatorio y debe ser una URL https (Telegram solo entrega webhooks por HTTPS)")

func validateWebhookURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return errInvalidURL
	}
	return nil
}
