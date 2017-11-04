package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Reading represents a meter reading
type Reading struct {
	Time    time.Time `json:"time"`
	State   float64   `json:"state"`
	Unit    string    `json:"unit"`
	MeterId string    `json:"meter_id,omitempty"`
	Sender  string    `json:"sender,omitempty"`
	Suit    string    `json:"suit,omitempty"`
}

// NewReading creates a new Reading object
func NewReading(date time.Time, unit string, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state float64) Reading {
	baseDailyConsumption = baseDailyConsumption / 24 // 24 hours in a day

	hourLowerBound := baseDailyConsumption - variability
	hourUpperBound := baseDailyConsumption + variability
	currentHour := RandomHourValue(hourLowerBound, hourUpperBound)
	//the current hour will be multiplied by all the profiles
	//e.g 1.5 for Jan, 1 for Sat and 1 for 12:00 hours with currentHour 8 should produce 12
	rawReading := currentHour * hourBase * weekBase * monthBase

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

// randomHourValue returns a random value for an hour given a variability
func RandomHourValue(lo float64, hi float64) float64 {
	// to allow the the values to be represented
	// multiply the numbers by 100, get a random value and divide the value by 100 to get desired value
	rand.Seed(time.Now().UnixNano())
	lowerBound := int(lo * 100)
	upperBound := int(hi * 100)
	randomNumber := rand.Intn(upperBound-lowerBound) + lowerBound
	return float64(randomNumber / 100)
}
