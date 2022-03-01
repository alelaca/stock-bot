package stooq

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/alelaca/stock-bot/src/rest"
)

func TestGetStockQuote_OK(t *testing.T) {
	expectedBody := `Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.US,2022-02-28,22:00:07,163.06,165.42,162.43,165.12,69201975`

	body := []byte(expectedBody)
	r := ioutil.NopCloser(bytes.NewReader(body))
	expectedResponse := http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}

	httpClient := rest.HTTPClientMock{
		Res: &expectedResponse,
	}

	stooqClient := Client{
		RestClient: httpClient,
	}

	stock, err := stooqClient.GetStockQuote("aapl.us")

	if err != nil {
		t.Errorf("Test failed with error %s", err.Error())
	}

	if stock.Code != "AAPL.US" {
		t.Errorf("Test failed. Stock code expected %s, got %s", "AAPL.US", stock.Code)
	}

	if stock.Value != 165.12 {
		t.Errorf("Test failed. Stock code expected %f, got %f", 165.12, stock.Value)
	}
}

func TestGetStockQuote_NoResult(t *testing.T) {
	expectedBody := `Symbol,Date,Time,Open,High,Low,Close,Volume`

	body := []byte(expectedBody)
	r := ioutil.NopCloser(bytes.NewReader(body))
	expectedResponse := http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}

	httpClient := rest.HTTPClientMock{
		Res: &expectedResponse,
	}

	stooqClient := Client{
		RestClient: httpClient,
	}

	stock, err := stooqClient.GetStockQuote("unknown")

	if err != nil {
		t.Errorf("Test failed with error %s", err.Error())
	}

	if stock != nil {
		t.Errorf("Test failed. Expected nil stock")
	}
}

func TestGetStockQuote_StooqResultInND(t *testing.T) {
	expectedBody := `Symbol,Date,Time,Open,High,Low,Close,Volume
APL.US,N/D,N/D,N/D,N/D,N/D,N/D,N/D`

	body := []byte(expectedBody)
	r := ioutil.NopCloser(bytes.NewReader(body))
	expectedResponse := http.Response{
		StatusCode: http.StatusOK,
		Body:       r,
	}

	httpClient := rest.HTTPClientMock{
		Res: &expectedResponse,
	}

	stooqClient := Client{
		RestClient: httpClient,
	}

	stock, err := stooqClient.GetStockQuote("unknown")

	if err != nil {
		t.Errorf("Test failed with error %s", err.Error())
	}

	if stock != nil {
		t.Errorf("Test failed. Expected nil stock")
	}
}
