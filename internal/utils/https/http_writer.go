package https

import (
	"encoding/json"
	"net/http"

	"github.com/olajoe/forecast_weather/internal/middlewares"
	"github.com/rs/zerolog"
)

func WriteError(w http.ResponseWriter, r *http.Request, res ErrorResponse) {
	logger := middlewares.GetLoggerFromContext(r.Context())
	logger.Error().Msg(res.Error())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.Error().Msg(err.Error())
	}
}

func WriteResponse(
	w http.ResponseWriter,
	logger *zerolog.Logger,
	statusCode int,
	payload any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger.Error().Msgf("json encode error: %s", err)
	}
}
