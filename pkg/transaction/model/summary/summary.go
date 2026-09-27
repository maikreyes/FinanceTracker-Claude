package summary

import "finance-tracker/pkg/transaction/model/transaction"

// Monthly es el resultado de sumar los movimientos del mes calendario
// actual, separados por tipo.
type Monthly struct {
	Egreso       float64
	Ingreso      float64
	CountEgreso  int
	CountIngreso int
}

// Balance calcula Ingreso - Egreso: "cuánto tengo disponible", la
// convención estándar.
func (m Monthly) Balance() float64 {
	return m.Ingreso - m.Egreso
}

// Latest es una vista de solo lectura del movimiento más reciente. Date
// queda como el string ISO tal como lo devuelve Notion: no hay necesidad
// de parsearlo a time.Time solo para mostrarlo.
type Latest struct {
	PageID        string
	Name          string
	Amount        float64
	Type          transaction.Type
	PaymentMethod transaction.PaymentMethod
	CategoryName  string
	Date          string
}
