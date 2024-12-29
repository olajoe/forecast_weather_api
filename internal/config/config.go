package config

import (
	"log"
	"strings"

	"github.com/olajoe/forecast_weather_api/pkg/logging"
	"github.com/spf13/viper"
)

type Configuration struct {
	Port     int
	Cors     CorsConfig
	LogLevel logging.Level `mapstructure:"log_level"`

	Tmd TmdConfig
}

type CorsConfig struct {
	Origins string
}

type TmdConfig struct {
	Url         string
	AccessToken string
}

func New() *Configuration {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Cannot read .env file, %s \n", err)
	}

	viper.AutomaticEnv()

	cfg := Configuration{
		Port: viper.GetInt("port"),
		Cors: CorsConfig{
			// Origins: viper.GetString("cors_origins"),
			Origins: "*",
		},
		LogLevel: logging.Level(viper.GetInt("log_level")),

		Tmd: TmdConfig{
			Url:         viper.GetString("tmd_url"),
			AccessToken: viper.GetString("tmd_access_token"),
		},
	}

	return &cfg
}
