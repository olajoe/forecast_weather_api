package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/imroc/req/v3"
	"github.com/olajoe/forecast_weather/internal/config"
	v1 "github.com/olajoe/forecast_weather/internal/routes/v1"
	"github.com/olajoe/forecast_weather/internal/validator"
	"github.com/olajoe/forecast_weather/internal/weather"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/olajoe/forecast_weather/internal/utils/https"
	"github.com/olajoe/forecast_weather/pkg/logging"
)

func main() {
	cfg := config.New()
	logger := logging.New(cfg.LogLevel)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	r := mux.NewRouter()
	corsHandler := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins(strings.Split(cfg.Cors.Origins, ",")),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "x-api-key", "x-secret-key"}),
	)

	r.HandleFunc("/healthz", https.HealthCheckHandler).Methods(http.MethodGet)

	v1Router := r.PathPrefix("/v1").Subrouter()
	v1Router.Use()

	// dependency
	client := req.C().SetTimeout(5 * time.Second)
	_validator := validator.NewValidator()
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.IgnoreUnknownKeys(true)

	// repository
	weatherRepo := weather.NewWeatherRepository(client, cfg)

	// usecase
	weatherUsecase := weather.NewWeatherUsecase(weatherRepo)

	// handler
	weatherHandler := weather.NewWeatherHandler(_validator, schemaDecoder, weatherUsecase)

	v1.RegisterRoutes(v1Router, weatherHandler)

	// Create signal channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	serve := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: corsHandler(r)}
	// Start server in goroutine
	go func() {
		logger.Info().Msgf("Server starting at port %d", cfg.Port)
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Msgf("Server failed: %s", err.Error())
		}
	}()
	// Wait for interrupt signal
	<-stop
	// Cancel root context on shutdown
	rootCancel()
	logger.Info().Msg("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer shutdownCancel()

	// Then shutdown HTTP server
	if err := serve.Shutdown(shutdownCtx); err != nil {
		// Change to WARN since timeout during shutdown is expected
		logger.Warn().Msgf("Server shutdown exceeded timeout: %s", err.Error())
	} else {
		logger.Info().Msg("Server shutdown completed successfully")
	}

	logger.Info().Msg("Server exiting")

}
