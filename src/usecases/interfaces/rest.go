package interfaces

import "github.com/alelaca/stock-bot/src/entities"

type StockClient interface {
	GetStockQuote(stockCode string) (*entities.Stock, error)
}
