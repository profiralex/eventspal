package eventspal

import "time"

type Event struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Date    time.Time `json:"start"`
	Weather Weather   `json:"weather"`
}

type Weather struct {
	Cloudcover  int    `json:"cloudcover"`
	LiftedIndex int    `json:"lifted_index"`
	PrecType    string `json:"prec_type"`
	PrecAmount  int    `json:"prec_amount"`
	Temp2m      int    `json:"temp2m"`
	Rh2m        string `json:"rh2m"`
	Weather     string `json:"weather"`
}
