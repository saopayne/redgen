package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gizak/termui"
	"github.com/olekukonko/tablewriter"
	"sort"
)

const DefaultProfile = "DefaultProfile"

var (
	defaultProfilePath  = filepath.Join(".", "profiles")
	defaultReadingsPath = filepath.Join(".", "readings")
	parseFolderPath     = filepath.Join(".", "parsers")
	parseableFileName   = "parseable.txt"
	unparseableFileName = "unparseable.txt"
	defaultProfileName  = "default_config.json"
)

type Profile struct {
	Name                 string             `json:"name"`
	BaseDailyConsumption float64            `json:"baseDailyConsumption"`
	HourlyProfiles       map[string]float64 `json:"hourlyProfiles"`
	WeeklyProfiles       map[string]float64 `json:"weeklyProfiles"`
	MonthlyProfiles      map[string]float64 `json:"monthlyProfiles"`
	Variability          float64            `json:"variability"`
	Unit                 string             `json:"unit"`
	Interval             float64            `json:"interval"`
	Start                time.Time          `json:"startAt"`
	Readings             []Reading          `json:"readings"`
}

type Start struct {
	Year  int    `json:"year,omitempty"`
	Day   int    `json:"day,omitempty"`
	Month string `json:"month"`
	Hour  int    `json:"hour"`
}

func Save(p Profile) {
	err := WriteProfileToFile(p, defaultProfilePath, SanitizeName(p.Name)+".json")
	if err != nil {
		log.Fatal(err)
	}
}

func SaveReadings(p Profile, destinationFolderPath string) {
	err := WriteReadingsToFile(p, filepath.Join(destinationFolderPath, SanitizeName(p.Name)+".json"))
	if err != nil {
		log.Fatal(err)
	}
}

func SanitizeName(name string) string {
	trimmedName := strings.TrimSpace(name)
	lowerCaseName := strings.ToLower(trimmedName)
	sanitizedName := strings.Replace(lowerCaseName, " ", "_", -1)
	re := regexp.MustCompile("[[:^ascii:]]")
	sanitizedName = re.ReplaceAllLiteralString(sanitizedName, "")

	return sanitizedName
}

func (p Profile) SetName(newName string) {
	p.Name = newName
	Save(p)
}

func (p Profile) SetVariability(newVariability float64) {
	p.Variability = newVariability
	Save(p)
}

func (p Profile) SetUnit(newUnit string) {
	p.Unit = newUnit
	Save(p)
}

func (p Profile) SetBaseDailyConsumption(baseConsumption float64) {
	p.BaseDailyConsumption = baseConsumption
	Save(p)
}

// StartAt mocks a clock based on the configuration file (Year,Month, Day and Hour are configurable)
// The clock will always start in the current year(if the year is not set), day 1, in 0 minutes, seconds and nseconds.
func (p Profile) StartAt() (time.Time, float64, error) {
	var (
		state float64
		date  time.Time
		err   error
		count = len(p.Readings)
	)

	if count != 0 {
		state = p.Readings[count-1].State
		lastWriteTime := p.Readings[count-1].Time
		if err != nil {
			log.Fatal(err.Error())
		}
		lastWriteTime = lastWriteTime.Add(time.Minute * time.Duration(p.Interval))
		return lastWriteTime, state, nil
	} else {
		return p.Start, state, nil
	}

	return date, state, nil
}

func WriteProfileToFile(profile Profile, path string, profileFile string) error {

	jsonBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}
	os.MkdirAll(path, os.ModePerm)
	sanitizedProfileName := SanitizeName(profileFile)
	err = ioutil.WriteFile(filepath.Join(path, sanitizedProfileName), jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteReadingsToFile(profile Profile, profileFile string) error {
	jsonBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stat(defaultReadingsPath); os.IsNotExist(err) {
		os.MkdirAll(defaultReadingsPath, os.ModePerm)
	}
	err = ioutil.WriteFile(profileFile, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewProfileFromJson(profileBytes []byte) (Profile, error) {
	profile := Profile{}
	err := json.Unmarshal(profileBytes, &profile)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func CreateDefaultProfile(name string) Profile {
	selectedName := ""
	if name == "" {
		selectedName = DefaultProfile
	} else {
		selectedName = name
	}
	var hourKeys []string
	for k := range defaultHourlyProfile {
		hourKeys = append(hourKeys, k)
	}
	sort.Strings(hourKeys)
	modifiedHourlyProfiles := map[string]float64{}
	for _, k := range hourKeys {
		modifiedHourlyProfiles[k] = defaultHourlyProfile[k]
	}

	return Profile{
		Name:                 selectedName,
		BaseDailyConsumption: 18,
		HourlyProfiles:       modifiedHourlyProfiles,
		WeeklyProfiles:       defaultWeeklyProfile,
		MonthlyProfiles:      defaultMonthlyProfile,
		Variability:          5,
		Interval:             15,
		Unit:                 "kW",
		Start:                time.Date(2017, 01, 01, 00, 00, 00, 00, time.UTC),
		Readings:             make([]Reading, 0),
	}
}

func GenerateReadings(profile Profile, path string) {
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
			hourBase  = profile.HourlyProfiles[date.Format("12")]
			weekBase  = profile.WeeklyProfiles[date.Format("Mon")]
			monthBase = profile.MonthlyProfiles[date.Format("Jan")]
		)
		reading := NewReading(date, profile.Unit, profile.Interval, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state)
		profile.Readings = append(profile.Readings, reading)

		PrintJSONReading(reading)

		state = profile.Readings[len(profile.Readings)-1].State
		SaveReadings(profile, path)

		time.Sleep(5 * time.Millisecond)
		date = date.Add(time.Duration(profile.Interval) * time.Minute)
	}
}

func GenerateSingleReading(profile Profile) Profile {
	var (
		baseDailyConsumption = profile.BaseDailyConsumption
		variability          = profile.Variability
	)

	date, state, err := profile.StartAt()

	if err != nil {
		log.Fatal(err.Error())
	}

	var (
		hourBase  = profile.HourlyProfiles[date.Format("12")]
		weekBase  = profile.WeeklyProfiles[date.Format("Mon")]
		monthBase = profile.MonthlyProfiles[date.Format("Jan")]
	)
	if len(profile.Readings) > 0 {
		state = profile.Readings[len(profile.Readings)-1].State
	}
	reading := NewReading(date, profile.Unit, profile.Interval, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state)
	profile.Readings = append(profile.Readings, reading)
	return profile
}

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
			timeLabels[i] = profile.Readings[i].Time.String()
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

func PlotMonthlyProfilesChart(profile Profile) {
	data := [][]string{}
	var monthKeys []string
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

func PlotWeeklyProfilesChart(profile Profile) {
	data := [][]string{}
	var weekKeys []string
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

func PlotHourlyProfilesChart(profile Profile) {
	data := [][]string{}
	var hourKeys []string
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
