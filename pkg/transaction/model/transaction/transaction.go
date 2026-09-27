package transaction

import "time"

// PaymentMethod mirrors the fixed option set of the "Payment Method"
// select property in the Notion "Expenses" data source.
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "Credit Card"
	PaymentMethodDebitCard  PaymentMethod = "Debit Card"
	PaymentMethodBank       PaymentMethod = "Bank"
	PaymentMethodCash       PaymentMethod = "Cash"
)

// Type mirrors the "Type" select property added to the Notion "Expenses"
// data source on 2026-09-20: the original schema only supported expenses
// (see .claude/rules/notion-rules.md).
type Type string

const (
	TypeExpense Type = "Egreso"
	TypeIncome  Type = "Ingreso"
)

// Transaction is the core entity for a parsed movement (ingreso o egreso).
// Field set y tipos calcan 1:1 el esquema real de la data source
// "Expenses" leído vía Notion MCP (ver .claude/rules/notion-rules.md).
type Transaction struct {
	Name           string        // title property "Name" — concepto del movimiento
	Amount         float64       // number property "Amount" (formato colombian_peso), siempre positivo
	Type           Type          // select property "Type"
	Date           time.Time     // date property "Date"
	CategoryID     string        // relation property "Category" — page ID en la data source "Category" (limit: 1), opcional
	PaymentMethod  PaymentMethod // select property "Payment Method"
	Notes          string        // text property "Notes", opcional
	ReceiptFileIDs []string      // file property "Receipt" — IDs de file_upload de Notion, opcional
}
