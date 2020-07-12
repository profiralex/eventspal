package events

import (
	"context"
	"fmt"
	"net/http"
	"eventspal/pkg/common"
	"eventspal/pkg/config"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	GetEventsByLatLong(ctx context.Context, lat float64, lon float64, radius int) ([]Event, error)
}

type ticketmasterClient struct {
	common.BaseClient
}

func NewClient(cfg config.Config) Client {
	return &ticketmasterClient{
		BaseClient: common.BaseClient{
			BaseURL: cfg.Events.BaseURL,
			Client:  &http.Client{},
			QueryParams: map[string]string{
				"apikey": cfg.Events.APIKey,
			},
		},
	}
}

func (c ticketmasterClient) GetEventsByLatLong(ctx context.Context, lat float64, lon float64, radius int) ([]Event, error) {
	url := c.BuildURL("/events.json")

	var response = struct {
		Embeded struct {
			Events []struct {
				Name     string  `json:"name"`
				URL      string  `json:"url"`
				Distance float64 `json:"distance"`
				Units    string  `json:"units"`

				Dates struct {
					Timezone         string `json:"timezone"`
					SpanMultipleDays bool   `json:"spanMultipleDays"`

					Distance float64 `json:"distance"`
					Units    float64 `json:"units"`

					Start struct {
						LocalDate      string `json:"localDate"`
						LocalTime      string `json:"localTime"`
						DateTBD        bool   `json:"dateTBD"`
						DateTBA        bool   `json:"dateTBA"`
						TimeTBA        bool   `json:"timeTBA"`
						NoSpecificTime bool   `json:"noSpecificTime"`
					} `json:"start"`

					End struct {
						LocalDate      string `json:"localDate"`
						Approximate    bool   `json:"approximate"`
						NoSpecificTime bool   `json:"noSpecificTime"`
					} `json:"end"`

					Status struct {
						Code string `json:"code"`
					} `json:"status"`
				} `json:"dates"`

				Embedded struct {
					Venues []struct {
						Location struct {
							Latitude  string `json:"latitude"`
							Longitude string `json:"longitude"`
						} `json:"location"`
					} `json:"venues"`
				} `json:"embedded"`
			} `json:"events"`
		} `json:"_embedded"`
	}{}
	body := map[string]string{}

	query := common.Query{
		"sort":    "distance,asc",
		"size":    "50",
		"radius":  fmt.Sprintf("%d", radius),
		"unit":    "km",
		"latlong": fmt.Sprintf("%.2f,%.2f", lat, lon),
	}
	err := c.ExecuteRequestAndGetResponse(ctx, http.MethodGet, url, query, nil, &body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query events service: %w", err)
	}

	var result []Event
	for _, data := range response.Embeded.Events {
		if data.Dates.Status.Code == "cancelled" {
			continue
		}

		if data.Dates.Start.DateTBA || data.Dates.Start.DateTBD || data.Dates.Start.NoSpecificTime {
			continue
		}

		dateParts := strings.Split(data.Dates.Start.LocalDate, "-")
		startYear, err := strconv.Atoi(dateParts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get start year: %w", err)
		}
		startMonth, err := strconv.Atoi(dateParts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get start month: %w", err)
		}
		startDay, err := strconv.Atoi(dateParts[2])
		if err != nil {
			return nil, fmt.Errorf("failed to get start day: %w", err)
		}

		timeParts := strings.Split(data.Dates.Start.LocalTime, ":")
		startHour, err := strconv.Atoi(timeParts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get start hour: %w", err)
		}
		startMinute, err := strconv.Atoi(timeParts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to get start hour: %w", err)
		}

		date := time.Date(startYear, time.Month(startMonth), startDay, startHour, startMinute, 0, 0, time.UTC)
		result = append(result, Event{
			Name:     data.Name,
			URL:      data.URL,
			Distance: data.Distance,
			Units:    data.Units,
			Date:     date,
		})
	}

	return result, err
}
