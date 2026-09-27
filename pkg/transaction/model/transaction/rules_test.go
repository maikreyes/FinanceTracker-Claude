package transaction

import (
	"math"
	"testing"
)

func TestValidAmount(t *testing.T) {
	cases := []struct {
		amount float64
		want   bool
	}{
		{13000, true},
		{0.5, true},
		{0, false},
		{-5, false},
		{math.NaN(), false},
		{math.Inf(1), false},
		{math.Inf(-1), false},
	}
	for _, tc := range cases {
		if got := ValidAmount(tc.amount); got != tc.want {
			t.Errorf("ValidAmount(%v) = %v, want %v", tc.amount, got, tc.want)
		}
	}
}
