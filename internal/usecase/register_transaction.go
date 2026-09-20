package usecase

import (
	"strings"
	"time"
)

// Registrar orquesta el registro de un movimiento: parsea el mensaje,
// resuelve la categoría (si vino) contra Notion, y escribe la página.
type Registrar struct {
	categories CategoryResolver
	writer     ExpenseWriter
}

func NewRegistrar(categories CategoryResolver, writer ExpenseWriter) *Registrar {
	return &Registrar{categories: categories, writer: writer}
}

// Register procesa un mensaje de Telegram y devuelve el ID de la página
// creada en Notion.
func (r *Registrar) Register(text string) (pageID string, err error) {
	parsed, err := ParseMessage(text, time.Now())
	if err != nil {
		return "", err
	}

	tx := parsed.Transaction

	if name := strings.TrimSpace(parsed.CategoryName); name != "" {
		categoryID, ok, err := r.categories.Resolve(name)
		if err != nil {
			return "", ErrWriteFailed{Err: err}
		}
		if !ok {
			return "", ErrUnrecognizedCategory{Got: name}
		}
		tx.CategoryID = categoryID
	}

	pageID, err = r.writer.CreateExpense(tx)
	if err != nil {
		return "", ErrWriteFailed{Err: err}
	}
	return pageID, nil
}
