package config

import (
	"fmt"
	"github.com/profiralex/goconfig"
)

type Config struct {
	App struct {
		DebugLevel string `cfg:"DEBUG_LEVEL" cfg-default:"info"`
		Port       string `cfg:"PORT" cfg-default:"8080"`
		ApiKey     string `cfg:"API_KEY" cfg-default:"Please set api key"`
		Version    string `cfg:"APP_VER" cfg-default:"v1"`
	}

	Weather struct {
		BaseURL string `cfg:"WEATHER_BASE_URL" cfg-default:"http://www.7timer.info/bin/api.pl"`
	}

	Events struct {
		BaseURL string `cfg:"EVENTS_BASE_URL" cfg-default:"https://app.ticketmaster.com/discovery/v2"`
		APIKey  string `cfg:"EVENTS_API_KEY"`
	}
}

var config Config

func Init() {
	config = Config{}
	err := goconfig.Load(&config, &goconfig.EnvProvider{}, false)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}
}

func GetConfig() Config {
	return config
}
