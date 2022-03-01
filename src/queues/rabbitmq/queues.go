package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/alelaca/stock-bot/src/entities/dtos"
	"github.com/alelaca/stock-bot/src/queues"
	"github.com/streadway/amqp"
)

type Handler struct {
	Connection   *amqp.Connection
	commandQueue queues.Config
	postsQueue   queues.Config
}

func InitializeRabbitMQHandler(connection *amqp.Connection) (*Handler, error) {
	commandQueue := queues.Config{
		Name: "command-queue",
	}
	postsQueue := queues.Config{
		Name: "posts-queue",
	}

	return &Handler{
		Connection:   connection,
		commandQueue: commandQueue,
		postsQueue:   postsQueue,
	}, nil
}

// Receives command messages from queue
func (h *Handler) PollCommandMessages(messageHandler queues.QueueMessageHandler) error {
	return h.receiveMessages(h.commandQueue, messageHandler)
}

// Sends command results to queue
func (h *Handler) NotifyCommandResult(post dtos.PostDTO) error {
	postBody, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return h.sendMessage(h.postsQueue, postBody)
}

// Receives messages from specified queue
// Handle any messages received with parameter function
func (h *Handler) receiveMessages(queue queues.Config, messageHandler queues.QueueMessageHandler) error {
	ch, err := h.Connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			queueMessage := queues.QueueMessage{
				Message: msg.Body,
			}

			err := messageHandler(queueMessage)
			if err != nil {
				fmt.Println(err.Error())
				msg.Reject(true)
				continue
			}

			msg.Ack(false)
		}
	}()

	<-forever

	return nil
}

// Sends a message to the specified queue
func (h *Handler) sendMessage(queue queues.Config, body []byte) error {
	ch, err := h.Connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		return err
	}

	return nil
}
