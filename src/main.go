package main

import (
	"fmt"
	"net/http"

	"github.com/alelaca/stock-bot/src/adapters/queues"
	"github.com/alelaca/stock-bot/src/queues/rabbitmq"
	"github.com/alelaca/stock-bot/src/rest/stooq"
	"github.com/alelaca/stock-bot/src/usecases/stock"
	"github.com/streadway/amqp"
)

func main() {
	stooq := stooq.InitializeClient(http.DefaultClient)

	stockHandler := stock.Handler{
		StockClient: stooq,
	}

	connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(fmt.Sprintf("error initializing rabbitmq connection, log: %s", err.Error()))
	}

	queuesHandler, err := rabbitmq.InitializeRabbitMQHandler(connection)

	postsCommandWorker := queues.InitializeWorker(queuesHandler, &stockHandler)
	postsCommandWorker.StartPollingCommandMessages()
}
