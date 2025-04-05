package money

import (
	"errors"
	"testing"
)

func TestParseCurrency(t *testing.T) {
	tt := map[string]struct {
		currencyCode string
		expected     Currency
	}{
		"lowercase hundredth EUR": {
			currencyCode: "eur",
			expected:     Currency{code: "EUR", precision: 2},
		},
		"hundredth EUR":  {"EUR", Currency{code: "EUR", precision: 2}},
		"thousandth BHD": {"BHD", Currency{"BHD", 3}},
		"tenth VND":      {"VND", Currency{"VND", 1}},
		"integer IRR":    {"IRR", Currency{"IRR", 0}},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := ParseCurrency(tc.currencyCode)

			if err != nil {
				t.Errorf("expected no errors, got %s", err.Error())
			}

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestParseCurrency_UnknownCurrency(t *testing.T) {
	tt := map[string]struct {
		currencyCode string
		err          error
	}{
		"4 letter code": {
			currencyCode: "XXXX",
			err:          ErrInvalidCurrencyCode,
		},
		"unkown 2 letter code":     {"XX", ErrInvalidCurrencyCode},
		"unicode letter code":      {"\u00C95X", ErrInvalidCurrencyCode},
		"mixed char letter code":   {"X5X", ErrInvalidCurrencyCode},
		"numbers only letter code": {"777", ErrInvalidCurrencyCode},
		"symbol letter code":       {"A!A", ErrInvalidCurrencyCode},
		"empty letter code":        {"", ErrInvalidCurrencyCode},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := ParseCurrency(tc.currencyCode)

			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}
		})
	}
}
