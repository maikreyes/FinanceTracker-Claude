package usecase

import (
	"strconv"
	"strings"
	"time"

	"finance-tracker/internal/domain"
)

var typeKeywords = map[string]domain.TransactionType{
	"egreso":  domain.TransactionTypeExpense,
	"ingreso": domain.TransactionTypeIncome,
}

// paymentMethodPhrases se prueban en este orden para que una frase más
// larga ("tarjeta de credito") se reconozca antes que una más corta
// ("tarjeta credito") que la contiene.
var paymentMethodPhrases = []struct {
	words  []string
	method domain.PaymentMethod
}{
	{[]string{"tarjeta", "de", "credito"}, domain.PaymentMethodCreditCard},
	{[]string{"tarjeta", "de", "crédito"}, domain.PaymentMethodCreditCard},
	{[]string{"tarjeta", "de", "debito"}, domain.PaymentMethodDebitCard},
	{[]string{"tarjeta", "de", "débito"}, domain.PaymentMethodDebitCard},
	{[]string{"tarjeta", "credito"}, domain.PaymentMethodCreditCard},
	{[]string{"tarjeta", "crédito"}, domain.PaymentMethodCreditCard},
	{[]string{"tarjeta", "debito"}, domain.PaymentMethodDebitCard},
	{[]string{"tarjeta", "débito"}, domain.PaymentMethodDebitCard},
	{[]string{"efectivo"}, domain.PaymentMethodCash},
	{[]string{"banco"}, domain.PaymentMethodBank},
	{[]string{"tranferencia"}, domain.PaymentMethodBank},
	{[]string{"transferencia"}, domain.PaymentMethodBank},
}

// ParsedMessage es la salida cruda de ParseMessage: Transaction trae todo
// lo que ya mapea 1:1 a Notion; CategoryName es el texto de categoría tal
// como vino en el mensaje (vacío si no hubo), sin resolver — eso lo hace
// Registrar contra un CategoryResolver real.
type ParsedMessage struct {
	Transaction  domain.Transaction
	CategoryName string
}

// ParseMessage interpreta un mensaje con la gramática documentada en
// specs/features/001-registrar-movimiento-por-mensaje/spec.md:
//
//	<Tipo> [Concepto] <Monto> <MedioDePago> [Categoria]
//
// now es la fecha a asignar a Transaction.Date (fecha de recepción del
// mensaje, el formato no trae fecha explícita).
func ParseMessage(text string, now time.Time) (ParsedMessage, error) {
	tokens := strings.Fields(text)
	if len(tokens) < 3 {
		return ParsedMessage{}, ErrInvalidFormat{Reason: "se esperan al menos tipo, monto y medio de pago"}
	}

	txType, ok := typeKeywords[strings.ToLower(tokens[0])]
	if !ok {
		return ParsedMessage{}, ErrUnrecognizedType{Got: tokens[0]}
	}

	amountIndex := -1
	var amount float64
	for i := 1; i < len(tokens); i++ {
		if v, err := strconv.ParseFloat(tokens[i], 64); err == nil {
			amountIndex = i
			amount = v
			break
		}
	}
	if amountIndex == -1 {
		return ParsedMessage{}, ErrInvalidFormat{Reason: "no se encontró un monto numérico"}
	}

	concepto := strings.Join(tokens[1:amountIndex], " ")
	if concepto == "" {
		concepto = string(txType)
	}

	rest := tokens[amountIndex+1:]
	method, consumed, ok := matchPaymentMethod(rest)
	if !ok {
		got := ""
		if len(rest) > 0 {
			got = rest[0]
		}
		return ParsedMessage{}, ErrUnrecognizedPaymentMethod{Got: got}
	}

	categoryName := strings.Join(rest[consumed:], " ")

	return ParsedMessage{
		Transaction: domain.Transaction{
			Name:          concepto,
			Amount:        amount,
			Type:          txType,
			Date:          now,
			PaymentMethod: method,
		},
		CategoryName: categoryName,
	}, nil
}

func matchPaymentMethod(tokens []string) (domain.PaymentMethod, int, bool) {
	for _, phrase := range paymentMethodPhrases {
		n := len(phrase.words)
		if n > len(tokens) || !phraseMatches(tokens[:n], phrase.words) {
			continue
		}
		return phrase.method, n, true
	}
	return "", 0, false
}

func phraseMatches(got, want []string) bool {
	for i, w := range want {
		if strings.ToLower(got[i]) != w {
			return false
		}
	}
	return true
}
