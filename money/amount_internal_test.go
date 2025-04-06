package money

import (
	"errors"
	"testing"
)

func TestNewAmount(t *testing.T) {
	tt := map[string]struct {
		quantity Decimal
		currency Currency
		want     Amount
		err      error
	}{
		"1.20 EUR": {
			quantity: Decimal{120, 2},
			currency: Currency{"EUR", 2},
			want: Amount{
				quantity: Decimal{120, 2},
				currency: Currency{"EUR", 2},
			},
		},
		"5.500 EUR": {
			quantity: Decimal{5500, 3},
			currency: Currency{"EUR", 2},
			err:      ErrTooPrecise,
		},
		"5.500 BHD": {
			quantity: Decimal{5500, 3},
			currency: Currency{"BHD", 3},
			want: Amount{
				quantity: Decimal{5500, 3},
				currency: Currency{"BHD", 3},
			},
		},
		"8.3 IRR": {
			quantity: Decimal{83, 1},
			currency: Currency{"IRR", 0},
			err:      ErrTooPrecise,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := NewAmount(tc.quantity, tc.currency)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
