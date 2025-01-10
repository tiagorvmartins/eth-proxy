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
	"encoding/json"
)

type Request struct {
	Path    string	      `json:"Path" binding:"omitempty"`
	JsonRpc string        `json:"jsonrpc" binding:"required"`
	Method  string        `json:"method" binding:"required"`
	Params  []interface{} `json:"params"`
	Id      *int          `json:"id" binding:"required"`
}

type ProviderRequest struct {
	JsonRpc string        `json:"jsonrpc" binding:"required"`
	Method  string        `json:"method" binding:"required"`
	Params  []interface{} `json:"params"`
	Id      *int          `json:"id" binding:"required"`
}

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


		// Parse the JSON into a struct
		var data Request
		err = json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Error().Msgf("[%s] Failed to parse JSON with Path on Consumer: %s", id, err)
		}

		var dataToSend ProviderRequest
		err = json.Unmarshal(msg.Body, &dataToSend)
		if err != nil {
			log.Error().Msgf("[%s] Failed to parse JSON to Send on Consumer: %s", id, err)
		}

		jsonData, err := json.Marshal(dataToSend)
		if err != nil {
			log.Error().Msgf("[%s] Failed to marshal data to JSON: %s", id, err)
		}

		req, err := http.NewRequest("POST", providerUrl+string(data.Path), bytes.NewReader(jsonData))
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
