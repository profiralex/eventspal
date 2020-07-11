package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"small_container/pkg/config"
	"small_container/pkg/events"
	"small_container/pkg/logs"
)

func main() {
	config.Init()
	cfg := config.GetConfig()
	logs.Init(cfg)

	client := events.NewClient(cfg)
	data, err := client.GetEventsByLatLong(context.Background(), 37, -122, 10)
	if err != nil {
		log.Errorf("Failed to get events data: %s", err)
		return
	}

	log.WithFields(log.Fields{"data": data}).Infof("Successfully events weather data %d", len(data))
}
