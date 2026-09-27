package command

// BotCommand es una entrada del menú de comandos de Telegram.
type BotCommand struct {
	Command     string // sin la "/": Telegram la agrega sola en el menú
	Description string
}
