package ports

import (
	"context"

	"finance-tracker/pkg/telegram/model/keyboard"
)

// Messenger es lo que el handler de Telegram necesita del canal de
// mensajería.
type Messenger interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
	// SendKeyboard manda un mensaje con botones y devuelve su message_id,
	// para poder reconocerlo cuando el usuario toque un botón.
	SendKeyboard(ctx context.Context, chatID int64, text string, buttons []keyboard.InlineButton) (messageID int, err error)
	EditMessageText(ctx context.Context, chatID int64, messageID int, text string) error
	// AnswerCallbackQuery le avisa al cliente que ya se procesó el toque
	// del botón; sin esto queda con el ícono de "cargando".
	AnswerCallbackQuery(ctx context.Context, callbackQueryID string) error
	// DownloadPhoto descarga el contenido real de una foto a partir de su
	// file_id: los file_id no son URLs permanentes.
	DownloadPhoto(ctx context.Context, fileID string) ([]byte, error)
}
