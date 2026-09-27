package ports

import "finance-tracker/pkg/telegram/model/chatstate"

// ChatStore persiste el ChatState de cada chat entre updates. Get de un
// chat sin estado pendiente devuelve el ChatState cero, no un error: "sin
// nada pendiente" es un resultado válido, no una falla.
//
// Un ChatState que sale de Get es una copia: modificarlo no cambia lo
// guardado hasta que se llame a Set.
type ChatStore interface {
	Get(chatID int64) (chatstate.ChatState, error)
	Set(chatID int64, state chatstate.ChatState) error
	Delete(chatID int64) error
	// Take devuelve el estado y lo borra en una sola operación atómica. Es
	// lo que hay que usar para consumir una espera una única vez: con
	// updates concurrentes (webhook), Get seguido de Delete dejaría que
	// dos toques del mismo botón registren dos veces.
	Take(chatID int64) (chatstate.ChatState, error)
	// MarkUpdateSeen reserva el update_id y dice si es la primera vez que
	// se ve. Telegram reentrega un update si el webhook tarda en
	// responder, y procesarlo dos veces registraría el movimiento dos
	// veces.
	MarkUpdateSeen(updateID int) (firstTime bool, err error)
}
