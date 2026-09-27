package services

import (
	"strings"

	"finance-tracker/pkg/transaction/model/register"
)

// Register procesa un mensaje de Telegram y devuelve el resultado del
// registro en Notion.
func (s *Services) Register(text string) (register.Result, error) {
	return s.registerText(text, nil)
}

// RegisterWithPhoto es como Register, pero además sube photo (los bytes de
// una foto de recibo) y la adjunta a la propiedad "Receipt" del movimiento.
// text sigue siendo el caption con el formato de siempre: esta función no
// interpreta la imagen.
func (s *Services) RegisterWithPhoto(text string, photo []byte, filename, contentType string) (register.Result, error) {
	return s.registerText(text, &receiptFile{data: photo, filename: filename, contentType: contentType})
}

func (s *Services) registerText(text string, receipt *receiptFile) (register.Result, error) {
	parsed, err := ParseMessage(text, s.now())
	if err != nil {
		return register.Result{}, err
	}

	tx := parsed.Transaction
	categoryName := strings.TrimSpace(parsed.CategoryName)
	if categoryName != "" {
		categoryID, err := s.resolveCategory(categoryName)
		if err != nil {
			return register.Result{}, err
		}
		tx.CategoryID = categoryID
	}

	return s.finish(tx, categoryName, receipt, false)
}
