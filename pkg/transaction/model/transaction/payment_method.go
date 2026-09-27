package transaction

import "strings"

// paymentMethodPhrases se prueban en este orden para que una frase más
// larga ("tarjeta de credito") se reconozca antes que una más corta
// ("tarjeta credito") que la contiene.
var paymentMethodPhrases = []struct {
	words  []string
	method PaymentMethod
}{
	{[]string{"tarjeta", "de", "credito"}, PaymentMethodCreditCard},
	{[]string{"tarjeta", "de", "crédito"}, PaymentMethodCreditCard},
	{[]string{"tarjeta", "de", "debito"}, PaymentMethodDebitCard},
	{[]string{"tarjeta", "de", "débito"}, PaymentMethodDebitCard},
	{[]string{"tarjeta", "credito"}, PaymentMethodCreditCard},
	{[]string{"tarjeta", "crédito"}, PaymentMethodCreditCard},
	{[]string{"tarjeta", "debito"}, PaymentMethodDebitCard},
	{[]string{"tarjeta", "débito"}, PaymentMethodDebitCard},
	{[]string{"efectivo"}, PaymentMethodCash},
	{[]string{"banco"}, PaymentMethodBank},
	{[]string{"tranferencia"}, PaymentMethodBank},
	{[]string{"transferencia"}, PaymentMethodBank},
}

// validPaymentMethods es la lista "de exhibición" para mensajes de error:
// una forma canónica de cada opción, no todos los alias que acepta el
// parser.
var validPaymentMethods = []string{
	"efectivo", "tarjeta de credito", "tarjeta de debito", "banco", "transferencia",
}

// ValidPaymentMethods devuelve los valores aceptados para <MedioDePago>.
func ValidPaymentMethods() []string { return validPaymentMethods }

// MatchPaymentMethod reconoce una frase de medio de pago al inicio de
// tokens y devuelve cuántos tokens consumió.
func MatchPaymentMethod(tokens []string) (PaymentMethod, int, bool) {
	for _, phrase := range paymentMethodPhrases {
		n := len(phrase.words)
		if n > len(tokens) || !phraseMatches(tokens[:n], phrase.words) {
			continue
		}
		return phrase.method, n, true
	}
	return "", 0, false
}

// ParsePaymentMethod valida una respuesta de texto libre contra las frases
// conocidas. A diferencia de MatchPaymentMethod, exige que el texto
// completo sea la frase, para el paso de medio de pago del flujo guiado.
func ParsePaymentMethod(text string) (PaymentMethod, bool) {
	tokens := strings.Fields(text)
	method, consumed, ok := MatchPaymentMethod(tokens)
	if !ok || consumed != len(tokens) {
		return "", false
	}
	return method, true
}

func phraseMatches(got, want []string) bool {
	for i, w := range want {
		if strings.ToLower(got[i]) != w {
			return false
		}
	}
	return true
}

// MatchPaymentMethodName valida el valor de "medio_de_pago" que devuelve
// el intérprete de IA contra las 4 opciones reales de Notion. Espera
// exactamente uno de esos valores en inglés (ver el prompt en
// pkg/groq/services).
func MatchPaymentMethodName(s string) (PaymentMethod, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "cash":
		return PaymentMethodCash, true
	case "credit card":
		return PaymentMethodCreditCard, true
	case "debit card":
		return PaymentMethodDebitCard, true
	case "bank":
		return PaymentMethodBank, true
	default:
		return "", false
	}
}
