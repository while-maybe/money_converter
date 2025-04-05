package money

import (
	"errors"
	"testing"
)

func TestParseDecimal(t *testing.T) {
	tt := map[string]struct {
		decimal  string
		expected Decimal
		err      error
	}{
		"2 decimal digits": {
			decimal: "1.52",
			expected: Decimal{
				subunits:  152,
				precision: 2,
			},
			err: nil,
		},
		"no decimal digits": {
			decimal:  "152",
			expected: Decimal{152, 0},
			err:      nil,
		},
		"suffix 0 as decimal digits": {
			decimal:  "8.200",
			expected: Decimal{82, 1},
			err:      nil,
		},
		"prefix 0 as decimal digits": {
			decimal:  "9.006",
			expected: Decimal{9006, 3},
			err:      nil,
		},
		"multiple of 10": {
			decimal:  "500",
			expected: Decimal{500, 0},
			err:      nil,
		},
		"invalid decimal part": {
			decimal: "500.mistake",
			// expected: n/a
			err: ErrInvalidDecimal,
		},
		"Not a number": {
			decimal: "notnumeric",
			// expected: n/a
			err: ErrInvalidDecimal,
		},
		"empty string": {
			decimal: "",
			// expected: n/a
			err: ErrInvalidDecimal,
		},
		"too large": {
			decimal: "1234567890123",
			err:     ErrTooLarge,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			got, err := ParseDecimal(tc.decimal)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
