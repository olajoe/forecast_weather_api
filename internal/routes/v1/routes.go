package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olajoe/forecast_weather/internal/weather"
)

func RegisterRoutes(r *mux.Router, weatherHandler *weather.WeatherHandler) {
	corporateApi := r.PathPrefix("/weathers").Subrouter()
	corporateApi.HandleFunc("/daily/coordinates", weatherHandler.GetWeatherDailyByCoordinates).Methods(http.MethodGet)
}
