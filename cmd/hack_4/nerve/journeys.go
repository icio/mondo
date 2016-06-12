package main

import (
	"github.com/icio/mondo/cmd/hack_4/tflol"
)

type Journey struct {
	Origin      string `json:"origin"`
	Destination string `json:"dest"`
	Cost        int    `json:"cost"`
}

type Journeys struct {
}

func (j *Journeys) RegisterJourneys(journeys []Journey) {

}
