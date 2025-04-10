package ecbank

import (
	"moneyconverter/money"
	"testing"
)

func TestExchangeRate(t *testing.T) {
	tt := map[string]struct {
		envelope envelope
		source   string
		target   string
		want     money.ExchangeRate
		err      error
	}{
		"EUR to USD": {
			envelope: envelope{Rates: []currencyRate{{Currency: "USD", Rate: 1.5}}},
			source:   "EUR",
			target:   "USD",
			want:     mustParseExchangeRate(t, "1.5"),
			err:      nil,
		},
		"EUR to EUR": {
			envelope: envelope{Rates: []currencyRate{{Currency: "EUR", Rate: 1}}},
			source:   "EUR",
			target:   "EUR",
			want:     mustParseExchangeRate(t, "1"),
			err:      nil,
		},
		"CAD to EUR": {
			envelope: envelope{Rates: []currencyRate{{Currency: "CAD", Rate: 1.5}}},
			source:   "CAD",
			target:   "EUR",
			want:     mustParseExchangeRate(t, "0.6666666667"),
			err:      nil,
		},
		"CAD to USD": {
			envelope: envelope{Rates: []currencyRate{{Currency: "USD", Rate: 4}, {Currency: "CAD", Rate: 2}}},
			source:   "CAD",
			target:   "USD",
			want:     mustParseExchangeRate(t, "2"),
			err:      nil,
		},
		"CAD to XYZ": {
			envelope: envelope{Rates: []currencyRate{{Currency: "XYZ", Rate: 9}, {Currency: "CAD", Rate: 2}}},
			source:   "CAD",
			target:   "XYZ",
			want:     mustParseExchangeRate(t, "4.5"),
			err:      nil,
		},
	}
	// TODO add tc for errors: missing source, missing target, unable to parse currency

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := tc.envelope.exchangeRate(tc.source, tc.target)

			if err != tc.err {
				// fmt.Printf("\n\ngot: %s\nwant: %s\n\n", err.Error(), tc.err.Error())
				t.Errorf("unable to marshal: %s", err.Error())
			}

			if got != tc.want {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func mustParseExchangeRate(t *testing.T, rate string) money.ExchangeRate {
	t.Helper()

	excRate, err := money.ParseDecimal(rate)
	if err != nil {
		t.Fatalf("unable to parse exchange rate %s", rate)
	}
	return money.ExchangeRate(excRate)
}
