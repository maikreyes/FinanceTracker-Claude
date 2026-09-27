package services

import (
	"errors"
	"math"
	"testing"
	"time"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/receipt"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

func TestParseMessage_MontoNegativoOCero(t *testing.T) {
	for _, text := range []string{"egreso mecato -5000 efectivo", "egreso mecato 0 efectivo"} {
		_, err := ParseMessage(text, time.Now())
		var errAmount failure.ErrInvalidAmount
		if !errors.As(err, &errAmount) {
			t.Errorf("ParseMessage(%q) error = %v, want failure.ErrInvalidAmount", text, err)
		}
	}
}

func TestParseMessage_ConceptoNanNoEsMonto(t *testing.T) {
	parsed, err := ParseMessage("egreso nan 13000 efectivo", time.Now())
	if err != nil {
		t.Fatalf("ParseMessage error inesperado: %v", err)
	}
	if parsed.Transaction.Name != "nan" || parsed.Transaction.Amount != 13000 {
		t.Errorf("Name = %q, Amount = %v, want \"nan\" y 13000", parsed.Transaction.Name, parsed.Transaction.Amount)
	}
}

func TestRegistrar_RegisterFromFields_MontoInvalido(t *testing.T) {
	writer := &fakeExpenseWriter{id: "page-1"}
	r := NewServices(fakeCategoryResolver{}, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	for _, amount := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		_, err := r.RegisterFromFields(register.ConversationInput{Type: transaction.TypeExpense, Amount: amount, PaymentMethod: transaction.PaymentMethodCash})
		var errAmount failure.ErrInvalidAmount
		if !errors.As(err, &errAmount) {
			t.Errorf("RegisterFromFields(monto %v) error = %v, want failure.ErrInvalidAmount", amount, err)
		}
	}
	if writer.got.Name != "" {
		t.Error("un monto inválido no debería llegar a escribir en Notion")
	}
}

func TestRegistrar_RegisterFromPhotoAI_MontoInvalidoEsMontoNoEncontrado(t *testing.T) {
	interpreter := &fakeReceiptInterpreter{result: receipt.Interpreted{Amount: amountPtr(-100)}}
	r := NewServices(fakeCategoryResolver{}, &fakeExpenseWriter{}, &fakeReceiptUploader{}, interpreter)

	_, err := r.RegisterFromPhotoAI([]byte("x"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	var errAmount failure.ErrReceiptAmountNotFound
	if !errors.As(err, &errAmount) {
		t.Errorf("error = %v, want failure.ErrReceiptAmountNotFound", err)
	}
}

func TestRegistrar_UsaFechaDeBogota(t *testing.T) {
	// 2026-09-21 02:00 UTC sigue siendo 2026-09-20 21:00 en Bogotá.
	utc := time.Date(2026, 9, 21, 2, 0, 0, 0, time.UTC)
	writer := &fakeExpenseWriter{id: "page-1"}
	r := NewServices(fakeCategoryResolver{}, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})
	r.now = func() time.Time { return utc.In(bogota) }

	if _, err := r.Register("egreso mecato 13000 efectivo"); err != nil {
		t.Fatalf("Register error inesperado: %v", err)
	}
	if got := writer.got.Date.Format("2006-01-02"); got != "2026-09-20" {
		t.Errorf("Date = %s, want 2026-09-20 (día de Bogotá)", got)
	}
}
