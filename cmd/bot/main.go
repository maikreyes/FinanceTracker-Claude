package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"finance-tracker/internal/infrastructure/notion"
	"finance-tracker/internal/infrastructure/telegram"
	"finance-tracker/internal/usecase"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN no está seteado")
	}

	writer := notion.NewClient(os.Getenv("NOTION_API_KEY"), os.Getenv("NOTION_DATABASE_ID"))
	resolver := notion.NewCategoryResolver(os.Getenv("NOTION_API_KEY"), os.Getenv("NOTION_CATEGORY_DATA_SOURCE_ID"))
	registrar := usecase.NewRegistrar(resolver, writer)
	bot := telegram.NewBot(token)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("bot arrancado, long-polling contra Telegram...")
	runPollingLoop(ctx, bot, registrar)
	log.Println("bot detenido")
}

// runPollingLoop hace long-polling contra getUpdates hasta que ctx se
// cancele. Ante error de red/API reintenta con backoff exponencial acotado
// en vez de terminar el proceso.
func runPollingLoop(ctx context.Context, bot *telegram.Bot, registrar *usecase.Registrar) {
	const (
		initialBackoff = 2 * time.Second
		maxBackoff     = 30 * time.Second
	)

	offset := 0
	backoff := initialBackoff

	for {
		if ctx.Err() != nil {
			return
		}

		updates, err := bot.GetUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("error en getUpdates: %v (reintenta en %s)", err, backoff)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return
			}
			if backoff < maxBackoff {
				backoff *= 2
			}
			continue
		}
		backoff = initialBackoff

		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message == nil || update.Message.Text == "" {
				continue
			}
			handleMessage(ctx, bot, registrar, update.Message)
		}
	}
}

// handleMessage responde con el ID de página de Notion en éxito, o el
// error crudo en falla. Es un placeholder funcional — la feature 004 define
// el copy exacto y evita exponer errores internos de Notion al chat.
func handleMessage(ctx context.Context, bot *telegram.Bot, registrar *usecase.Registrar, msg *telegram.Message) {
	pageID, err := registrar.Register(msg.Text)
	if err != nil {
		log.Printf("error registrando mensaje de chat %d: %v", msg.Chat.ID, err)
		if sendErr := bot.SendMessage(ctx, msg.Chat.ID, fmt.Sprintf("Error: %v", err)); sendErr != nil {
			log.Printf("error respondiendo: %v", sendErr)
		}
		return
	}

	reply := fmt.Sprintf("Registrado. Notion: https://notion.so/%s", pageID)
	if err := bot.SendMessage(ctx, msg.Chat.ID, reply); err != nil {
		log.Printf("error respondiendo: %v", err)
	}
}
