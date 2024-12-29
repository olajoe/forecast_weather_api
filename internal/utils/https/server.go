package https

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type RegisterRouteFunc func(route *mux.Router)

type Server struct {
	route  *mux.Router
	logger *zerolog.Logger
}

func (s *Server) ListenAndServe(port int) {
	logger := s.logger
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.route,
	}

	serverCtx := gracefulShutdown(logger, func(ctx context.Context) {
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Fatal().Msgf("HTTP server shutdown error: %s", err)
		}
	})

	logger.Info().Msgf("server listen and serve at %d", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Msgf("HTTP server listen and serve error: %s", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	logger.Info().Msg("graceful shutdown completed")
}

func gracefulShutdown(logger *zerolog.Logger, shutdown func(ctx context.Context)) context.Context {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Shutdown signal with grace period of 30 seconds
	shutdownCtx, cancelFunc := context.WithTimeout(serverCtx, 5*time.Second)
	defer cancelFunc()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		logger.Info().Msg("Received shutdown signal. Initiating graceful shutdown...")
		shutdown(shutdownCtx)
		serverStopCtx()
	}()

	return serverCtx
}

func NewServer(logger *zerolog.Logger, registerRoute RegisterRouteFunc) *Server {
	router := mux.NewRouter()

	registerRoute(router)

	return &Server{
		route:  router,
		logger: logger,
	}
}
