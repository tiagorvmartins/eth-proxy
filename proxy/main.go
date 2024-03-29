package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tiagorvmartins/eth-proxy/proxy/handlers"
	"github.com/tiagorvmartins/eth-proxy/proxy/utils"
)

func main() {
	connectionString := os.Getenv("RMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")
	providerNamesStr := os.Getenv("PROVIDER_NAMES")
	if providerNamesStr == "" {
		panic("No providers configured!")
	}

	availableProviders := strings.Split(providerNamesStr, ",")

	threadsPerConsumerStr := os.Getenv("NUM_THREADS_PER_CONSUMER")
	threadsPerConsumer, err := strconv.Atoi(threadsPerConsumerStr)
	if err != nil {
		panic(err)
	}

	go func() {
		// From: https://github.com/kubernetes/kubernetes/blob/master/test/images/agnhost/liveness/server.go
		startedOrUpdated := time.Now()
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			duration := time.Since(startedOrUpdated)
			if duration.Seconds() > 60 {
				w.WriteHeader(500)
				w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
			} else {
				w.WriteHeader(200)
				w.Write([]byte("ok"))
				startedOrUpdated = time.Now()
			}
		})
		if err := http.ListenAndServe(":8090", nil); err != nil {
			panic(err)
		}
	}()

	forever := make(chan bool)
	for _, provider := range availableProviders {
		for i := 0; i < threadsPerConsumer; i++ {
			consumer := utils.RMQConsumer{
				Queue:            queueName,
				ConnectionString: connectionString,
				MsgHandler:       handlers.Handler(provider),
				Id:               fmt.Sprintf("%s_%d", provider, i),
			}
			go consumer.Consume()
		}
	}
	<-forever
}
