package events

import "time"

type Event struct {
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Distance float64   `json:"distance"`
	Units    string    `json:"units"`
	Date     time.Time `json:"start"`
}
