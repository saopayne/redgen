package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gizak/termui"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/olekukonko/tablewriter"
	"sort"
)

const DefaultProfile = "DefaultProfile"

var (
	defaultProfilePath = filepath.Join(".", "profiles")
	defaultProfileName = "default_config.json"
)

// Profile represents a demonstration account profile
type Profile struct {
	Name                 string             `json:"name"`
	BaseDailyConsumption float64            `json:"baseDailyConsumption"`
	HourlyProfiles       map[string]float64 `json:"hourlyProfiles"`
	WeeklyProfiles       map[string]float64 `json:"weeklyProfiles"`
	MonthlyProfiles      map[string]float64 `json:"monthlyProfiles"`
	Variability          float64            `json:"variability"`
	Unit                 string             `json:"unit"`
	Interval             time.Duration      `json:"interval"`
	Start                Start              `json:"startAt"`
	Readings             []Reading          `json:"readings"`
}

// Start represents a partially mocked clock
type Start struct {
	Year  int    `json:"year,omitempty"`
	Day   int 	 `json:"day,omitempty"`
	Month string `json:"month"`
	Hour  int    `json:"hour"`
}

// Save writes the profile JSON into a file, so it can be recovered later
func (p Profile) Save() {
	err := WriteProfileToFile(p, strings.ToLower(p.Name)+"_readings.json")
	if err != nil {
		log.Fatal(err)
	}
}

func (p Profile) SetName(newName string) {
	p.Name = newName
	p.Save()
}

func (p Profile) SetVariability(newVariability float64) {
	p.Variability = newVariability
	p.Save()
}

func (p Profile) SetUnit(newUnit string) {
	p.Unit = newUnit
	p.Save()
}

func (p Profile) SetBaseDailyConsumption(baseConsumption float64) {
	p.BaseDailyConsumption = baseConsumption
	p.Save()
}

// StartAt mocks a clock based on the configuration file (Year,Month, Day and Hour are configurable)
// The clock will always start in the current year(if the year is not set), day 1, in 0 minutes, seconds and nseconds.
func (p Profile) StartAt() (time.Time, float64, error) {
	var (
		state    float64
		date     time.Time
		err      error
		readings = len(p.Readings)
	)

	if readings != 0 {
		state = p.Readings[readings-1].State
		date, err = time.Parse("2010-01-02 15:04:05.999999999 MST", p.Readings[readings-1].Time)
		if err != nil {
			log.Fatal(err.Error())
		}
		date = date.Add(p.Interval * time.Minute)
	} else {
		if p.Start.Hour > 23 {
			return time.Time{}, state, fmt.Errorf("invalid starting hour in configuration file: %d", p.Start.Hour)
		}
		month, ok := months[p.Start.Month]
		if ok {
			startYear := p.Start.Year
			startDay := p.Start.Day
			if validation.IsEmpty(p.Start.Year) {
				startYear = time.Now().Year()
			}
			if validation.IsEmpty(p.Start.Day) {
				startDay = 1
			}
			return time.Date(startYear, month, startDay, p.Start.Hour, 0, 0, 0, time.UTC), state, nil
		}
		return time.Time{}, state, fmt.Errorf("invalid starting month in configuration file: %s", p.Start.Month)
	}

	return date, state, nil
}

// createProfile marshalls a Profile object into a JSON file
func WriteProfileToFile(profile Profile, profileFile string) error {
	jsonBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}

	os.MkdirAll(defaultProfilePath, os.ModePerm)

	err = ioutil.WriteFile(filepath.Join(defaultProfilePath, profileFile), jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// newProfileFromJson unmarshalls the bytes of configuration file into a Profile object
func NewProfileFromJson(profileBytes []byte) (Profile, error) {
	profile := Profile{}
	err := json.Unmarshal(profileBytes, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

// createDefaultProfile returns a Profile object with the default configuration
func CreateDefaultProfile() Profile {
	return Profile{
		Name:                 DefaultProfile,
		BaseDailyConsumption: 18,
		HourlyProfiles:       defaultHourlyProfile,
		WeeklyProfiles:       defaultWeeklyProfile,
		MonthlyProfiles:      defaultMonthlyProfile,
		Variability:          5,
		Interval:             15,
		Unit:                 "kW",
		Start: Start{
			Year:  2017,
			Month: "Jan",
			Day: 1,
			Hour:  06,
		},
		Readings: make([]Reading, 0),
	}
}

// startDemonstration generate a reading based on the configured interval
func GenerateReadings(profile Profile) {
	var (
		baseDailyConsumption = profile.BaseDailyConsumption
		variability          = profile.Variability
	)

	date, state, err := profile.StartAt()
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		var (
			hourBase  = profile.HourlyProfiles[date.Format("15")]
			weekBase  = profile.WeeklyProfiles[date.Format("Mon")]
			monthBase = profile.MonthlyProfiles[date.Format("Jan")]
		)

		reading := NewReading(date.String(), profile.Unit, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state)
		profile.Readings = append(profile.Readings, reading)

		PrintJSONReading(reading)

		state = profile.Readings[len(profile.Readings)-1].State
		profile.Save()

		time.Sleep(5 * time.Second)
		date = date.Add(profile.Interval * time.Minute)
	}
}

// plotReadingsChart displays the readings generated for a configuration with ASCII art
func PlotReadingsChart(profile Profile) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	stateReadings := (func() []float64 {
		n := len(profile.Readings)
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = profile.Readings[i].State
		}
		return ps
	})()

	dataLabels := (func() []string {
		readingsLength := len(profile.Readings)
		timeLabels := make([]string, readingsLength)
		for i := range timeLabels {
			timeLabels[i] = profile.Readings[i].Time
		}
		return timeLabels
	})()

	chartLabel := fmt.Sprintf("Andy Readings Chart for Profile  -%s-", profile.Name)
	lc := termui.NewLineChart()
	lc.BorderLabel = chartLabel
	lc.Mode = "dot"
	lc.Data = stateReadings
	lc.DataLabels = dataLabels[:]
	lc.Width = 140
	lc.Height = 18
	lc.X = 0
	lc.Y = 0
	lc.AxesColor = termui.ColorWhite
	lc.LineColor = termui.ColorGreen | termui.AttrBold

	p := termui.NewPar(":PRESS q TO QUIT READINGS CHART")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Text Box"
	p.BorderFg = termui.ColorCyan

	termui.Render(lc, p)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()
}

// plotMonthlyProfilesChart displays the monthly values
// set for this configuration in a table format with ASCII art
func PlotMonthlyProfilesChart(profile Profile) {
	data := [][]string{}
	var monthKeys []string
	// to preserve ordering of the values
	for k := range profile.MonthlyProfiles {
		monthKeys = append(monthKeys, k)
	}
	sort.Strings(monthKeys)

	for _, month := range monthKeys {
		value := profile.MonthlyProfiles[month]
		monthValue := fmt.Sprintf("%.6f", value)
		data = append(data, []string{month, monthValue})
	}

	table := tablewriter.NewWriter(os.Stdout)
	profileMonthHeader := fmt.Sprintf("Month for config --%s--", profile.Name)
	table.SetHeader([]string{profileMonthHeader, "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

// plotWeeklyProfilesChart displays the weekly values
// set for this configuration in a table format with ASCII art
func PlotWeeklyProfilesChart(profile Profile) {
	data := [][]string{}
	var weekKeys []string
	// to preserve ordering of the values
	for k := range profile.WeeklyProfiles {
		weekKeys = append(weekKeys, k)
	}
	sort.Strings(weekKeys)

	for _, weekDay := range weekKeys {
		value := profile.WeeklyProfiles[weekDay]
		weekDayValue := fmt.Sprintf("%.6f", value)
		data = append(data, []string{weekDay, weekDayValue})
	}

	table := tablewriter.NewWriter(os.Stdout)
	profileWeekHeader := fmt.Sprintf("Week day for config --%s--", profile.Name)
	table.SetHeader([]string{profileWeekHeader, "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}

// plotHourlyProfilesChart displays the hourly values
// set for this configuration in a table format with ASCII art
func PlotHourlyProfilesChart(profile Profile) {
	data := [][]string{}
	var hourKeys []string
	// to preserve iteration ordering
	for k := range profile.HourlyProfiles {
		hourKeys = append(hourKeys, k)
	}
	sort.Strings(hourKeys)
	for _, hour := range hourKeys {
		value := profile.HourlyProfiles[hour]
		hourValue := fmt.Sprintf("%.6f", value)
		data = append(data, []string{hour, hourValue})
	}

	table := tablewriter.NewWriter(os.Stdout)
	profileHourHeader := fmt.Sprintf("Hour for config --%s--", profile.Name)
	table.SetHeader([]string{profileHourHeader, "Value"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCenterSeparator("|")
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
