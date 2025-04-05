package money

import (
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
		// "prefix 0 as decimal digits": {...}
		// "multiple of 10": {...},#B
		// "invalid decimal part": {...},
		// "Not a number":
		// "empty string":
		// "too large"
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			got, err := ParseDecimal(tc.decimal)

			if err != tc.err {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
