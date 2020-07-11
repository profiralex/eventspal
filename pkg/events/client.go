package events

import (
	"context"
	"fmt"
	"net/http"
	"small_container/pkg/common"
	"small_container/pkg/config"
)

type Event struct {
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Distance float64 `json:"distance"`
	Units    string  `json:"units"`
}

type Client interface {
	GetEventsByLatLong(ctx context.Context, lat float64, lon float64, maxDistance float64) ([]Event, error)
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

func (c ticketmasterClient) GetEventsByLatLong(ctx context.Context, lat float64, lon float64, maxDistance float64) ([]Event, error) {
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
		"latlong": fmt.Sprintf("%.2f,%.2f", lat, lon),
	}
	err := c.ExecuteRequestAndGetResponse(ctx, http.MethodGet, url, query, nil, &body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query events service: %w", err)
	}

	var result []Event
	for _, data := range response.Embeded.Events {
		if data.Distance > maxDistance {
			continue
		}

		result = append(result, Event{
			Name:     data.Name,
			URL:      data.URL,
			Distance: data.Distance,
			Units:    data.Units,
		})
	}

	return result, err
}
