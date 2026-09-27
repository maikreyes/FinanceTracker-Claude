package transaction

import "testing"

func TestParsePaymentMethod(t *testing.T) {
	cases := []struct {
		text string
		want PaymentMethod
		ok   bool
	}{
		{"efectivo", PaymentMethodCash, true},
		{"Efectivo", PaymentMethodCash, true},
		{"tarjeta credito", PaymentMethodCreditCard, true},
		{"tarjeta de credito", PaymentMethodCreditCard, true},
		{"tarjeta debito", PaymentMethodDebitCard, true},
		{"banco", PaymentMethodBank, true},
		{"transferencia", PaymentMethodBank, true},
		{"tranferencia", PaymentMethodBank, true},
		{"bitcoin", "", false},
		{"efectivo extra", "", false}, // sobra texto — no matchea exacto
		{"", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.text, func(t *testing.T) {
			got, ok := ParsePaymentMethod(tc.text)
			if ok != tc.ok {
				t.Fatalf("ParsePaymentMethod(%q) ok = %v, want %v", tc.text, ok, tc.ok)
			}
			if got != tc.want {
				t.Errorf("ParsePaymentMethod(%q) = %q, want %q", tc.text, got, tc.want)
			}
		})
	}
}
