package queues

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alelaca/stock-bot/src/entities/dtos"
	"github.com/alelaca/stock-bot/src/queues"
	"github.com/alelaca/stock-bot/src/usecases/stock"
)

type Worker struct {
	QueuesHandler queues.QueuesHandler
	StockHandler  stock.Usecases
}

func InitializeWorker(queuesHandler queues.QueuesHandler, stockHandler stock.Usecases) *Worker {
	return &Worker{
		QueuesHandler: queuesHandler,
		StockHandler:  stockHandler,
	}
}

func (w *Worker) StartPollingCommandMessages() {
	for {
		w.QueuesHandler.PollCommandMessages(w.postsCommandHandler)
	}
}

// Handle a stock command message
// Checks for the stock code and gets the value from a stock client
// Sends the message to posts queue
func (w *Worker) postsCommandHandler(msg queues.QueueMessage) error {
	postDTO := dtos.PostDTO{}
	if err := json.Unmarshal(msg.Message, &postDTO); err != nil {
		return err
	}

	// rabbitmq topic only sends "/stock=" to this queue
	stockCode := strings.TrimPrefix(postDTO.Message, "/stock=")

	stock, err := w.StockHandler.GetStock(stockCode)
	if err != nil {
		return err
	}

	newPost := dtos.PostDTO{
		Sender: "stock-bot",
		Room:   postDTO.Room,
	}
	if stock == nil {
		newPost.Message = fmt.Sprintf("The stock '%s' doesn't exist", stockCode)
	} else {
		newPost.Message = fmt.Sprintf("%s quote is $%.2f per share", stock.Code, stock.Value)
	}

	return w.QueuesHandler.NotifyCommandResult(newPost)
}
