package services

import (
	"math"
	"strconv"
	"strings"
	"time"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

// ParseMessage interpreta un mensaje con la gramática documentada en
// specs/features/001-registrar-movimiento-por-mensaje/spec.md:
//
//	<Tipo> [Concepto] <Monto> <MedioDePago> [Categoria]
//
// now es la fecha a asignar a Transaction.Date (fecha de recepción del
// mensaje, el formato no trae fecha explícita).
func ParseMessage(text string, now time.Time) (register.ParsedMessage, error) {
	tokens := strings.Fields(text)
	if len(tokens) < 3 {
		return register.ParsedMessage{}, failure.ErrInvalidFormat{Reason: "se esperan al menos tipo, monto y medio de pago"}
	}

	txType, ok := transaction.TypeFromKeyword(tokens[0])
	if !ok {
		return register.ParsedMessage{}, failure.ErrUnrecognizedType{Got: tokens[0]}
	}

	amountIndex := -1
	var amount float64
	for i := 1; i < len(tokens); i++ {
		// ParseFloat acepta "nan"/"inf"/"infinity": un concepto con esas
		// palabras no es un monto.
		if v, err := strconv.ParseFloat(tokens[i], 64); err == nil && !math.IsNaN(v) && !math.IsInf(v, 0) {
			amountIndex = i
			amount = v
			break
		}
	}
	if amountIndex == -1 {
		return register.ParsedMessage{}, failure.ErrInvalidFormat{Reason: "no se encontró un monto numérico"}
	}
	if !transaction.ValidAmount(amount) {
		return register.ParsedMessage{}, failure.ErrInvalidAmount{Got: amount}
	}

	concepto := strings.Join(tokens[1:amountIndex], " ")
	if concepto == "" {
		concepto = string(txType)
	}

	rest := tokens[amountIndex+1:]
	method, consumed, ok := transaction.MatchPaymentMethod(rest)
	if !ok {
		got := ""
		if len(rest) > 0 {
			got = rest[0]
		}
		return register.ParsedMessage{}, failure.ErrUnrecognizedPaymentMethod{Got: got}
	}

	return register.ParsedMessage{
		Transaction: transaction.Transaction{
			Name:          concepto,
			Amount:        amount,
			Type:          txType,
			Date:          now,
			PaymentMethod: method,
		},
		CategoryName: strings.Join(rest[consumed:], " "),
	}, nil
}
