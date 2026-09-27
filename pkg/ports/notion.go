package ports

import (
	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

// CategoryResolver resuelve el nombre de una categoría (tal como viene en
// el mensaje) al ID de su página en la data source "Category" de Notion.
type CategoryResolver interface {
	Resolve(name string) (categoryID string, ok bool, err error)
	// ListNames devuelve los nombres de las categorías reales existentes
	// en Notion.
	ListNames() ([]string, error)
	// Create crea una categoría nueva y devuelve su ID. Si ya existe una
	// con ese nombre (case-insensitive), no crea una página duplicada:
	// devuelve el ID de la existente.
	Create(name string) (categoryID string, err error)
}

// ExpenseWriter crea la página del movimiento en la data source
// "Expenses" de Notion.
type ExpenseWriter interface {
	CreateExpense(tx transaction.Transaction) (pageID string, err error)
}

// ReceiptUploader sube un archivo (foto de un recibo) y devuelve un ID
// para referenciarlo en la propiedad "Receipt" del movimiento.
type ReceiptUploader interface {
	Upload(data []byte, filename, contentType string) (fileUploadID string, err error)
}

// MonthlySummarizer suma los movimientos del mes calendario actual,
// separados por tipo.
type MonthlySummarizer interface {
	SumThisMonth() (summary.Monthly, error)
}

// LatestTransactionFinder busca el movimiento registrado más
// recientemente. El bool distingue "no hay ningún movimiento todavía"
// (false, sin error) de una falla real.
type LatestTransactionFinder interface {
	FindLatest() (summary.Latest, bool, error)
}
