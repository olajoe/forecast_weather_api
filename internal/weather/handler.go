package weather

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/olajoe/forecast_weather_api/internal/utils/https"
	"github.com/olajoe/forecast_weather_api/pkg/logging"
)

type WeatherHandler struct {
	validate       *validator.Validate
	schemaDecoder  *schema.Decoder
	weatherUsecase WeatherUsecase
}

func NewWeatherHandler(
	validate *validator.Validate,
	schemaDecoder *schema.Decoder,
	weatherUsecase WeatherUsecase,
) *WeatherHandler {
	return &WeatherHandler{
		validate:       validate,
		schemaDecoder:  schemaDecoder,
		weatherUsecase: weatherUsecase,
	}
}

func (h *WeatherHandler) GetWeatherDailyByCoordinates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.Ctx(ctx)

	var queries struct {
		Lat      float32  `schema:"lat,required"`
		Lon      float32  `schema:"lon,required"`
		Date     string   `schema:"date"`     // YYYY-MM-DD
		Duration int      `schema:"duration"` // default 1 days, max 126 days
		Fields   []string `schema:"fields"`   // fields=tc_max,tc_min,rh,slp,psfc,cloudlow,cloudmed,cloudhigh,cond
		/* fields // https://data.tmd.go.th/nwpapi/doc/apidoc/location/forecast_daily.html
		tc_max
		tc_min
		rh
		slp
		psfc
		cloudlow
		cloudmed
		cloudhigh
		cond
		*/
	}

	if err := h.schemaDecoder.Decode(&queries, r.URL.Query()); err != nil {
		https.WriteError(w, r, https.NewErrorResponseBadRequest(err))
		return
	}

	if err := h.validate.Struct(queries); err != nil {
		https.WriteError(w, r, https.NewErrorResponseBadRequest(err))
		return
	}

	queriesData := buildGetWeatherDailyCordinatesQuery(
		queries.Lat,
		queries.Lon,
		queries.Date,
		queries.Duration,
		queries.Fields,
	)

	result, err := h.weatherUsecase.GetWeatherDailyByCoordinates(queriesData)
	if err != nil {
		https.WriteError(w, r, https.NewErrorResponseInternalServerError(err))
		return
	}

	https.WriteResponse(w, logger, http.StatusOK, map[string]interface{}{
		"data": result,
	})
}
