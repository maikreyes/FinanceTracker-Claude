package transaction

import (
	"math"
	"strings"
)

var typeKeywords = map[string]Type{
	"egreso":  TypeExpense,
	"ingreso": TypeIncome,
}

// validTypes es la lista "de exhibición" para mensajes de error.
var validTypes = []string{"egreso", "ingreso"}

// ValidTypes devuelve los valores aceptados para <Tipo>.
func ValidTypes() []string { return validTypes }

// TypeFromKeyword reconoce "egreso" o "ingreso" sin distinguir mayúsculas.
func TypeFromKeyword(word string) (Type, bool) {
	t, ok := typeKeywords[strings.ToLower(word)]
	return t, ok
}

// ValidAmount dice si amount sirve como monto de un movimiento: positivo y
// finito. El signo de un movimiento lo da su Type, nunca el número.
func ValidAmount(amount float64) bool {
	return amount > 0 && !math.IsNaN(amount) && !math.IsInf(amount, 0)
}
