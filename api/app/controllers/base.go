package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tiagorvmartins/eth-proxy/api/app/models"
	"github.com/tiagorvmartins/eth-proxy/api/utils"
)

func BasePath(c *gin.Context) {
	var msg models.Request
	request_id := c.GetString("x-identifier-id")

	if binderr := c.ShouldBindJSON(&msg); binderr != nil {
		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	connectionString := os.Getenv("RMQ_URL")

	rmqProducer := utils.RMQProducer{
		Queue:            os.Getenv("QUEUE_NAME"),
		ConnectionString: connectionString,
		QueueCallback:    c.GetString("x-token"),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Str("request_id", request_id).
			Msg("Error occurred while marshaling message")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}
	reply := rmqProducer.PublishMessage("application/json", []byte(msgBytes), request_id)
	c.JSON(http.StatusOK, reply)
}
