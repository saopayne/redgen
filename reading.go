package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Reading struct {
	Time    time.Time `json:"time"`
	State   float64   `json:"state"`
	Unit    string    `json:"unit"`
	MeterId string    `json:"meter_id,omitempty"`
	Sender  string    `json:"sender,omitempty"`
	Suit    string    `json:"suit,omitempty"`
}

func NewReading(date time.Time, unit string, interval, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state float64) Reading {
	baseDailyConsumptionDiv := baseDailyConsumption / 24 // 24 hours in a day
	var currentHour float64

	if variability > 0 {
		hourLowerBound := baseDailyConsumptionDiv - (variability / 10)
		hourUpperBound := baseDailyConsumptionDiv + (variability / 10)
		currentHour = RandomFloat64(hourLowerBound, hourUpperBound)
	} else {
		currentHour = baseDailyConsumptionDiv
	}
	hourBasedInterval := 60 / interval
	rawReading := float64(currentHour*hourBase*weekBase*monthBase) / hourBasedInterval
	if rawReading < 0 {
		rawReading = 0
	}
	return Reading{
		Time:  date,
		State: state + rawReading,
		Unit:  unit,
	}
}

func PrintJSONReading(reading Reading) {
	jsonBytes, _ := json.MarshalIndent(reading, "", "  ")
	fmt.Println(string(jsonBytes))
}

func RandomFloat64(lo float64, hi float64) float64 {
	rand.Seed(time.Now().UnixNano())
	lowerBound := int(lo * 10000000)
	upperBound := int(hi * 10000000)
	boundDifference := upperBound - lowerBound
	if boundDifference < 0 {
		boundDifference = 0
	}
	randomNumber := rand.Intn(boundDifference) + lowerBound
	return float64(randomNumber) / 10000000
}
