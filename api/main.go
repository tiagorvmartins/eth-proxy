package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/tiagorvmartins/eth-proxy/api/app"
)

func init() {
	// Set gin mode
	mode := os.Getenv("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	// Setup the app
	app := app.SetupApp()

	// Read ADDR and port
	addr := os.Getenv("GIN_ADDR")
	port := os.Getenv("GIN_PORT")
	log.Info().Msgf("Starting service on http//:%s:%s", addr, port)
	if err := app.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		log.Fatal().Err(err).Msg("Error occurred while setting up the server")
	}
}
