package money_test

import (
	"moneyconverter/money"
	"reflect"
	"testing"
)

// stubRate is a very simple stub for the exchangeRates.
type stubRate struct {
	rate string
	err  error
}

// FetchExchangeRate implements the interface ratesFetcher with the same signature but fields are unused for tests purposes.
func (m stubRate) FetchExchangeRate(_, _ money.Currency) (money.ExchangeRate, error) {
	rate, _ := money.ParseDecimal(m.rate)
	return money.ExchangeRate(rate), m.err
}

func TestConvert(t *testing.T) {
	tt := map[string]struct {
		amount   money.Amount
		to       money.Currency
		validate func(t *testing.T, got money.Amount, err error)
	}{
		"34.98USD to EUR": {
			amount: mustParseAmount(t, "34.98", "USD"),
			to:     mustParseCurrency(t, "EUR"),
			validate: func(t *testing.T, got money.Amount, err error) {
				if err != nil {
					t.Errorf("expected no error, got %s", err.Error())
				}

				expected := mustParseAmount(t, "69.96", "EUR")
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected %v, got %v", expected, got)
				}
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			stub := stubRate{rate: "2", err: nil}

			got, err := money.Convert(tc.amount, tc.to, stub)
			tc.validate(t, got, err)
		})
	}
}
