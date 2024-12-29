package weather

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/olajoe/forecast_weather_api/internal/config"
	"github.com/olajoe/forecast_weather_api/internal/utils/https"
)

type WeatherRepository interface {
	GetWeatherDailyByCoordinates(queryParams map[string]string) (*WeatherForecastDailyCoordinatesResponse, error)
}

type weatherRepository struct {
	client      *req.Client
	baseUrl     string
	accessToken string
}

func NewWeatherRepository(
	client *req.Client,
	cfg *config.Configuration,
) WeatherRepository {
	return &weatherRepository{
		client:      client,
		baseUrl:     cfg.Tmd.Url,
		accessToken: cfg.Tmd.AccessToken,
	}
}

func (r *weatherRepository) GetWeatherDailyByCoordinates(queryParams map[string]string) (*WeatherForecastDailyCoordinatesResponse, error) {
	var resultBody WeatherForecastDailyCoordinatesResponse
	var errResp https.ErrorResponse

	resp, err := r.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", r.accessToken)).
		SetQueryParams(queryParams).
		SetSuccessResult(&resultBody).
		SetErrorResult(&errResp).
		Get(fmt.Sprintf("%s/forecast/location/daily/at", r.baseUrl))

	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, errResp
	}

	return &resultBody, nil
}
