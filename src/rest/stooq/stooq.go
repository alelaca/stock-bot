package stooq

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/alelaca/stock-bot/src/apperrors"
	"github.com/alelaca/stock-bot/src/entities"
	"github.com/alelaca/stock-bot/src/rest"
)

const (
	URLStooq               = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
	FileStockCodePostion   = 0
	FileStockValuePosition = 6
)

type Client struct {
	RestClient rest.HTTPClient
}

func InitializeClient(restClient rest.HTTPClient) *Client {
	return &Client{
		RestClient: restClient,
	}
}

// Gets stock data from stooq for a particular stock code
func (c *Client) GetStockQuote(stockCode string) (*entities.Stock, error) {
	url := fmt.Sprintf(URLStooq, stockCode)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, apperrors.CreateInternalServerError(fmt.Sprintf("error creating new request to %s", url), err.Error())
	}

	res, err := c.RestClient.Do(req)
	if err != nil {
		return nil, err
	}

	return parseResponseBody(res.Body)
}

// Given a csv response body, it parses it to get the stock results
func parseResponseBody(body io.Reader) (*entities.Stock, error) {
	lines, err := csv.NewReader(body).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) < 2 {
		return nil, nil // not found
	}

	stockValueStr := lines[1][FileStockValuePosition]

	if stockValueStr == "N/D" {
		return nil, nil // not found
	}

	stockValue, err := strconv.ParseFloat(stockValueStr, 64)
	if err != nil {
		return nil, err
	}

	return &entities.Stock{
		Code:  lines[1][FileStockCodePostion],
		Value: stockValue,
	}, nil
}
