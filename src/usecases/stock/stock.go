package stock

import (
	"github.com/alelaca/stock-bot/src/entities"
	"github.com/alelaca/stock-bot/src/usecases/interfaces"
)

type Usecases interface {
	GetStock(stockCode string) (*entities.Stock, error)
}

type Handler struct {
	StockClient interfaces.StockClient
}

// Get the stock for a particular stock code
func (h *Handler) GetStock(stockCode string) (*entities.Stock, error) {
	return h.StockClient.GetStockQuote(stockCode)
}
