package weather

import (
	"context"
	"eventspal/pkg/common"
	"eventspal/pkg/config"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Client interface {
	GetWeatherByLatLong(ctx context.Context, lat float64, lon float64) ([]Weather, error)
}

type client7Timer struct {
	common.BaseClient
}

func NewClient(cfg config.Config) Client {
	return &client7Timer{
		BaseClient: common.BaseClient{
			BaseURL: cfg.Weather.BaseURL,
			Client:  &http.Client{},
		},
	}
}

func (c client7Timer) GetWeatherByLatLong(ctx context.Context, lat float64, lon float64) ([]Weather, error) {
	url := c.BaseURL

	var response = struct {
		Init       string `json:"init"`
		Dataseries []struct {
			Timepoint   int    `json:"timepoint"`
			Cloudcover  int    `json:"cloudcover"`
			LiftedIndex int    `json:"lifted_index"`
			PrecType    string `json:"prec_type"`
			PrecAmount  int    `json:"prec_amount"`
			Temp2m      int    `json:"temp2m"`
			Rh2m        string `json:"rh2m"`
			Weather     string `json:"weather"`
		} `json:"dataseries"`
	}{}
	body := map[string]string{}

	query := common.Query{
		"product": "civil",
		"output":  "json",
		"lat":     fmt.Sprintf("%.2f", lat),
		"lon":     fmt.Sprintf("%.2f", lon),
	}
	err := c.ExecuteRequestAndGetResponse(ctx, http.MethodGet, url, query, nil, &body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to query weather service: %w", err)
	}

	startYear, err := strconv.Atoi(response.Init[:4])
	if err != nil {
		return nil, fmt.Errorf("failed to get start year: %w", err)
	}
	startMonth, err := strconv.Atoi(response.Init[4:6])
	if err != nil {
		return nil, fmt.Errorf("failed to get start month: %w", err)
	}
	startDay, err := strconv.Atoi(response.Init[6:8])
	if err != nil {
		return nil, fmt.Errorf("failed to get start day: %w", err)
	}
	startHour, err := strconv.Atoi(response.Init[8:10])
	if err != nil {
		return nil, fmt.Errorf("failed to get start hour: %w", err)
	}

	startTime := time.Date(startYear, time.Month(startMonth), startDay, startHour, 0, 0, 0, time.UTC)
	var result []Weather
	for _, data := range response.Dataseries {
		result = append(result, Weather{
			Time:        startTime.Add(time.Duration(data.Timepoint) * time.Hour),
			Cloudcover:  data.Cloudcover,
			LiftedIndex: data.LiftedIndex,
			PrecType:    data.PrecType,
			PrecAmount:  data.PrecAmount,
			Temp2m:      data.Temp2m,
			Rh2m:        data.Rh2m,
			Weather:     data.Weather,
		})
	}

	return result, err
}
