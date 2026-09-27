package services

import (
	"strings"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

// RegisterFromPhotoAI interpreta photo con IA, sin caption, a diferencia
// de RegisterWithPhoto. txType viene de la respuesta del usuario a los
// botones Ingreso/Egreso: la IA no decide el tipo, una imagen sola no
// alcanza para inferirlo con confianza. Solo el monto es obligatorio: sin
// él, no se registra nada. Categoría/medio de pago que la IA no determinó
// (o sugirió algo que no matchea) se dejan vacíos y no bloquean el
// registro, a diferencia de una categoría escrita a mano por el usuario.
func (s *Services) RegisterFromPhotoAI(photo []byte, filename, contentType string, txType transaction.Type) (register.Result, error) {
	categoryNames, err := s.categories.ListNames()
	if err != nil {
		return register.Result{}, failure.ErrWriteFailed{Err: err}
	}

	interpreted, err := s.interpreter.Interpret(photo, categoryNames)
	if err != nil {
		return register.Result{}, failure.ErrReceiptInterpretationFailed{Err: err}
	}
	if interpreted.Amount == nil || !transaction.ValidAmount(*interpreted.Amount) {
		return register.Result{}, failure.ErrReceiptAmountNotFound{}
	}

	name := strings.TrimSpace(interpreted.Concept)
	if name == "" {
		name = string(txType)
	}

	tx := transaction.Transaction{
		Name:   name,
		Amount: *interpreted.Amount,
		Type:   txType,
		Date:   s.now(),
	}

	if method, ok := transaction.MatchPaymentMethodName(interpreted.PaymentMethod); ok {
		tx.PaymentMethod = method
	}

	categoryName := ""
	if raw := strings.TrimSpace(interpreted.CategoryName); raw != "" {
		// Un error de red al resolver la sugerencia de la IA no justifica
		// perder el registro entero: se guarda sin categoría, igual que si
		// la sugerencia no hubiera matcheado.
		if categoryID, ok, err := s.categories.Resolve(raw); err == nil && ok {
			tx.CategoryID = categoryID
			categoryName = raw
		}
	}

	return s.finish(tx, categoryName, &receiptFile{data: photo, filename: filename, contentType: contentType}, true)
}
