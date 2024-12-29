package logging

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (r *loggingResponseWriter) WriteHeader(
	statusCode int,
) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.StatusCode = statusCode
}

type HandlerFunc = func(http.Handler) http.Handler

func NewMiddleware(
	logger *zerolog.Logger,
) HandlerFunc {
	return func(
		next http.Handler,
	) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			ctx := r.Context()
			ctx = logger.WithContext(ctx)
			lw := loggingResponseWriter{ResponseWriter: w}

			method := r.Method
			endpoint := r.URL.Path
			startTime := time.Now()
			traceID := r.Header.Get("x-correlation-id")

			defer func() {
				if lw.StatusCode == 0 {
					lw.WriteHeader(http.StatusInternalServerError)
				}

				log := logger.Info()
				if lw.StatusCode >= http.StatusInternalServerError {
					log = logger.Error()
				}

				log.
					Str("method", method).
					Str("path", endpoint).
					Str("traceID", traceID).
					Int("statusCode", lw.StatusCode).
					Str("duration", time.Since(startTime).String()).
					Send()
			}()

			next.ServeHTTP(&lw, r.WithContext(ctx))
		})
	}
}
