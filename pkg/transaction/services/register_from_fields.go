package services

import (
	"strings"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

// RegisterFromFields registra un movimiento armado a partir de
// ConversationInput. No pasa por ParseMessage: los campos ya vienen
// separados, y reconstruir una línea de texto para volver a parsearla
// reintroduciría la ambigüedad posicional que el flujo guiado evita.
func (s *Services) RegisterFromFields(input register.ConversationInput) (register.Result, error) {
	if !transaction.ValidAmount(input.Amount) {
		return register.Result{}, failure.ErrInvalidAmount{Got: input.Amount}
	}

	name := strings.TrimSpace(input.Concept)
	if name == "" {
		name = string(input.Type)
	}

	tx := transaction.Transaction{
		Name:          name,
		Amount:        input.Amount,
		Type:          input.Type,
		Date:          s.now(),
		PaymentMethod: input.PaymentMethod,
	}

	categoryName := strings.TrimSpace(input.CategoryName)
	if categoryName != "" {
		categoryID, err := s.resolveCategory(categoryName)
		if err != nil {
			return register.Result{}, err
		}
		tx.CategoryID = categoryID
	}

	return s.finish(tx, categoryName, nil, false)
}
