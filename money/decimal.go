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

const (
	// maxDecimal value is a thousand billion, using the shortscale 10^12
	maxDecimal = 1e12
	// ErrInvalidDecimal is returned if the decimal is malformed.
	ErrInvalidDecimal = Error("unable to convert the decimal")
	// ErrTooLarge is returned if the quantity is too large - this would cause floating point precision errors
	ErrTooLarge = Error("quantity over 10^12 is too large")
)

// String implements Stringer and returns the decimal formatted as digits and optionally a decimal point followed by digits
// mind the simplify method at the bottom which already has a pointer receiver
func (d *Decimal) String() string {
	if d.precision == 0 {
		return fmt.Sprintf("%d", d.subunits)
	}

	centsPerUnit := pow10(d.precision)
	frac := d.subunits % centsPerUnit
	integer := d.subunits / centsPerUnit

	decimalFormat := "%d.%0" + strconv.Itoa(int(d.precision)) + "d"
	return fmt.Sprintf(decimalFormat, integer, frac)
}

// ParseDecimal convert a string into its decimal representation.
// It assumes there is up to a decimal separator, and that the separator is '.' (full stop).
func ParseDecimal(value string) (Decimal, error) {
	beforeSep, afterSep, _ := strings.Cut(value, ".")

	parsed, err := strconv.ParseInt(beforeSep+afterSep, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	// fmt.Printf("\nvalue:%s\nbefore and after: %s %s\nparsed:%d\n\n", value, beforeSep, afterSep, parsed)
	if parsed > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	precision := byte(len(afterSep))
	result := Decimal{subunits: parsed, precision: precision}
	result.simplify()

	return result, nil
}

// simplify removes the trailing 0s after the . and decreases precision from a Decimal
func (d *Decimal) simplify() {
	// using %10 returns the last digit in base 10 of a number.
	// If the precision is positive, that digit belongs to the right side of the decimal separator
	for d.subunits%10 == 0 && d.precision > 0 {
		d.subunits /= 10
		d.precision--
	}
}
