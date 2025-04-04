package money

// Amount defines a quantity of money in a given currency.
type Amount struct {
	quantity Decimal
	currency Currency
}
