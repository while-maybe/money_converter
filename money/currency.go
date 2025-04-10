package money

import (
	"regexp"
	"strings"
)

// Currency defines the code of a currency and its decimal precision.
type Currency struct {
	code      string
	precision byte
}

// ErrInvalidCurrencyCode is returned when the currency to parse is not a standard 3-letter code
const ErrInvalidCurrencyCode = Error("invalid currency code")

// ParseCurrency returns the currency associated to a name and may return ErrInvalidCurrencyCode
func ParseCurrency(code string) (Currency, error) {
	code = strings.ToUpper(code)

	validCode := regexp.MustCompile(`^[A-Z]{3}$`)

	if len(code) != 3 || !validCode.MatchString(code) {
		return Currency{}, ErrInvalidCurrencyCode
	}

	switch code {
	case "IRR":
		return Currency{code: code, precision: 0}, nil
	case "CNY", "VND":
		return Currency{code: code, precision: 1}, nil
	case "BHD", "IQD", "KWD", "LYD", "OMR", "TND":
		return Currency{code: code, precision: 3}, nil
	default:
		return Currency{code: code, precision: 2}, nil
	}
}

// String implements Stringer.
func (c Currency) String() string {
	return c.code
}

// ISOCode returns the 3 letter ISO code representing the currency name
// it is useful should the String() implementation ever need to change
func (c Currency) ISOCode() string {
	return c.code
}
