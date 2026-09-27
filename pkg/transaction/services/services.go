package services

import (
	"time"

	"finance-tracker/pkg/ports"
)

// Services orquesta el registro de un movimiento: parsea el mensaje,
// resuelve la categoría (si vino) contra Notion, y escribe la página.
type Services struct {
	categories  ports.CategoryResolver
	writer      ports.ExpenseWriter
	uploader    ports.ReceiptUploader
	interpreter ports.ReceiptInterpreter
	now         func() time.Time
}

func NewServices(categories ports.CategoryResolver, writer ports.ExpenseWriter, uploader ports.ReceiptUploader, interpreter ports.ReceiptInterpreter) *Services {
	return &Services{categories: categories, writer: writer, uploader: uploader, interpreter: interpreter, now: nowBogota}
}

var _ ports.TransactionServices = (*Services)(nil)
