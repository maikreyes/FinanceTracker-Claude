package update

// Update es lo que Telegram entrega por getUpdates o por webhook.
type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"` // se llena cuando el usuario toca un botón inline
}

type Message struct {
	MessageID int         `json:"message_id"`
	Chat      Chat        `json:"chat"`
	Text      string      `json:"text"`
	Caption   string      `json:"caption"` // texto de una foto, no viaja en "text"
	Photo     []PhotoSize `json:"photo"`   // varias resoluciones de la misma foto; la última es la de mayor calidad
}

// CallbackQuery es lo que Telegram manda cuando el usuario toca un botón
// inline. Message es el mensaje que tenía el botón: se necesita su
// Chat/MessageID para poder editarlo después.
type CallbackQuery struct {
	ID      string   `json:"id"`
	Data    string   `json:"data"`
	Message *Message `json:"message"`
}

type PhotoSize struct {
	FileID string `json:"file_id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Chat struct {
	ID int64 `json:"id"`
}
