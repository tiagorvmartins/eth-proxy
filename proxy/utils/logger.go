package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	host, err := os.Hostname()
	if err != nil {
		log.Logger = log.With().Str("host", "unknown").Logger()
	} else {
		log.Logger = log.With().Str("host", host).Logger()
	}

	log.Logger = log.With().Str("service", "rabbitmq-consumer").Logger()

	log.Logger = log.With().Caller().Logger()
}
