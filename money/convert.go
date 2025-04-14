package money

import (
	"fmt"
	"math"
)

// Convert applies the change rate to convert an amount to a target currency.
func Convert(amount Amount, to Currency, rates ratesFetcher) (Amount, error) {
	// fetch the exchange rate for the day

	r, err := rates.FetchExchangeRate(amount.currency, to)
	if err != nil {
		return Amount{}, fmt.Errorf("cannot get exchange rate: %w", err)
	}

	// convert to the target rate currency applying the fetched change rate.
	convertedValue := applyExchangeRate(amount, to, r)

	// validate the converted amount is within bounds
	if err := convertedValue.validate(); err != nil {
		return Amount{}, err
	}

	return convertedValue, nil
}

type ratesFetcher interface {
	// FetchExchangeRate fetches the ExchangeRate for the day and returns it.
	FetchExchangeRate(source, target Currency) (ExchangeRate, error)
}

// ExchangeRate represents a rate to convert from a currency to another.
type ExchangeRate Decimal

// pow10 is a quick implementation of how to raise 10 to a given power.
// It's optimized for small powers, and slow for unusually high powers.
func pow10(power byte) int64 {
	switch power {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	default:
		return int64(math.Pow(10, float64(power)))
	}
}

// applyExchangeRate returns a new Amount representing the input multiplied by the rate.
// The precision of the returned value is that of the target Currency.
// This function does not guarantee that the output amount is supported.
func applyExchangeRate(a Amount, target Currency, rate ExchangeRate) Amount {
	converted := multiply(a.quantity, rate)

	switch {
	case converted.precision > target.precision:
		converted.subunits /= pow10(converted.precision - target.precision)
	case converted.precision < target.precision:
		converted.subunits *= pow10(target.precision - converted.precision)
	}

	converted.precision = target.precision

	return Amount{
		quantity: converted,
		currency: target,
	}
}

// multiply a Decimal with an ExchangeRate and returns the product
func multiply(d Decimal, r ExchangeRate) Decimal {
	dec := Decimal{
		subunits:  d.subunits * r.subunits,
		precision: d.precision + r.precision,
	}

	// Clean the representation a bit. Remove trailing zeros.
	dec.simplify()

	return dec
}
