package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	HeaderCorrelationID = "X-Correlation-Id"
)

type LoggerContextKey struct{}

type ILoggerMiddleware interface {
	LogResponse(next http.Handler) http.Handler
}

type LoggerMiddleware struct {
	logger *zerolog.Logger
}

func (l *LoggerMiddleware) LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		traceID := r.Header.Get(HeaderCorrelationID)
		method := r.Method
		path := r.URL.Path

		if traceID == "" {
			traceID = uuid.NewString()
			r.Header.Add(HeaderCorrelationID, traceID)
		}

		// Create a custom response writer to capture status code
		crw := &customResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		ctx := context.WithValue(r.Context(), LoggerContextKey{}, l.logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(crw, r)

		e := l.logger.Info()
		responseStatus := crw.statusCode
		if responseStatus >= http.StatusBadRequest || responseStatus < http.StatusOK {
			e = l.logger.Error()
		}

		// TODO Improve log message
		e.Str("traceID", traceID)
		e.Str("duration", time.Since(startTime).String())
		e.Int("statusCode", responseStatus)
		e.Str("method", method)
		e.Str("path", path)
		e.Msg("log response")
	})
}

// customResponseWriter captures the status code for logging
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {
	crw.statusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

func NewLoggerMiddleware(logger *zerolog.Logger) ILoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

func GetLoggerFromContext(ctx context.Context) *zerolog.Logger {
	logger, ok := ctx.Value(LoggerContextKey{}).(*zerolog.Logger)
	if !ok {
		logger = &log.Logger
		logger.Warn().Msg("not found logger in context")
	}

	return logger
}
