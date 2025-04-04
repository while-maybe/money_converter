package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Decimal can represent a floating points number with a fixed precision.
// example: 1.52 = 152 * 10^(-2) will be stored as {152, 2}
type Decimal struct {
	// subunits is the amount of sub units, Multiply it by the precision to get the real value
	subunits int64
	// number of "subunits" in a unit, expressed as a power of 10.
	precision byte
}

const maxDecimal = 1e12

func ParseDecimal(value string) (Decimal, error) {
	beforeSep, afterSep, _ := strings.Cut(value, ".")

	parsed, err := strconv.ParseInt(beforeSep+afterSep, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	if parsed > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(afterSep))

	return Decimal{subunits: parsed, precision: precision}, nil
}

const (
	// ErrInvalidDecimal is returned if the decimal is malformed.
	ErrInvalidDecimal = Error("unable to convert the decimal")
	// ErrTooLarge is returned if the quantity is too large - this would cause floating point precision errors
	ErrTooLarge = Error("quantity over 10^12 is too large")
)
