package eventspal

import (
	"context"
	"eventspal/pkg/config"
	"eventspal/pkg/events"
	"eventspal/pkg/weather"
	"fmt"
)

type Service interface {
	GetWeatherAndEventsForLocation(ctx context.Context, lat float64, long float64) ([]Event, error)
}

type eventsPalService struct {
	weatherClient weather.Client
	eventsClient  events.Client
}

func NewService(cfg config.Config) Service {
	return &eventsPalService{
		weatherClient: weather.NewClient(cfg),
		eventsClient:  events.NewClient(cfg),
	}
}

func (e eventsPalService) GetWeatherAndEventsForLocation(ctx context.Context, lat float64, long float64) ([]Event, error) {
	weatherData, err := e.weatherClient.GetWeatherByLatLong(ctx, lat, long)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}
	minWeatherDate := weatherData[0].Time
	maxWeatherDate := weatherData[len(weatherData)-1].Time

	eventsData, err := e.eventsClient.GetEventsByLatLong(ctx, lat, long, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var result []Event
	for _, eventData := range eventsData {
		if eventData.Date.Before(minWeatherDate) || eventData.Date.After(maxWeatherDate) {
			continue
		}

		eventWeather := weatherData[0]
		for _, potentialEventWeather := range weatherData {
			if potentialEventWeather.Time.After(eventData.Date) {
				eventWeather = potentialEventWeather
				break
			}
		}

		result = append(result, Event{
			Name: eventData.Name,
			URL:  eventData.Name,
			Date: eventData.Date,
			Weather: Weather{
				Cloudcover:  eventWeather.Cloudcover,
				LiftedIndex: eventWeather.LiftedIndex,
				PrecType:    eventWeather.PrecType,
				PrecAmount:  eventWeather.PrecAmount,
				Temp2m:      eventWeather.Temp2m,
				Rh2m:        eventWeather.Rh2m,
				Weather:     eventWeather.Weather,
			},
		})
	}

	return result, nil
}
