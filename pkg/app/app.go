// Package app arma las dependencias reales del bot. Es el único lugar que
// conoce a la vez los servicios concretos (Telegram, Notion, Groq) y el
// handler, para que cada punto de entrada (long-polling, webhook) no
// repita el cableado.
package app

import (
	"log"

	"finance-tracker/pkg/config"
	groq "finance-tracker/pkg/groq/services"
	notion "finance-tracker/pkg/notion/services"
	"finance-tracker/pkg/ports"
	"finance-tracker/pkg/telegram/handler"
	"finance-tracker/pkg/telegram/model/command"
	telegram "finance-tracker/pkg/telegram/services"
	transaction "finance-tracker/pkg/transaction/services"
)

// Commands es el menú de comandos que se registra en Telegram: solo
// descubribilidad, el bot reconoce los comandos por texto igual.
var Commands = []command.BotCommand{
	{Command: "egreso", Description: "Registrar un egreso paso a paso"},
	{Command: "ingreso", Description: "Registrar un ingreso paso a paso"},
	{Command: "resumen", Description: "Ver el balance y los totales del mes"},
	{Command: "ultimo", Description: "Ver el último movimiento registrado"},
	{Command: "categorias", Description: "Ver o agregar categorías"},
	{Command: "ayuda", Description: "Formato de mensaje y comandos disponibles"},
}

// App agrupa lo que un punto de entrada necesita: el handler y el cliente
// de Telegram concreto (para GetUpdates, webhook y menú de comandos, que
// no forman parte de ports.Messenger).
type App struct {
	Handler  *handler.Handler
	Telegram *telegram.Services
}

// New cablea los servicios con store como almacén del estado pendiente. Una
// whitelist vacía deja al bot sin procesar ningún mensaje (fail closed).
func New(cfg *config.Config, store ports.ChatStore) *App {
	if len(cfg.AllowedChatIDs) == 0 {
		log.Println("advertencia: TELEGRAM_ALLOWED_CHAT_IDS vacío, el bot no va a procesar ningún mensaje (fail closed)")
	} else {
		log.Printf("whitelist cargada: %d chat_id autorizado(s)", len(cfg.AllowedChatIDs))
	}

	tg := telegram.NewServices(cfg)
	notionServices := notion.NewServices(cfg)
	groqServices := groq.NewServices(cfg)

	// Notion cubre varios puertos: escribir el movimiento, subir el recibo,
	// sumar el mes, buscar el último y resolver categorías.
	transactions := transaction.NewServices(notionServices, notionServices, notionServices, groqServices)

	return &App{
		Handler:  handler.NewHandler(cfg, tg, transactions, notionServices, notionServices, notionServices, store),
		Telegram: tg,
	}
}
