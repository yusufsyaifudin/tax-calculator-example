package restapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// This package must be imported to make swaggerFiles working.
	_ "github.com/yusufsyaifudin/tax-calculator-example/docs"
)

// Config is an configuration needed by this restapi Server to run.
type Config struct {
	Address string
	Test    bool
}

var conf *Config
var logger = log.With().Str("pkg", "restapi").Caller().Logger()
var stopped = false

// Router is a gin.Engine type. This exported because we need to run a httptest server in integration test (main_test.go).
var Router *gin.Engine

// init will initiating some logic when this file is called.
func init() {
	gin.SetMode(gin.ReleaseMode)
	Router = gin.New()
	Router.Use(middleware())
}

// Configure will start a listener to HTTP.
// @title Tax Calculator Example
// @version 1.0
// @description This is a sample API to generate tax foreach user.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func Configure(config *Config) {
	conf = config

	registerRoute()

	if !config.Test {
		logger.Level(zerolog.Disabled)
		Router.Use(Logger())
	}

	for _, route := range Router.Routes() {
		if config.Test {
			continue
		}

		log.Info().
			Str("method", route.Method).
			Str("path", route.Path).
			Str("handler", route.Handler).
			Msg("")
	}
}

// Run will start the server
func Run() error {
	return Router.Run(conf.Address)
}

// registerRoute will register all route for this application.
func registerRoute() {
	parent := context.Background()

	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := Router.Group("/api/v1")

	protectedEndpointMiddleware := ChainMiddleware(middlewareAuthTokenCheck)

	v1.POST("/register", WrapGin(parent, register))
	v1.POST("/login", WrapGin(parent, login))

	v1.POST("/tax", WrapGin(parent, protectedEndpointMiddleware(createNewTax)))
	v1.GET("/tax", WrapGin(parent, protectedEndpointMiddleware(getTaxes)))
}

// Shutdown gracefully when some signal from OS tell that system should be down.
func Shutdown() {
	log.Info().Msg("not receiving requests anymore")
	stopped = true
}

func middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// handle panic
		defer func() {
			if err := recover(); err != nil {
				errStr := fmt.Sprint(err)

				log.Error().
					Str("err", errStr).
					Msgf("panic while executing %s", ctx.Request.URL.String())

				ctx.Status(http.StatusInternalServerError)
			}
		}()

		// check if flash is shutting down
		// if it's the case then don't receive anymore requests
		if stopped {
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		log.Info().
			Str("endpoint", ctx.Request.RequestURI).
			Dur("duration", duration).
			Int("status", ctx.Writer.Status()).
			Msg("finished handling api request")
	}
}
