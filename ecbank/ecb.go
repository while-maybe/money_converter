package ecbank

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"moneyconverter/money"
	"net/http"
	"net/url"
	"time"
)

const (
	ErrCallingServer      = ecBankError("error calling server")
	ErrTimeout            = ecBankError("timed out when waiting for response")
	ErrUnexpectedFormat   = ecBankError("unexpected response format")
	ErrChangeRateNotFound = ecBankError("couldn't find the exchange rate")
	ErrClientSide         = ecBankError("client side error when contacting ECB")
	ErrServerSide         = ecBankError("server side error when contacting ECB")
	ErrUnknownStatusCode  = ecBankError("unknown status code contacting ECB")
)

// Client can call the bank to retrieve exchange rates.
type Client struct {
	client *http.Client
}

// NewClient builds a client that can fetch exchange rates within a given timeout.
func NewClient(timeout time.Duration) Client {
	return Client{
		client: &http.Client{Timeout: timeout},
	}
}

// FetchExchangeRate fetches the ExchangeRate for the day and returns in.
func (c Client) FetchExchangeRate(source, target money.Currency) (money.ExchangeRate, error) {
	dataBuffer := bytes.NewBuffer(make([]byte, 0, 4096))
	err := readFromCache(dataBuffer)

	if err != nil {
		fmt.Print("[API CALL] ")
		const path = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

		resp, err := c.client.Get(path)

		if err != nil {
			var urlError *url.Error
			if ok := errors.As(err, &urlError); ok && urlError.Timeout() {
				return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrTimeout, err.Error())
			}

			return money.ExchangeRate{}, fmt.Errorf("%w: %s", ErrCallingServer, err.Error())
		}
		defer resp.Body.Close()

		if err = checkStatusCode(resp.StatusCode); err != nil {
			return money.ExchangeRate{}, err
		}

		err = writeToCache(dataBuffer, resp.Body)
		if err != nil {
			return money.ExchangeRate{}, err
		}
	}

	rate, err := readRateFromResponse(source.ISOCode(), target.ISOCode(), dataBuffer)
	if err != nil {
		return money.ExchangeRate{}, err
	}

	return rate, nil
}

// writeToCache creates a buffer and attempts to write to file cache
func writeToCache(buf *bytes.Buffer, data io.ReadCloser) error {
	cache := newCache()
	err := cache.writeCache(io.TeeReader(data, buf))
	if err != nil {
		return fmt.Errorf("couldn't write to cache: %w", err)
	}
	return nil
}

// readFromCache creates a buffer and attempts to read from file cache
func readFromCache(buf *bytes.Buffer) error {
	cache := newCache()
	err := cache.readCache(buf)
	if err != nil {
		return fmt.Errorf("couldn't read from cache: %w", err)
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
