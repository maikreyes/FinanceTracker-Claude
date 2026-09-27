package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"finance-tracker/pkg/app"
	"finance-tracker/pkg/config"
	store "finance-tracker/pkg/store/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no se cargó .env, se usa el entorno del proceso: %v", err)
	}

	cfg := config.NewConfig()
	if err := cfg.CheckTelegram(); err != nil {
		log.Fatal(err)
	}

	// El estado pendiente vive en memoria de proceso: un único proceso
	// persistente procesa los updates de forma secuencial.
	a := app.New(cfg, store.NewMemory())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.Telegram.SetMyCommands(ctx, app.Commands); err != nil {
		log.Printf("advertencia: no se pudo registrar el menú de comandos: %v", err)
	}

	log.Println("bot arrancado, long-polling contra Telegram...")
	runPollingLoop(ctx, a)
	log.Println("bot detenido")
}

// runPollingLoop hace long-polling contra getUpdates hasta que ctx se
// cancele. Ante error de red/API reintenta con backoff exponencial acotado
// en vez de terminar el proceso. Toda la lógica de negocio vive en
// Handler.Dispatch: este loop solo trae updates y se los pasa.
func runPollingLoop(ctx context.Context, a *app.App) {
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

		updates, err := a.Telegram.GetUpdates(ctx, offset)
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

		for _, u := range updates {
			offset = u.UpdateID + 1
			a.Handler.Dispatch(ctx, u)
		}
	}
}
