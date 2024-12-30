package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/olajoe/forecast_weather_api/internal/weather"
)

func RegisterRoutes(r *mux.Router, weatherHandler *weather.WeatherHandler) {
	corporateApi := r.PathPrefix("/weathers").Subrouter()
	corporateApi.HandleFunc("/daily/coordinates", weatherHandler.GetWeatherForecastDailyByCoordinates).Methods(http.MethodGet)
	corporateApi.HandleFunc("/daily/place", weatherHandler.GetWeatherForecastDailyByPlace).Methods(http.MethodGet)
}
