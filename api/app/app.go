package app

import (
	"fmt"
	"os"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"github.com/tiagorvmartins/eth-proxy/api/app/middlewares"
	routers "github.com/tiagorvmartins/eth-proxy/api/app/routes"
)

func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	app := gin.New()

	// Gin Prometheus exporter
	app.Use(ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexEndpoint: "^/metrics|^/healthz"}))
	app.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	app.Use(gin.Recovery())

	proxyPingCheck := checks.NewPingCheck("http://proxy:8090/healthz", "GET", 1000, nil, nil)
	rabbitmqCheckQueueEndpoint := fmt.Sprintf("http://guest:guest@rabbitmq:15672/api/queues/%%2F/%s", os.Getenv("QUEUE_NAME"))
	rabbitPingCheck := checks.NewPingCheck(rabbitmqCheckQueueEndpoint, "GET", 1000, nil, nil)
	healthcheck.New(app, config.DefaultConfig(), []checks.Check{proxyPingCheck, rabbitPingCheck})

	// We disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	// Add Token, Identifier, IdentifierLogger and CORS middleware
	log.Info().Msg("Adding cors, request id and request logging middleware")
	app.Use(
		middlewares.Identifier(),
		middlewares.IdentifierLogger(),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	// Setup routers
	log.Info().Msg("Setting up routers")
	routers.SetupRouters(app)

	return app
}
