package api

import (
	"log"
	"net/http"
	"sync"

	"finance-tracker/pkg/app"
	"finance-tracker/pkg/config"
	store "finance-tracker/pkg/store/services"
	"finance-tracker/pkg/telegram/handler"
)

var (
	initOnce   sync.Once
	webhook    *handler.Handler
	initFailed error
)

// initApp cablea las dependencias una sola vez por instancia de la
// función: Vercel reutiliza la instancia entre invocaciones, y así se
// conserva la caché de categorías. Un error de configuración se recuerda:
// las variables de entorno no cambian mientras la instancia vive.
func initApp() error {
	initOnce.Do(func() {
		cfg := config.NewConfig()
		if initFailed = cfg.CheckWebhook(); initFailed != nil {
			return
		}
		webhook = app.New(cfg, store.NewUpstash(cfg)).Handler
	})
	return initFailed
}

// Handler es la función serverless de Vercel: Telegram le manda un POST por
// cada update.
func Handler(w http.ResponseWriter, r *http.Request) {
	if err := initApp(); err != nil {
		// El detalle se loguea, no se devuelve: el endpoint es público.
		log.Printf("configuración inválida: %v", err)
		http.Error(w, "Server misconfigured", http.StatusInternalServerError)
		return
	}
	webhook.WebhookHandler(w, r)
}
