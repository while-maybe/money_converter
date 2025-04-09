package ecbank

// ecBankError defines a sentinel error.
type ecBankError string

// ecBankError implements the error interface.
func (e ecBankError) Error() string {
	return string(e)
}
