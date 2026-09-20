package usecase

import (
	"errors"
	"testing"
	"time"

	"finance-tracker/internal/domain"
)

func TestParseMessage_RealExamples(t *testing.T) {
	now := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		text         string
		wantType     domain.TransactionType
		wantName     string
		wantAmount   float64
		wantMethod   domain.PaymentMethod
		wantCategory string
	}{
		{"Egreso Gasolina 5000 tarjeta credito", domain.TransactionTypeExpense, "Gasolina", 5000, domain.PaymentMethodCreditCard, ""},
		{"Egreso Gasolina 5000 tarjeta debito", domain.TransactionTypeExpense, "Gasolina", 5000, domain.PaymentMethodDebitCard, ""},
		{"Egreso mecato 13000 efectivo", domain.TransactionTypeExpense, "mecato", 13000, domain.PaymentMethodCash, ""},
		{"egreso mecato 13000 tranferencia", domain.TransactionTypeExpense, "mecato", 13000, domain.PaymentMethodBank, ""},
		{"ingreso 3000000 banco", domain.TransactionTypeIncome, "Ingreso", 3000000, domain.PaymentMethodBank, ""},
		{"ingreso 12000 efectivo", domain.TransactionTypeIncome, "Ingreso", 12000, domain.PaymentMethodCash, ""},
		{"egreso mecato 13000 efectivo comida", domain.TransactionTypeExpense, "mecato", 13000, domain.PaymentMethodCash, "comida"},
	}

	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			got, err := ParseMessage(tc.text, now)
			if err != nil {
				t.Fatalf("ParseMessage(%q) error inesperado: %v", tc.text, err)
			}
			if got.Transaction.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", got.Transaction.Type, tc.wantType)
			}
			if got.Transaction.Name != tc.wantName {
				t.Errorf("Name = %q, want %q", got.Transaction.Name, tc.wantName)
			}
			if got.Transaction.Amount != tc.wantAmount {
				t.Errorf("Amount = %v, want %v", got.Transaction.Amount, tc.wantAmount)
			}
			if got.Transaction.PaymentMethod != tc.wantMethod {
				t.Errorf("PaymentMethod = %q, want %q", got.Transaction.PaymentMethod, tc.wantMethod)
			}
			if got.CategoryName != tc.wantCategory {
				t.Errorf("CategoryName = %q, want %q", got.CategoryName, tc.wantCategory)
			}
			if !got.Transaction.Date.Equal(now) {
				t.Errorf("Date = %v, want %v", got.Transaction.Date, now)
			}
		})
	}
}

func TestParseMessage_SeguridadSocialDosPalabras(t *testing.T) {
	now := time.Now()
	got, err := ParseMessage("egreso salud 50000 efectivo seguridad social", now)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if got.CategoryName != "seguridad social" {
		t.Errorf("CategoryName = %q, want %q", got.CategoryName, "seguridad social")
	}
}

func TestParseMessage_Errors(t *testing.T) {
	now := time.Now()

	t.Run("tipo no reconocido", func(t *testing.T) {
		_, err := ParseMessage("compra mecato 13000 efectivo", now)
		var target ErrUnrecognizedType
		if !errors.As(err, &target) {
			t.Fatalf("error = %v, want ErrUnrecognizedType", err)
		}
	})

	t.Run("medio de pago no reconocido", func(t *testing.T) {
		_, err := ParseMessage("egreso mecato 13000 bitcoin", now)
		var target ErrUnrecognizedPaymentMethod
		if !errors.As(err, &target) {
			t.Fatalf("error = %v, want ErrUnrecognizedPaymentMethod", err)
		}
	})

	t.Run("sin monto", func(t *testing.T) {
		_, err := ParseMessage("egreso mecato efectivo", now)
		var target ErrInvalidFormat
		if !errors.As(err, &target) {
			t.Fatalf("error = %v, want ErrInvalidFormat", err)
		}
	})

	t.Run("mensaje demasiado corto", func(t *testing.T) {
		_, err := ParseMessage("egreso 5000", now)
		var target ErrInvalidFormat
		if !errors.As(err, &target) {
			t.Fatalf("error = %v, want ErrInvalidFormat", err)
		}
	})
}
