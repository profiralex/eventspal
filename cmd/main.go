package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"small_container/pkg/config"
	"small_container/pkg/logs"
	"small_container/pkg/weather"
)

func main() {
	config.Init()
	cfg := config.GetConfig()
	logs.Init(cfg)

	client := weather.NewClient(cfg)
	data, err := client.GetWeatherByLatLong(context.Background(), 23.6, 46.5)
	if err != nil {
		log.Errorf("Failed to get weather data: %s", err)
	}

	log.WithFields(log.Fields{"data": data}).Info("Successfully queried weather data")
}
