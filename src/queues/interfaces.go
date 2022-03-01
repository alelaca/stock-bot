package queues

import "github.com/alelaca/stock-bot/src/entities/dtos"

type QueueMessageHandler = func(message QueueMessage) error

type QueueMessage struct {
	Message []byte
}

type Config struct {
	Name string
}

type QueuesHandler interface {
	PollCommandMessages(messageHandler QueueMessageHandler) error
	NotifyCommandResult(post dtos.PostDTO) error
}
