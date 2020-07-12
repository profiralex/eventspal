package weather

import "time"

type Weather struct {
	Time        time.Time `json:"time"`
	Cloudcover  int       `json:"cloudcover"`
	LiftedIndex int       `json:"lifted_index"`
	PrecType    string    `json:"prec_type"`
	PrecAmount  int       `json:"prec_amount"`
	Temp2m      int       `json:"temp2m"`
	Rh2m        string    `json:"rh2m"`
	Weather     string    `json:"weather"`
}


