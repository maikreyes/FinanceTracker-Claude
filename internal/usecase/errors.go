package usecase

import "fmt"

// ErrInvalidFormat indica que el mensaje no sigue la gramática esperada:
// <Tipo> [Concepto] <Monto> <MedioDePago> [Categoria].
type ErrInvalidFormat struct {
	Reason string
}

func (e ErrInvalidFormat) Error() string {
	return fmt.Sprintf("formato de mensaje inválido: %s", e.Reason)
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
// coincide con ninguna categoría existente en Notion.
type ErrUnrecognizedCategory struct {
	Got string
}

func (e ErrUnrecognizedCategory) Error() string {
	return fmt.Sprintf("categoría no reconocida: %q", e.Got)
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
