package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"eventspal/pkg/config"
	"eventspal/pkg/eventspal"
	"eventspal/pkg/logs"
)

func main() {
	config.Init()
	cfg := config.GetConfig()
	logs.Init(cfg)

	service := eventspal.NewService(cfg)
	data, err := service.GetWeatherAndEventsForLocation(context.Background(), 37, -122)
	if err != nil {
		log.Errorf("Failed to get events data: %s", err)
		return
	}

	log.WithFields(log.Fields{"data": data}).Infof("Successfully events %d", len(data))
}
