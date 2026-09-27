package ports

import (
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

// TransactionServices registra movimientos a partir de las distintas
// entradas del bot. Implementada por pkg/transaction/services.
type TransactionServices interface {
	Register(text string) (register.Result, error)
	RegisterWithPhoto(text string, photo []byte, filename, contentType string) (register.Result, error)
	RegisterFromPhotoAI(photo []byte, filename, contentType string, txType transaction.Type) (register.Result, error)
	RegisterFromFields(input register.ConversationInput) (register.Result, error)
}
