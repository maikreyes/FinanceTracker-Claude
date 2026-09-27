package failure

import (
	"fmt"
	"strings"
)

// ErrInvalidFormat indica que el mensaje no sigue la gramática esperada:
// <Tipo> [Concepto] <Monto> <MedioDePago> [Categoria].
type ErrInvalidFormat struct {
	Reason string
}

func (e ErrInvalidFormat) Error() string {
	return fmt.Sprintf("formato de mensaje inválido: %s", e.Reason)
}

// ErrInvalidAmount indica que el monto es cero, negativo o no es un número
// finito.
type ErrInvalidAmount struct {
	Got float64
}

func (e ErrInvalidAmount) Error() string {
	return fmt.Sprintf("monto inválido: %v (debe ser un número mayor a 0)", e.Got)
}

// ErrUnrecognizedType indica que el primer token del mensaje no es
// "egreso" ni "ingreso".
type ErrUnrecognizedType struct {
	Got string
}

func (e ErrUnrecognizedType) Error() string {
	return fmt.Sprintf("tipo no reconocido: %q (esperado \"egreso\" o \"ingreso\")", e.Got)
}

// ErrUnrecognizedPaymentMethod indica que ningún token coincidió con una
// opción válida de "Payment Method".
type ErrUnrecognizedPaymentMethod struct {
	Got string
}

func (e ErrUnrecognizedPaymentMethod) Error() string {
	return fmt.Sprintf("medio de pago no reconocido: %q", e.Got)
}

// ErrUnrecognizedCategory indica que el texto de categoría del mensaje no
// coincide con ninguna categoría existente en Notion. Valid trae la lista
// de categorías reales al momento del error (consultadas en vivo), para
// que quien maneje el error pueda mostrarla sin volver a consultar Notion.
type ErrUnrecognizedCategory struct {
	Got   string
	Valid []string
}

func (e ErrUnrecognizedCategory) Error() string {
	if len(e.Valid) == 0 {
		return fmt.Sprintf("categoría no reconocida: %q", e.Got)
	}
	return fmt.Sprintf("categoría no reconocida: %q (válidas: %s)", e.Got, strings.Join(e.Valid, ", "))
}

// ErrWriteFailed envuelve un error de comunicación con la API de Notion
// (falla al resolver categoría o al crear la página).
type ErrWriteFailed struct {
	Err error
}

func (e ErrWriteFailed) Error() string {
	return fmt.Sprintf("no se pudo guardar el movimiento: %v", e.Err)
}

func (e ErrWriteFailed) Unwrap() error {
	return e.Err
}

// ErrReceiptAmountNotFound indica que la IA no logró identificar un monto
// en la foto del recibo — distinto de una falla de la API,
// es la imagen la que no trae un dato claro.
type ErrReceiptAmountNotFound struct{}

func (e ErrReceiptAmountNotFound) Error() string {
	return "no se identificó un monto en la foto del recibo"
}

// ErrReceiptInterpretationFailed envuelve una falla real (red, cuota, error
// del servicio) al llamar al modelo de IA — no expone el error real al
// chat, mismo patrón que ErrWriteFailed.
type ErrReceiptInterpretationFailed struct {
	Err error
}

func (e ErrReceiptInterpretationFailed) Error() string {
	return fmt.Sprintf("no se pudo interpretar la foto del recibo: %v", e.Err)
}

func (e ErrReceiptInterpretationFailed) Unwrap() error {
	return e.Err
}
