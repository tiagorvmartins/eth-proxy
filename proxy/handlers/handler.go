package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

func Handler(name string) func(queue string, msg amqp.Delivery, err error, ch *amqp.Channel, id string) {
	return func(queue string, msg amqp.Delivery, err error, ch *amqp.Channel, id string) {
		if err != nil {
			log.Err(err).Msgf("[%s] Error occurred in RMQ consumer", id)
		}

		providerUrl := os.Getenv(fmt.Sprintf("%s_URL", name))
		if providerUrl == "" {
			panic(fmt.Sprintf("[%s] URL for provider %s is not configured!", id, name))
		}

		log.Info().Msgf("[%s] Message received on Consumer on '%s' queue: %s", id, queue, string(msg.Body))

		req, err := http.NewRequest("POST", providerUrl, bytes.NewReader(msg.Body))
		if err != nil {
			log.Error().Msgf("[%s] Failed to mashal request on Consumer: %s", id, err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{Timeout: 10 * time.Second}
		res, err := client.Do(req)
		if err != nil {
			log.Error().Msgf("[%s] Failed to send request on Consumer: %s", id, err)
		}
		log.Info().Msgf("[%s] Response received with statusCode on Consumer: %d", id, res.StatusCode)

		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Msgf("[%s] Failed to read all body from response on Consumer: %s", id, err)
		}
		log.Info().Msgf("[%s] Response received with body on Consumer: %s", id, string(resBody))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		ch.PublishWithContext(ctx,
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				Body:          resBody,
				CorrelationId: msg.CorrelationId, // The same CorrelationId as received in the request
			},
		)
	}
}
