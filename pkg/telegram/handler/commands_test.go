package handler

import (
	"testing"

	"finance-tracker/pkg/transaction/model/transaction"
)

func TestIsAllowed(t *testing.T) {
	allowed := []int64{111, 222}

	if !isAllowed(111, allowed) {
		t.Error("111 debería estar autorizado")
	}
	if isAllowed(333, allowed) {
		t.Error("333 no debería estar autorizado")
	}
	if isAllowed(111, nil) {
		t.Error("con whitelist vacía, nada debería estar autorizado (fail closed)")
	}
}

func TestTipoCallbackData_RoundTrip(t *testing.T) {
	for _, want := range []transaction.Type{transaction.TypeIncome, transaction.TypeExpense} {
		data := tipoCallbackData(want)
		got, ok := parseTipoCallback(data)
		if !ok {
			t.Fatalf("parseTipoCallback(%q) = ok=false, want true", data)
		}
		if got != want {
			t.Errorf("parseTipoCallback(%q) = %q, want %q", data, got, want)
		}
	}
}

func TestParseTipoCallback_Desconocido(t *testing.T) {
	for _, data := range []string{"tipo:algo_raro", ""} {
		if _, ok := parseTipoCallback(data); ok {
			t.Errorf("parseTipoCallback(%q) debería devolver ok=false", data)
		}
	}
}

func TestNormalizeCommand(t *testing.T) {
	cases := map[string]string{
		"/Egreso":               "/egreso",
		"  /ingreso  ":          "/ingreso",
		"/ayuda@FinanceMoikBot": "/ayuda",
		"":                      "",
	}
	for in, want := range cases {
		if got := normalizeCommand(in); got != want {
			t.Errorf("normalizeCommand(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseStartCommand(t *testing.T) {
	cases := []struct {
		text string
		want transaction.Type
		ok   bool
	}{
		{"/egreso", transaction.TypeExpense, true},
		{"/ingreso", transaction.TypeIncome, true},
		{"/Egreso", transaction.TypeExpense, true},
		{"/egreso@FinanceMoikBot", transaction.TypeExpense, true},
		{"  /ingreso  ", transaction.TypeIncome, true},
		{"/resumen", "", false},
		{"egreso", "", false}, // sin "/" no es comando, es el formato de una línea
		{"", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			got, ok := parseStartCommand(tc.text)
			if ok != tc.ok || got != tc.want {
				t.Errorf("parseStartCommand(%q) = (%q, %v), want (%q, %v)", tc.text, got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestMatchesCommand(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"/ayuda", true},
		{"/Ayuda", true},
		{"  /ayuda  ", true},
		{"/ayuda@FinanceMoikBot", true},
		{"/ultimo", false},
		{"ayuda", false}, // sin "/" no es comando
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			if got := matchesCommand(tc.text, "/ayuda"); got != tc.want {
				t.Errorf("matchesCommand(%q, \"/ayuda\") = %v, want %v", tc.text, got, tc.want)
			}
		})
	}
}
