package weather

import (
	"fmt"
	"strings"
)

type Location struct {
	Province string `json:"province,omitempty"`
	AreaType string `json:"areatype,omitempty"`
	Tambon   string `json:"tambon,omitempty"`
	Amphoe   string `json:"amphoe,omitempty"`
	Region   string `json:"region,omitempty"`
	Geocode  string `json:"geocode,omitempty"`

	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type ForecastData struct {
	TcMin     *float64 `json:"tc_min"`
	TcMax     *float64 `json:"tc_max"`
	Rh        *float64 `json:"rh"`
	Slp       *float64 `json:"slp"`
	Psfc      *float64 `json:"psfc"`
	Rain      *float64 `json:"rain"`
	Ws10m     *float64 `json:"ws10m"`
	Wd10m     *float64 `json:"wd10m"`
	Ws        *float64 `json:"ws"`
	Wd        *float64 `json:"wd"`
	CloudLow  *float64 `json:"cloudlow"`
	CloudMed  *float64 `json:"cloudmed"`
	CloudHigh *float64 `json:"cloudhigh"`
	Swdown    *float64 `json:"swdown"`
	Cond      *float64 `json:"cond"`
}

type Forecast struct {
	Time string       `json:"time"`
	Data ForecastData `json:"data"`
}

type WeatherForecastDaily struct {
	Location  Location   `json:"location,omitempty"`
	Forecasts []Forecast `json:"forecasts"`
}

type WeatherForecastDailyResponse struct {
	WeatherForecasts []WeatherForecastDaily `json:"WeatherForecasts"`
}

type GetWeatherDailyQuery struct {
	//at
	Lat float32 `schema:"lat,omitempty"`
	Lon float32 `schema:"lon,omitempty"`

	// place
	Province string `schema:"province,omitempty"`
	Amphoe   string `schema:"amphoe,omitempty"`
	Tambon   string `schema:"tambon,omitempty"`
	Subarea  bool   `schema:"subarea,omitempty"` // 0 or 1 default 0

	// region
	Region string `schema:"region,omitempty"`

	Date     string `schema:"date"`
	Duration int    `schema:"duration"`
	Fields   string `schema:"fields"`
}

func buildGetWeatherDailyCordinatesQuery(
	lat float32,
	lon float32,
	date string,
	duration int,
	fields []string,
) GetWeatherDailyQuery {
	return GetWeatherDailyQuery{
		Lat:      lat,
		Lon:      lon,
		Date:     date,
		Duration: duration,
		Fields:   strings.Join(fields, ","),
	}
}

func buildGetWeatherDailyPlaceQuery(
	province string,
	amphoe string,
	tambon string,
	subarea bool,
	date string,
	duration int,
	fields []string,
) GetWeatherDailyQuery {
	return GetWeatherDailyQuery{
		Province: province,
		Amphoe:   amphoe,
		Tambon:   tambon,
		Subarea:  subarea,
		Date:     date,
		Duration: duration,
		Fields:   strings.Join(fields, ","),
	}
}

func buildGetWeatherDailyByCoordinatesQueryParams(queries GetWeatherDailyQuery) map[string]string {
	return map[string]string{
		"lat":      fmt.Sprintf("%f", queries.Lat),
		"lon":      fmt.Sprintf("%f", queries.Lon),
		"date":     queries.Date,
		"duration": fmt.Sprintf("%d", queries.Duration),
		"fields":   queries.Fields,
	}
}

func buildGetWeatherDailyByPlaceQueryParams(queries GetWeatherDailyQuery) map[string]string {
	result := map[string]string{}
	if queries.Province != "" {
		result["province"] = queries.Province
	}
	if queries.Amphoe != "" {
		result["amphoe"] = queries.Amphoe
	}
	if queries.Tambon != "" {
		result["tambon"] = queries.Tambon
	}
	if queries.Subarea {
		result["subarea"] = fmt.Sprintf("%t", queries.Subarea)
	}
	if queries.Date != "" {
		result["date"] = queries.Date
	}
	if queries.Duration != 0 {
		result["duration"] = fmt.Sprintf("%d", queries.Duration)
	}
	if queries.Fields != "" {
		result["fields"] = queries.Fields
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
