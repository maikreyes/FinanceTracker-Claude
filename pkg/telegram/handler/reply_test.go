package handler

import (
	"errors"
	"strings"
	"testing"
	"time"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/summary"
	"finance-tracker/pkg/transaction/model/transaction"
)

func TestSuccessReply(t *testing.T) {
	result := register.Result{
		PageID: "page-123",
		Transaction: transaction.Transaction{
			Name:          "mecato",
			Amount:        13000,
			Type:          transaction.TypeExpense,
			Date:          time.Now(),
			PaymentMethod: transaction.PaymentMethodCash,
		},
		CategoryName: "comida",
	}

	reply := successReply(result)

	for _, want := range []string{"Egreso", "mecato", "13000", "Cash", "comida", "page-123"} {
		if !strings.Contains(reply, want) {
			t.Errorf("successReply() = %q, no contiene %q", reply, want)
		}
	}
}

func TestSuccessReply_SinCategoria(t *testing.T) {
	result := register.Result{
		PageID: "page-456",
		Transaction: transaction.Transaction{
			Name:          "Ingreso",
			Amount:        3000000,
			Type:          transaction.TypeIncome,
			PaymentMethod: transaction.PaymentMethodBank,
		},
	}

	reply := successReply(result)

	if strings.Contains(reply, "Categoría:") {
		t.Errorf("successReply() = %q, no debería mencionar categoría", reply)
	}
}

func TestErrorReply(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want []string
	}{
		{
			"tipo no reconocido",
			failure.ErrUnrecognizedType{Got: "compra"},
			[]string{"compra", "egreso", "ingreso"},
		},
		{
			"medio de pago no reconocido",
			failure.ErrUnrecognizedPaymentMethod{Got: "bitcoin"},
			[]string{"bitcoin", "efectivo"},
		},
		{
			"categoria no reconocida con lista",
			failure.ErrUnrecognizedCategory{Got: "inventada", Valid: []string{"Comida", "Ahorro"}},
			[]string{"inventada", "Comida", "Ahorro"},
		},
		{
			"categoria no reconocida sin lista",
			failure.ErrUnrecognizedCategory{Got: "inventada"},
			[]string{"inventada"},
		},
		{
			"formato invalido",
			failure.ErrInvalidFormat{Reason: "no se encontró un monto numérico"},
			[]string{"egreso", "ingreso"},
		},
		{
			"falla de escritura no expone el error real",
			failure.ErrWriteFailed{Err: errors.New("notion: unexpected status 500: secreto interno")},
			[]string{"Intenta de nuevo"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reply := errorReply(tc.err)
			for _, want := range tc.want {
				if !strings.Contains(reply, want) {
					t.Errorf("errorReply(%v) = %q, no contiene %q", tc.err, reply, want)
				}
			}
		})
	}

	t.Run("falla de escritura no filtra detalles internos", func(t *testing.T) {
		reply := errorReply(failure.ErrWriteFailed{Err: errors.New("secreto interno")})
		if strings.Contains(reply, "secreto interno") {
			t.Errorf("errorReply() = %q, no debería exponer el error interno", reply)
		}
	})
}

func TestResumenReply(t *testing.T) {
	summary := summary.Monthly{Egreso: 13000, Ingreso: 3012000, CountEgreso: 2, CountIngreso: 1}
	reply := resumenReply(summary)
	for _, want := range []string{"13000", "3012000", "2999000", "2", "1"} {
		if !strings.Contains(reply, want) {
			t.Errorf("resumenReply() = %q, no contiene %q", reply, want)
		}
	}
}

func TestResumenReply_SinMovimientos(t *testing.T) {
	reply := resumenReply(summary.Monthly{})
	if !strings.Contains(reply, "Balance: 0") {
		t.Errorf("resumenReply(vacío) = %q, debería mostrar Balance: 0", reply)
	}
}

func TestUltimoReply(t *testing.T) {
	lt := summary.Latest{
		PageID:        "page-789",
		Name:          "mecato",
		Amount:        13000,
		Type:          transaction.TypeExpense,
		PaymentMethod: transaction.PaymentMethodCash,
		CategoryName:  "comida",
		Date:          "2026-09-20",
	}

	reply := ultimoReply(lt)

	for _, want := range []string{"Egreso", "mecato", "13000", "Cash", "comida", "2026-09-20", "page-789"} {
		if !strings.Contains(reply, want) {
			t.Errorf("ultimoReply() = %q, no contiene %q", reply, want)
		}
	}
}

func TestUltimoReply_SinCategoriaNiFecha(t *testing.T) {
	lt := summary.Latest{PageID: "page-1", Name: "Ingreso", Amount: 3000000, Type: transaction.TypeIncome}

	reply := ultimoReply(lt)

	if strings.Contains(reply, "Categoría:") {
		t.Errorf("ultimoReply() = %q, no debería mencionar categoría", reply)
	}
	if strings.Contains(reply, "Fecha:") {
		t.Errorf("ultimoReply() = %q, no debería mencionar fecha", reply)
	}
}
