package ecbank

import (
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
			<Cube currency='JPY' rate='160.65'/>
			<Cube currency='BGN' rate='1.9558'/>
			<Cube currency='CZK' rate='25.139'/>
			<Cube currency='DKK' rate='7.4648'/>
			<Cube currency='GBP' rate='0.85644'/>
			<Cube currency='HUF' rate='407.00'/>
			<Cube currency='PLN' rate='4.2690'/>
			<Cube currency='RON' rate='6.0000'/>
			<Cube currency='SEK' rate='10.9775'/>
			<Cube currency='CHF' rate='0.9349'/>
			<Cube currency='ISK' rate='145.10'/>
			<Cube currency='NOK' rate='11.9505'/>
			<Cube currency='TRY' rate='41.6215'/>
			<Cube currency='AUD' rate='1.8073'/>
			<Cube currency='BRL' rate='6.4211'/>
			<Cube currency='CAD' rate='1.5512'/>
			<Cube currency='CNY' rate='8.0359'/>
			<Cube currency='HKD' rate='8.5059'/>
			<Cube currency='IDR' rate='18522.86'/>
			<Cube currency='ILS' rate='4.1168'/>
			<Cube currency='INR' rate='94.4125'/>
			<Cube currency='KRW' rate='1618.79'/>
			<Cube currency='MXN' rate='22.4776'/>
			<Cube currency='MYR' rate='4.9176'/>
			<Cube currency='NZD' rate='1.9534'/>
			<Cube currency='PHP' rate='62.776'/>
			<Cube currency='SGD' rate='1.4766'/>
			<Cube currency='THB' rate='37.986'/>
			<Cube currency='ZAR' rate='21.1803'/>
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
