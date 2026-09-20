package usecase

import "finance-tracker/internal/domain"

// CategoryResolver resuelve el nombre de una categoría (tal como viene en
// el mensaje) al ID de su página en la data source "Category" de Notion.
// Implementada en internal/infrastructure/notion — el usecase no depende
// del cliente concreto, solo de esta interfaz.
type CategoryResolver interface {
	Resolve(name string) (categoryID string, ok bool, err error)
}

// ExpenseWriter crea la página del movimiento en la data source
// "Expenses" de Notion. Implementada por notion.Client.
type ExpenseWriter interface {
	CreateExpense(tx domain.Transaction) (pageID string, err error)
}
