package handler

import (
	"errors"
	"fmt"
	"strings"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

const ayudaText = `Formato de un mensaje: tipo [concepto] monto medio_de_pago [categoria]
Ejemplo: egreso mecato 13000 efectivo comida

Tipos válidos: Egreso, Ingreso
Medios de pago válidos: efectivo, tarjeta de credito, tarjeta de debito, banco, transferencia

También podés:
- Mandar una foto de un recibo (con o sin texto)
- /egreso o /ingreso: registrar paso a paso
- /resumen: balance y totales del mes
- /ultimo: ver el último movimiento
- /categorias: ver o agregar categorías`

// successReply confirma qué se guardó, con los datos reales de la
// transacción registrada y un link a la página creada en Notion. Si vino
// de IA, antepone un disclaimer — la confianza no es la misma que un
// registro por texto exacto.
func successReply(result register.Result) string {
	tx := result.Transaction
	var lines []string
	if result.Interpreted {
		lines = append(lines, "Interpretado por IA — revisá los datos en Notion si algo no coincide.")
	}
	lines = append(lines,
		fmt.Sprintf("Registrado: %s", tx.Type),
		fmt.Sprintf("Concepto: %s", tx.Name),
		fmt.Sprintf("Monto: %.0f", tx.Amount),
	)
	if tx.PaymentMethod != "" {
		lines = append(lines, fmt.Sprintf("Medio de pago: %s", tx.PaymentMethod))
	}
	if result.CategoryName != "" {
		lines = append(lines, fmt.Sprintf("Categoría: %s", result.CategoryName))
	}
	lines = append(lines, fmt.Sprintf("Notion: https://notion.so/%s", result.PageID))
	return strings.Join(lines, "\n")
}

// errorReply mapea los errores tipados del usecase a texto accionable en
// español. ErrWriteFailed nunca expone el error real de Notion al chat —
// solo se loguea server-side.
func errorReply(err error) string {
	var errType failure.ErrUnrecognizedType
	if errors.As(err, &errType) {
		return fmt.Sprintf(
			"No entendí el tipo %q. El mensaje empieza con \"egreso\" o \"ingreso\".\nEjemplo: egreso mecato 13000 efectivo comida",
			errType.Got,
		)
	}

	var errMethod failure.ErrUnrecognizedPaymentMethod
	if errors.As(err, &errMethod) {
		return fmt.Sprintf(
			"No reconocí el medio de pago %q. Válidos: %s",
			errMethod.Got, strings.Join(transaction.ValidPaymentMethods(), ", "),
		)
	}

	var errCategory failure.ErrUnrecognizedCategory
	if errors.As(err, &errCategory) {
		if len(errCategory.Valid) == 0 {
			return fmt.Sprintf("No reconocí la categoría %q.", errCategory.Got)
		}
		return fmt.Sprintf(
			"No reconocí la categoría %q. Válidas: %s",
			errCategory.Got, strings.Join(errCategory.Valid, ", "),
		)
	}

	var errAmountInvalid failure.ErrInvalidAmount
	if errors.As(err, &errAmountInvalid) {
		return "El monto tiene que ser un número mayor a 0 (ej. 13000)."
	}

	var errFormat failure.ErrInvalidFormat
	if errors.As(err, &errFormat) {
		return fmt.Sprintf(
			"No entendí el mensaje. Formato: tipo concepto monto medio_de_pago [categoria].\nTipos válidos: %s.\nEjemplo: egreso mecato 13000 efectivo comida",
			strings.Join(transaction.ValidTypes(), ", "),
		)
	}

	var errAmount failure.ErrReceiptAmountNotFound
	if errors.As(err, &errAmount) {
		return "No pude leer un monto en la foto. Reenviala con el texto del movimiento como caption."
	}

	var errInterpret failure.ErrReceiptInterpretationFailed
	if errors.As(err, &errInterpret) {
		return "No se pudo interpretar la foto en este momento. Tocá el botón de nuevo para reintentar, o mandala con caption."
	}

	var errWrite failure.ErrWriteFailed
	if errors.As(err, &errWrite) {
		return "No se pudo guardar el movimiento. Intenta de nuevo en un momento."
	}

	return "Ocurrió un error inesperado registrando el movimiento."
}

// resumenReply arma el texto del resumen. Balance = Ingreso - Egreso
// ("cuánto tengo disponible", la convención estándar).
func resumenReply(summary summary.Monthly) string {
	return fmt.Sprintf(
		"Resumen de este mes:\nBalance: %.0f\nEgresos: %.0f (%d movimiento(s))\nIngresos: %.0f (%d movimiento(s))",
		summary.Balance(),
		summary.Egreso, summary.CountEgreso,
		summary.Ingreso, summary.CountIngreso,
	)
}

func ultimoReply(lt summary.Latest) string {
	lines := []string{
		"Último movimiento:",
		fmt.Sprintf("Tipo: %s", lt.Type),
		fmt.Sprintf("Concepto: %s", lt.Name),
		fmt.Sprintf("Monto: %.0f", lt.Amount),
	}
	if lt.PaymentMethod != "" {
		lines = append(lines, fmt.Sprintf("Medio de pago: %s", lt.PaymentMethod))
	}
	if lt.CategoryName != "" {
		lines = append(lines, fmt.Sprintf("Categoría: %s", lt.CategoryName))
	}
	if lt.Date != "" {
		lines = append(lines, fmt.Sprintf("Fecha: %s", lt.Date))
	}
	lines = append(lines, fmt.Sprintf("Notion: https://notion.so/%s", lt.PageID))
	return strings.Join(lines, "\n")
}
