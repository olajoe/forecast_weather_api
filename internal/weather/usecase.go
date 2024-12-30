package weather

import (
	"fmt"
	"time"
)

type WeatherUsecase interface {
	GetWeatherDailyByCoordinates(queries GetWeatherDailyQuery) (map[string]interface{}, error)
	GetWeatherDailyByPlace(queries GetWeatherDailyQuery) (map[string]interface{}, error)
}

type weatherUsecase struct {
	weatherRepository WeatherRepository
}

func NewWeatherUsecase(weatherRepository WeatherRepository) WeatherUsecase {
	return &weatherUsecase{
		weatherRepository: weatherRepository,
	}
}

func (u *weatherUsecase) GetWeatherDailyByCoordinates(queries GetWeatherDailyQuery) (map[string]interface{}, error) {
	queryParams := buildGetWeatherDailyByCoordinatesQueryParams(queries)

	forecastResponse, err := u.weatherRepository.GetWeatherDailyByCoordinates(queryParams)
	if err != nil {
		return nil, err
	}

	result, err := mapWeatherForecastDailyResponseToResult(forecastResponse)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *weatherUsecase) GetWeatherDailyByPlace(queries GetWeatherDailyQuery) (map[string]interface{}, error) {
	queryParams := buildGetWeatherDailyByPlaceQueryParams(queries)

	forecastResponse, err := u.weatherRepository.GetWeatherDailyByPlace(queryParams)
	if err != nil {
		return nil, err
	}

	result, err := mapWeatherForecastDailyResponseToResult(forecastResponse)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func mapWeatherForecastDailyResponseToResult(response *WeatherForecastDailyResponse) (map[string]interface{}, error) {
	var result map[string]interface{}

	for _, forecast := range response.WeatherForecasts {
		result = map[string]interface{}{
			"location":  fulfillLocationValue(forecast.Location),
			"forecasts": []interface{}{},
		}

		for _, forecastItem := range forecast.Forecasts {
			tData, err := time.Parse(time.RFC3339, forecastItem.Time)
			if err != nil {
				return nil, err
			}

			forecastData := fulfillForecastDataValue(forecastItem.Data)

			result["forecasts"] = append(result["forecasts"].([]interface{}), map[string]interface{}{
				"time": tData.Format(time.DateOnly),
				"data": forecastData,
			})

		}

	}

	return result, nil
}

func fulfillLocationValue(location Location) map[string]interface{} {
	result := map[string]interface{}{
		"lat": location.Lat,
		"lon": location.Lon,
	}

	if location.Province != nil {
		result["province"] = location.Lat
	}

	if location.Amphoe != nil {
		result["amphoe"] = location.Amphoe
	}

	if location.Tambon != nil {
		result["tambon"] = location.Tambon
	}

	if location.Region != nil {
		result["region"] = location.Region
	}

	if location.Geocode != nil {
		result["geocode"] = location.Geocode
	}

	if location.AreaType != nil {
		result["areatype"] = location.AreaType
	}

	return result
}

func fulfillForecastDataValue(forecastData ForecastData) map[string]interface{} {
	result := map[string]interface{}{}

	if forecastData.TcMin != nil {
		result["tcMin"] = fmt.Sprintf("%v °C", *forecastData.TcMin)
	}

	if forecastData.TcMax != nil {
		result["tcMax"] = fmt.Sprintf("%v °C", *forecastData.TcMax)
	}

	if forecastData.Rh != nil {
		result["rh"] = fmt.Sprintf("%v %%", *forecastData.Rh)
	}

	if forecastData.Slp != nil {
		result["slp"] = fmt.Sprintf("%v hPa", *forecastData.Slp)
	}

	if forecastData.Psfc != nil {
		result["psfc"] = fmt.Sprintf("%v Pa", *forecastData.Psfc)
	}

	if forecastData.Rain != nil {
		result["rain"] = fmt.Sprintf("%v mm", *forecastData.Rain)
	}

	if forecastData.Ws10m != nil {
		result["ws10m"] = fmt.Sprintf("%v m/s", *forecastData.Ws10m)
	}

	if forecastData.Wd10m != nil {
		result["wd10m"] = fmt.Sprintf("%v °", *forecastData.Wd10m)
	}

	if forecastData.CloudLow != nil {
		result["cloudLow"] = fmt.Sprintf("%v %%", *forecastData.CloudLow)
	}

	if forecastData.CloudMed != nil {
		result["cloudMed"] = fmt.Sprintf("%v %%", *forecastData.CloudMed)
	}

	if forecastData.CloudHigh != nil {
		result["cloudHigh"] = fmt.Sprintf("%v %%", *forecastData.CloudHigh)
	}

	if forecastData.Swdown != nil {
		result["swDown"] = fmt.Sprintf("%v W/m^2", *forecastData.Swdown)
	}

	if forecastData.Cond != nil {
		result["cond"] = mapConditionToValue(*forecastData.Cond)
	}

	// if forecastData.Ws != nil {
	// 	result["ws"] = forecastData.Ws
	// }

	// if forecastData.Wd != nil {
	// 	result["wd"] = forecastData.Wd
	// }

	return result
}

func mapConditionToValue(condition float64) string {
	switch condition {
	case 1:
		return "ท้องฟ้าแจ่มใส (Clear)"
	case 2:
		return "มีเมฆบางส่วน (Partly cloudy)"
	case 3:
		return "เมฆเป็นส่วนมาก (Cloudy)"
	case 4:
		return "มีเมฆมาก (Overcast)"
	case 5:
		return "ฝนตกเล็กน้อย (Light rain)"
	case 6:
		return "ฝนปานกลาง (Moderate rain)"
	case 7:
		return "ฝนตกหนัก (Heavy rain)"
	case 8:
		return "ฝนฟ้าคะนอง (Thunderstorm)"
	case 9:
		return "อากาศหนาวจัด (Very cold)"
	case 10:
		return "อากาศหนาว (Cold)"
	case 11:
		return "อากาศเย็น (Cool)"
	case 12:
		return "อากาศร้อนจัด (Very hot)"
	default:
		return ""
	}
}
