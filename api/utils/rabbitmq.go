package utils

import (
	"context"
	"encoding/json"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type RMQProducer struct {
	Queue            string
	ConnectionString string
	QueueCallback    string
}

func (x RMQProducer) OnError(err error, msg string) {
	if err != nil {
		log.Err(err).Msgf("Error occurred while publishing message on '%s' queue. Error message: %s", x.Queue, msg)
	}
}

func (x RMQProducer) PublishMessage(contentType string, body []byte, request_id string) map[string]any {
	conn, err := amqp.Dial(x.ConnectionString)
	x.OnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	x.OnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		x.QueueCallback, // name
		false,           // durable
		false,           // delete when unused
		true,            // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	x.OnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	x.OnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",                      // exchange
		os.Getenv("QUEUE_NAME"), // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: request_id,
			ReplyTo:       q.Name,
			Body:          body,
		})
	x.OnError(err, "Failed to publish a message")

	var responseMessage map[string]any
	for msg := range msgs {
		if msg.CorrelationId == request_id {
			if err := json.Unmarshal(msg.Body, &responseMessage); err != nil {
				log.Printf("Failed to unmarshal response message: %v", err)
				continue
			}
			break
		}
	}
	return responseMessage
}
