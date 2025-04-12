package ecbank

import (
	"errors"
	"fmt"
	"moneyconverter/money"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEuroCentralBank_FetchExchangeRate_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender>
		<gesmes:name>European Central Bank</gesmes:name>
	</gesmes:Sender>
	<Cube>
		<Cube time='2025-04-08'>
			<Cube currency='USD' rate='2.0000'/>
			<Cube currency='RON' rate='6.0000'/>
			<Cube currency='SEK' rate='10.9775'/>
			<Cube currency='CHF' rate='0.9349'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	got, err := ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))
	want := mustParseDecimal(t, "3")

	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}

	if money.Decimal(got) != want {
		t.Errorf("FetchExchangeRate got %v, want %v", money.Decimal(got), want)
	}
}

func TestEuroCentralBank_FetchExchangeRate_ErrCallingServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, ``)
	}))
	defer ts.Close()

	ecb := Client{url: "does not exist"}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "XYZ"), mustParseCurrency(t, "ABC"))

	if !errors.Is(err, ErrCallingServer) {
		t.Errorf("expected error %s, got %s", ErrCallingServer, err)
	}
}
func TestEuroCentralBank_FetchExchangeRate_ErrUnexpectedFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, ``)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "XYZ"), mustParseCurrency(t, "ABC"))

	if !errors.Is(err, ErrUnexpectedFormat) {
		t.Errorf("expected error %s, got %s", ErrUnexpectedFormat, err)
	}
}

func TestEuroCentralBank_FetchExchangeRate_ErrChangeRateNotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender>
		<gesmes:name>European Central Bank</gesmes:name>
	</gesmes:Sender>
	<Cube>
		<Cube time='2025-04-08'>
			<Cube currency='USD' rate='2.0000'/>
			<Cube currency='RON' rate='6.0000'/>
			<Cube currency='SEK' rate='10.9775'/>
			<Cube currency='CHF' rate='0.9349'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "XYZ"), mustParseCurrency(t, "ABC"))

	if !errors.Is(err, ErrChangeRateNotFound) {
		t.Errorf("expected error %s, got %s", ErrChangeRateNotFound, err)
	}
}

func TestEuroCentralBank_FetchExchangeRate_ErrClientSide(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))

	if !errors.Is(err, ErrClientSide) {
		t.Errorf("expected error %s, got %s", ErrClientSide, err)
	}
}
func TestEuroCentralBank_FetchExchangeRate_ErrServerSide(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))

	if !errors.Is(err, ErrServerSide) {
		t.Errorf("expected error %s, got %s", ErrServerSide, err)
	}
}

func TestEuroCentralBank_FetchExchangeRate_ErrUnknownStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSeeOther)
	}))
	defer ts.Close()

	ecb := Client{url: ts.URL}

	_, err := ecb.FetchExchangeRate(mustParseCurrency(t, "USD"), mustParseCurrency(t, "RON"))

	if !errors.Is(err, ErrUnknownStatusCode) {
		t.Errorf("expected error %s, got %s", ErrUnknownStatusCode, err)
	}
}

func mustParseCurrency(t *testing.T, code string) money.Currency {
	t.Helper()

	currency, err := money.ParseCurrency(code)
	if err != nil {
		t.Fatalf("cannot parse currency %s code", code)
	}

	return currency
}

func mustParseDecimal(t *testing.T, decimal string) money.Decimal {
	t.Helper()

	dec, err := money.ParseDecimal(decimal)
	if err != nil {
		t.Fatalf("cannot parse decimal %s", decimal)
	}

	return dec
}
