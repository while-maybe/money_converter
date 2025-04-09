package ecbank

import (
	"fmt"
	"moneyconverter/money"
	"net/http"
)

const (
	ErrCallingServer = ecBankError("error calling server")
)

// Client can call the bank to retrieve exchange rates.
type Client struct {
}

// FetchExchangeRate fetches the ExchangeRate for the day and returns in.
func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	const euroxrefURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	resp, err := http.Get(euroxrefURL)
	defer resp.Body.Close()

	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrServerSide, err.Error())
	}

	return money.ExchangeRate{}, nil
}

const (
	clientErrorClass = 4
	serverErrorClass = 5
)

// checkStatusCode returns a different error depending on the returned status code.
func checkStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusOK:
		return nil
	case httpStatusClass(statusCode) == clientErrorClass:
		return fmt.Errorf("%w: %d", ErrClientSide, statusCode)
	case httpStatusClass(statusCode) == serverErrorClass:
		return fmt.Errorf("%w: %d", ErrServerSide, statusCode)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownStatusCode, statusCode)
	}
}

// httpStatusClass returns the class of an http status code.
func httpStatusClass(statusCode int) int {
	const httpErrorClassSize = 100
	return statusCode / httpErrorClassSize
}
