package main

import (
	"github.com/go-ozzo/ozzo-validation"

	"errors"
	"fmt"
	"strings"
)

// constants that define the units of energy
// useful when validating some parts of the configuration
const (
	w  = 1
	W  = w // Watts
	kW = w * 1000
	KW = kW // Kilowatts
	mW = kW * 1000
	MW = mW // Megawatts
	gW = mW * 100
	GW = gW // Gigawatts
	end
)

// getValueForUnit gives the expanded float value
// for a given unit of energy
// > 1w = 1
// > 1kW = 1000
func GetValueforUnit(unit string) float64 {
	unit = strings.ToUpper(unit)
	switch unit {
	case "W":
		return 1
	case "KW":
		return 1000
	case "MW":
		return 1000000
	case "GW":
		return 1000000000
	default:
		return 0
	}
}

// IsUnitValid checks if the given unit is one of [W, KW, MW or GW]
func IsUnitValid(value string) bool {
	return value < string(end)
}

// IsValueInList checks if a given string is present in a list of strings
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		// compare with case insensitivity such that
		// IsValueInList("Jan", "jan") will return true
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}

// IsIntValueInList checks if an int exists in a list of integers
func IsIntValueInList(value int, list []int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}

// ValidateHourlyProfiles confirms the value is
// greater than zero and that the hour key is a valid value
// It also checks that for higher units, an abnormal large number isn't set
func ValidateHourlyProfiles(p Profile) error {
	hoursOfDay := []string{"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11",
		"12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"}
	var err error
	if validation.IsEmpty(p.HourlyProfiles) {
		err = errors.New("hourly profiles must be set")
	}
	for hour, value := range p.HourlyProfiles {
		if !IsValueInList(hour, hoursOfDay) {
			// The hour set isn't a valid houe
			errorString := fmt.Sprintf("the hour %+v is not a valid hour, should be one of: %+v", hour, hoursOfDay)
			err = errors.New(errorString)
		}
		if value <= 0 {
			// The value set for an hour must be a minimum of 1
			errorString := fmt.Sprintf("the minimum for any hour should be 1 , hour %+v has the value: %d", hour, value)
			err = errors.New(errorString)
		}
		if p.Unit != "w" && p.Unit != "kW" && value > 100 {
			// The value set for an hour is too large
			errorString := fmt.Sprintf("the hour %+v is too large with value: %d", hour, value)
			err = errors.New(errorString)
		}
	}

	return err
}

// validateWeeklyProfiles checks that for a given profile,
// the day of the week set exists in a list of week days
// It also ensures that no zero value is set and that an abnormal large value isn't given for a unit
func ValidateWeeklyProfiles(p Profile) error {
	var err error
	if validation.IsEmpty(p.WeeklyProfiles) {
		err = errors.New("weekly profiles must be set")
	}
	for weekDay, value := range p.WeeklyProfiles {
		if _, ok := weekDays[weekDay]; !ok {
			// The week entered is not valid
			errorString := fmt.Sprintf("the value set for week %+v is not valid, must be one of: %+v", weekDay, weekDays)
			err = errors.New(errorString)
		}
		if value <= 0 {
			// The value set for a week must be a minimum of 1
			errorString := fmt.Sprintf("the minimum for any week should be 1 , week %+v has the value: %d", weekDay, value)
			err = errors.New(errorString)
		}
		if p.Unit != "w" && p.Unit != "kW" && value > 100 {
			// The value set for a week is too large
			errorString := fmt.Sprintf("the week %+v is too large with value: %d", weekDay, value)
			err = errors.New(errorString)
		}
	}

	return err
}

// validateMonthlyProfiles checks that for a given profile,
// the month set exists in a list of of possible months
// It also ensures that no zero value is set and that an abnormal large value isn't given for a unit
func ValidateMonthlyProfiles(p Profile) error {
	var err error
	if validation.IsEmpty(p.HourlyProfiles) {
		err = errors.New("monthly profiles must be set")
	}
	for aMonth, value := range p.MonthlyProfiles {
		if _, ok := months[aMonth]; !ok {
			// The week entered is not valid
			errorString := fmt.Sprintf("the value set for month %+v is not valid, must be one of: %+v", aMonth, months)
			err = errors.New(errorString)
		}
		if value <= 0 {
			// The value set for a month must be a minimum of 1
			errorString := fmt.Sprintf("the minimum for any month should be 1 , month %+v has the value: %d", aMonth, value)
			err = errors.New(errorString)
		}
		if p.Unit != "w" && p.Unit != "kW" && value > 100 {
			// The value set for a week is too large
			errorString := fmt.Sprintf("the month %+v is too large with value: %d", aMonth, value)
			err = errors.New(errorString)
		}
	}

	return err
}

// validateVariability checks that the variability value is non-negative
// and that for a variability value, it doesn't render the consumption to be a negative number if too large
func ValidateVariability(p Profile) error {
	var err error
	variability := p.Variability
	if validation.IsEmpty(variability) {
		err = errors.New("variability must be set")
	}
	if variability < 0 || variability >= 10 {
		errorString := "variability cannot be lower than 0 and greater than 10"
		err = errors.New(errorString)
	}
	unit := p.Unit
	baseDailyConsumption := GetValueforUnit(unit) * p.BaseDailyConsumption
	consumptionLimit := baseDailyConsumption / 24
	if consumptionLimit-p.Variability <= 0 {
		err = errors.New("either the variability is too high or the base consumption is too low")
	}

	return err
}

// validateBaseDailyConsumption ensures no missing base consumption value
func ValidateBaseDailyConsumption(p Profile) error {
	var err error
	unit := p.Unit
	if validation.IsEmpty(unit) {
		err = errors.New("base daily consumption must be set for a profile")
	}
	baseDailyConsumption := p.BaseDailyConsumption
	baseDailyConsumption = GetValueforUnit(unit) * baseDailyConsumption
	reasonableLimit := baseDailyConsumption / 24
	if reasonableLimit < 0 {
		err = fmt.Errorf("kindly set the base daily consumption to be above 24 for with this unit: %s ", unit)
	}

	return err
}

// validateInterval checks that the interval is set to either >= 1
func ValidateInterval(p Profile) error {
	var err error
	if p.Interval < 1 {
		errorString := "interval must be either be 1 or greater than 1 minute"
		err = errors.New(errorString)
	}

	return err
}

// validateStart confirms that the values set for the start are valid for the hour, month and year set
func ValidateStart(p Profile) error {
	months := []string{"Jan", "Feb", "Mar", "Apr", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	hoursOfDay := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
		12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}

	var err error
	start := p.Start
	if validation.IsEmpty(start) {
		err = errors.New("the start date of the profile must be set")
	}

	if !IsValueInList(start.Month, months) {
		// The week entered is not valid
		errorString := fmt.Sprintf("the value set for month %+v is not valid, must be one of: %+v", start.Month, months)
		err = errors.New(errorString)
	}

	if !IsIntValueInList(start.Hour, hoursOfDay) {
		// The hour set isn't a valid hour
		errorString := fmt.Sprintf("the hour %+v is not a valid hour, should be one of: %+v", start.Hour, hoursOfDay)
		err = errors.New(errorString)
	}

	if start.Year <= 1990 || start.Year > 2030 {
		err = errors.New("the year set must be within 1990 and 2030")
	}

	return err
}

// validateReadings only ensures that the readings have the correct unit and no empty time
func ValidateReadings(p Profile) error {
	var err error
	// no readings, might be the default configuration file
	if len(p.Readings) == 0 {
		err = nil
	}
	// has some readings
	totalReadings := len(p.Readings)
	if totalReadings > 0 {
		for i := 0; i < totalReadings; i++ {
			if GetValueforUnit(p.Readings[i].Unit) == 0 || len(p.Readings[i].Unit) != 2 {
				// the unit is not valid
				err = fmt.Errorf("the reading %d has an invalid unit %s :)", i+1, p.Readings[i].Unit)
			}
			if p.Readings[i].Time == "" {
				// no time set for the reading
				err = fmt.Errorf("the reading %d has no time:)", i+1)
			}
		}

	}
	return err
}

// validateName ensures names have a minimum length of 5 and max of 50
func ValidateName(p Profile) error {
	return validation.ValidateStruct(&p,
		// Name cannot be empty, and the length must between 5 and 50
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
	)
}

// validateUnit checks that the provided unit is a valid one
func ValidateUnit(p Profile) error {
	err := validation.ValidateStruct(&p,
		// Base daily consumption cannot be empty, and the length must be 2
		validation.Field(&p.Unit, validation.Required, validation.Length(2, 2)),
	)
	if err != nil {
		return err
	}
	isValid := IsUnitValid(p.Unit)
	if !isValid {
		err = errors.New("THE UNIT SHOULD BE ONE OF: [ w, kW, mW, gW ]")
	}
	return err
}

// validateProfile validates all the properties of the Profile struct
// calls the appropriate validation method for each field
func (p *Profile) ValidateProfile() error {
	err := ValidateName(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Profile name is valid")
	err = ValidateUnit(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Profile unit is valid")
	err = ValidateHourlyProfiles(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Hourly Profiles are valid")
	err = ValidateWeeklyProfiles(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Weekly Profiles are valid")
	err = ValidateMonthlyProfiles(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("Monthly Profiles are valid")
	err = ValidateStart(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("The Start time of the profile is valid")
	err = ValidateReadings(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("The readings are valid")
	err = ValidateVariability(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("The profile variability specified is valid")
	err = ValidateInterval(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("The profile interval specified is valid")
	err = ValidateBaseDailyConsumption(*p)
	if err != nil {
		panic(err)
	}
	fmt.Println("The base daily consumption is valid")
	return nil
}