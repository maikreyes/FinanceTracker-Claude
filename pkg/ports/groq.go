package ports

import "finance-tracker/pkg/transaction/model/receipt"

// ReceiptInterpreter interpreta la foto de un recibo con un modelo de IA y
// devuelve los datos del movimiento que pudo extraer. categoryNames son las
// categorías reales del usuario, para que la IA elija solo entre esas y no
// invente.
type ReceiptInterpreter interface {
	Interpret(imageData []byte, categoryNames []string) (receipt.Interpreted, error)
}
