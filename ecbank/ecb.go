package ecbank

import (
	"bytes"
	"fmt"
	"io"
	"moneyconverter/money"
	"net/http"
)

const (
	ErrCallingServer      = ecBankError("error calling server")
	ErrUnexpectedFormat   = ecBankError("unexpected response format")
	ErrChangeRateNotFound = ecBankError("couldn't find the exchange rate")
	ErrClientSide         = ecBankError("client side error when contacting ECB")
	ErrServerSide         = ecBankError("server side error when contacting ECB")
	ErrUnknownStatusCode  = ecBankError("unknown status code contacting ECB")
)

// Client can call the bank to retrieve exchange rates.
type Client struct {
	url string
}

// FetchExchangeRate fetches the ExchangeRate for the day and returns in.
func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	const euroxrefURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	if c.url == "" {
		c.url = euroxrefURL
	}

	resp, err := http.Get(c.url)
	if err != nil {
		return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrCallingServer, err.Error())
	}
	defer resp.Body.Close()

	if err = checkStatusCode(resp.StatusCode); err != nil {
		return money.ExchangeRate{}, err
	}

	ClearInvalidCache()

	// copy the stream into a buffer and attempt to create a new cache
	dataBuffer := bytes.NewBuffer(make([]byte, 0, 2048))
	err = writeToCache(dataBuffer, resp.Body)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	rate, err := readRateFromResponse(source.ISOCode(), target.ISOCode(), dataBuffer)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	return rate, nil
}

// writeToCache creates a buffer and attempts to write to cache
func writeToCache(buf *bytes.Buffer, data io.ReadCloser) error {
	cache := newCache()
	err := cache.writeCache(io.TeeReader(data, buf))
	if err != nil {
		return fmt.Errorf("Couldn't write to cache: %w", err)
	}
	return nil
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
