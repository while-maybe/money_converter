package ecbank

import (
	"moneyconverter/money"
)

type envelope struct {
	Rates []currencyRate `xml:"Cube>Cube>Cube"`
}

type currencyRate struct {
	Currency string             `xml:"currency,attr"`
	Rate     money.ExchangeRate `xml:"rate,attr"`
}

const baseCurrencyCode = "EUR"

// exchangeRates builds a map of all the supported exchange rates.
func (e envelope) exchangeRates() map[string]money.ExchangeRate {
	rates := make(map[string]money.ExchangeRate, len(e.Rates)+1)

	for _, c := range e.Rates {
		rates[c.Currency] = c.Rate
	}

	// represents EUR to EUR rate
	one, _ := money.ParseDecimal("1")

	rates[baseCurrencyCode] = money.ExchangeRate(one)

	return rates
}
