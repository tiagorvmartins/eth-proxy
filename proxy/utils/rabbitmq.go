package utils

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/rs/zerolog/log"
)

type RMQConsumer struct {
	Id               string
	Queue            string
	ConnectionString string
	MsgHandler       func(queue string, msg amqp.Delivery, err error, ch *amqp.Channel, id string)
}

func (x RMQConsumer) OnError(err error, msg string) {
	if err != nil {
		x.MsgHandler(x.Queue, amqp.Delivery{}, err, nil, x.Id)
	}
}

func (x RMQConsumer) Consume() {
	conn, err := amqp.Dial(x.ConnectionString)
	x.OnError(err, fmt.Sprintf("Failed to connect to RabbitMQ on Consumer with Id: %s", x.Id))
	defer conn.Close()

	ch, err := conn.Channel()
	x.OnError(err, fmt.Sprintf("Failed to open a channel on Consumer with Id: %s", x.Id))
	defer ch.Close()

	q, err := ch.QueueDeclare(
		x.Queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	x.OnError(err, fmt.Sprintf("Failed to declare a queue on Consumer with Id: %s", x.Id))

	// err = ch.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// x.OnError(err, fmt.Sprintf("Failed to set QoS on Consumer with Id: %s", x.Id))

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	x.OnError(err, fmt.Sprintf("Failed to register a consumer on Consumer with Id: %s", x.Id))

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			x.MsgHandler(x.Queue, d, nil, ch, x.Id)
		}
	}()
	log.Info().Msgf("Started listening on Consumer with Id '%s' for messages on '%s' queue", x.Id, x.Queue)
	<-forever
}
