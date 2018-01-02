package main

import (
	binary "encoding/binary"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	parseableFileFullPath   = filepath.Join(parseFolderPath, parseableFileName)
	unparseableFileFullPath = filepath.Join(parseFolderPath, unparseableFileName)
)

const maxInt32 = 1<<(32-1) - 1

func RunAppCommands() (string, error) {
	enteredCommand := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch enteredCommand {
	case redgenVersion.FullCommand():
		return helpVersion, nil
	case redgenConfigStart.FullCommand():
		return InitGenerator()
	case redgenConfigGenerate.FullCommand():
		return CmdGenerate(*redgenConfigGenerateArg)
	case redgenConfigInit.FullCommand():
		return CmdInit()
	case redgenConfigPreview.FullCommand():
		if *redgenConfigPreviewArg != "" {
			CmdPreviewAction(*redgenConfigPreviewArg, *redgenConfigPreviewArgTime)
			return "", nil
		}
		CmdPreviewAction(defaultProfileName, *redgenConfigPreviewArgTime)
		return "", nil
	case redgenConfigSend.FullCommand():
		if *redgenConfigSendArg != "" {
			CmdSendReadingsToServer(*redgenConfigSendArg)
			return "", nil
		}
		return helpMsg, nil
	case redgenConfigValidate.FullCommand():
		if *redgenConfigValidateArg != "" {
			CmdValidateAction(*redgenConfigValidateArg, true)
		}
		CmdValidateAction(*redgenConfigValidateArg, false)
		return "", nil
	case redgenConfigProfile.FullCommand():
		if *redgenConfigProfileArg != "" {
			CmdProfileAction(*redgenConfigProfileArg)
			return "", nil
		}
		return helpMsg, nil
	case redgenConfigShow.FullCommand():
		if *redgenConfigShowFileName != "" {
			if *redgenConfigShowDateFlag != "" {
				ShowDateConsumption(*redgenConfigShowFileName, *redgenConfigShowDateFlag)
			}
			return "", nil
		}
		return helpMsg, nil
	case redgenConfigClear.FullCommand():
		return "", nil
	default:
		return helpMsg, fmt.Errorf("unknown command: %s", enteredCommand)
	}
	return "", nil
}

func CmdInit() (string, error) {
	err := WriteProfileToFile(CreateDefaultProfile(""), defaultProfilePath, defaultProfileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	return "Default profile file created into ./profiles", nil
}

func CmdGenerate(arg string) (string, error) {
	if arg != "" {
		namesArr := strings.Split(arg, ".")
		extractedName := SanitizeName(namesArr[0])
		err := WriteProfileToFile(CreateDefaultProfile(extractedName), defaultProfilePath, arg)
		if err != nil {
			log.Fatal(err.Error())
		}
		return fmt.Sprintf("%s file created into ./profiles", arg), nil
	} else {
		return helpMsg, nil
	}
}

func CmdProfileAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	GenerateReadings(profile, defaultReadingsPath)
}

func CmdPreviewAction(filename string, flag string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	profile.Validate()
	fmt.Println("Generating sample consumption data for the profile")
	profile = GeneratePreviewData(profile, flag)

	startTime := profile.Readings[0].Time.String()
	formattedDay := startTime[:10]
	formattedMonth := startTime[:7]
	formattedYear := startTime[:4]
	if flag == "day" {
		fmt.Println("Show for day")
		ShowDayConsumption(profile, formattedDay)
	} else if flag == "month" {
		fmt.Println("Show for month")

		ShowMonthConsumption(profile, formattedMonth)
	} else if flag == "year" {
		fmt.Println("Show for year")
		ShowYearConsumption(profile, formattedYear)
	} else {
		fmt.Println("back to default, showing for day")
		ShowDayConsumption(profile, formattedDay)
	}
}

func GeneratePreviewData(profile Profile, timeFmt string) Profile {
	var (
		baseDailyConsumption = profile.BaseDailyConsumption
		variability          = profile.Variability
	)

	date, state, err := profile.StartAt()

	if err != nil {
		log.Fatal(err.Error())
	}

	var valuesToGenerate int
	switch timeFmt {
	case "day":
		valuesToGenerate = 100
		break
	case "month":
		valuesToGenerate = 3000
		break
	case "year":
		valuesToGenerate = 36000
		break
	}

	var readings []Reading
	for i := 0; i < valuesToGenerate; i++ {
		hourValue := strconv.Itoa(date.Hour())
		weekValue := date.Weekday().String()
		monthValue := date.Month().String()
		var (
			hourBase  = profile.HourlyProfiles[hourValue]
			weekBase  = profile.WeeklyProfiles[weekValue[:3]]
			monthBase = profile.MonthlyProfiles[monthValue[:3]]
		)
		reading := NewReading(date, profile.Unit, profile.Interval, baseDailyConsumption, hourBase, weekBase, monthBase, variability, state)
		readings = append(readings, reading)

		time.Sleep(5 * time.Millisecond)
		date = date.Add(time.Duration(profile.Interval) * time.Minute)
	}
	profile.Readings = readings
	return profile
}

func CmdValidateSingleFile(filename string) {
	fmt.Printf("attempting to validate profile with name: %s", filename)
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = profile.Validate()
	if err != nil {
		fmt.Println("The profile is not a valid profile")
	}
	fmt.Println("------The profile configuration is valid------------")
}

func CmdValidateAction(filename string, singleFile bool) {
	if singleFile {
		CmdValidateSingleFile(filename)
	} else {
		fileList, err := ioutil.ReadDir(defaultProfilePath)
		if err != nil {
			panic("could not complete the listing of profile configurations")
		}
		for _, file := range fileList {
			CmdValidateSingleFile(file.Name())
		}
	}
}

func CmdSendReadingsToServer(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = profile.Validate()
	if err != nil {
		log.Println("The profile is not valid", profile)
	}
	librarianService := new(LibrarianService)
	// send the readings to the api at this point
	resp, err := librarianService.sendReadingsAction(profile)
	fmt.Println("SendReadingsAction() httpResp", resp)
}

// InitGenerator starts the application and calls the function which runs
// every 5 seconds for now (should be 15 mins ideally)
func InitGenerator() (string, error) {
	t := time.NewTicker(time.Second * 5)
	for {
		SendReadingsOnStart()
		<-t.C
	}
	return "", nil
}

func SendReadingsOnStart() {
	parseableFileNamesList := GetAppendedParsedFileNames()
	// check if a file with the last reading exists, if not create and leave empty
	// generate readings at the specified interval in the config and add to the /readings/filename_readings.json file
	// when you stop the app, and start again, it should get the last reading, compare the time interval of the last reading with current time
	// account for the time lost by updating the state with an estimated value for the time it was offline but it should continue appending
	// reading from now (just add the state from say State: 10 -> 12 -> [...offline for three missed readings] -> 20 [14->16->18 skipped] but state added
	for _, filename := range parseableFileNamesList {
		profile, _ := GetProfileFromFile(filename)
		profile = GenerateSingleReading(profile)
		librarianService := new(LibrarianService)
		resp, err := librarianService.sendReadingsAction(profile)
		if err == nil {
			if resp != nil {
				if resp.StatusCode == 200 || resp.StatusCode == 201 {
					SaveReadings(profile, defaultReadingsPath)
				} else {
					log.Println("Sent reading responded with a status other than 200 OK success", resp.StatusCode)
				}
			}
		} else {
			log.Println("Encountered an error while sending reading to the API", err)
		}
		time.Sleep(time.Second * 1)
	}
}

func GetProfileFromFile(filename string) (Profile, bool) {
	if _, err := os.Stat(defaultReadingsPath); os.IsNotExist(err) {
		os.Mkdir(defaultReadingsPath, os.ModePerm)
	}

	readingsFile := filepath.Join(defaultReadingsPath, filename)
	if _, err := os.Stat(readingsFile); os.IsNotExist(err) {
		_, err = os.Create(readingsFile)
		profile, _ := GetProfileFromJson(filepath.Join(defaultProfilePath, filename))
		_ = WriteProfileToFile(profile, defaultReadingsPath, filename)
		return profile, false
	} else {
		profile, _ := GetProfileFromJson(filepath.Join(defaultReadingsPath, filename))
		return profile, true
	}
}

func GetProfileFromJson(filepath string) (Profile, error) {
	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = profile.Validate()
	return profile, err
}

func GetAppendedParsedFileNames() (validFileNames []string) {
	parseableNamesList, _ := SplitParseableConfigFiles()
	os.MkdirAll(parseFolderPath, os.ModePerm)
	if _, err := os.Stat(parseableFileFullPath); os.IsNotExist(err) {
		_, err = os.Create(parseableFileFullPath)
	}
	parseableNamesBytes := EncodeFileNamesAsBytes(parseableNamesList)
	err := ioutil.WriteFile(parseableFileFullPath, parseableNamesBytes, 0644)
	checkWithPanic(err)

	parseableFileBytes, err := ioutil.ReadFile(parseableFileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	parseableFileNamesList := DecodeBytesAsFileNames(parseableFileBytes)

	return parseableFileNamesList
}

func SplitParseableConfigFiles() (parseableList []string, unparseableList []string) {
	files, err := ioutil.ReadDir(defaultProfilePath)
	if err != nil {
		log.Fatal(err)
	}
	var parseableNamesList []string
	var unparseableNamesList []string
	for _, f := range files {
		fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, f.Name()))
		if err != nil {
			log.Fatal(err.Error())
		}
		profile, err := NewProfileFromJson(fileBytes)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = profile.Validate()
		if err != nil {
			unparseableNamesList = append(unparseableNamesList, f.Name())
			panic(err)
		} else {
			parseableNamesList = append(parseableNamesList, f.Name())
		}
	}
	return parseableNamesList, unparseableNamesList
}

func checkWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func writeLen(b []byte, l int) []byte {
	if 0 > l || l > maxInt32 {
		panic("writeLen: invalid length")
	}
	var lb [4]byte
	binary.BigEndian.PutUint32(lb[:], uint32(l))
	return append(b, lb[:]...)
}

func readLen(b []byte) ([]byte, int) {
	if len(b) < 4 {
		panic("readLen: invalid length")
	}
	l := binary.BigEndian.Uint32(b)
	if l > maxInt32 {
		panic("readLen: invalid length")
	}
	return b[4:], int(l)
}
func DecodeBytesAsFileNames(b []byte) []string {
	b, ls := readLen(b)
	s := make([]string, ls)
	for i := range s {
		b, ls = readLen(b)
		s[i] = string(b[:ls])
		b = b[ls:]
	}
	return s
}

func EncodeFileNamesAsBytes(s []string) []byte {
	var b []byte
	b = writeLen(b, len(s))
	for _, ss := range s {
		b = writeLen(b, len(ss))
		b = append(b, ss...)
	}
	return b
}

func ShowDateConsumption(filename string, date string) {
	profile, err := GetProfileFromJson(filepath.Join(defaultReadingsPath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	splitDate := strings.Split(date, "-")
	if len(splitDate) == 1 {
		ShowYearConsumption(profile, date)
	} else if len(splitDate) == 2 {
		ShowMonthConsumption(profile, date)
	} else if len(splitDate) == 3 {
		ShowDayConsumption(profile, date)
	}
}

func GetValueForEnergyUnit(unit string) float64 {
	var unitValueMap = map[string]float64{
		"mW": 1 / 1000,
		"W":  1,
		"kW": 1000,
		"MW": 1000000,
		"GW": 100000000,
	}
	if unitValueMap[unit] != 0 {
		return unitValueMap[unit]
	}
	return 0
}

func ShowDayConsumption(profile Profile, day string) {
	mDay := day
	intervalMultiplier := 60 / profile.Interval
	var newReadingValues []Reading
	for _, reading := range profile.Readings {
		if strings.Contains(reading.Time.String(), mDay) {
			reading.State = reading.State * intervalMultiplier
			newReadingValues = append(newReadingValues, reading)
		}
	}
	// calculate the differences between two hours for each reading
	readingHourMap := make(map[int]float64)
	for _, reading := range newReadingValues {
		if readingHourMap[reading.Time.Hour()] == 0 {
			readingHourMap[reading.Time.Hour()] = reading.State
		}
	}

	var hourKeys []int
	var hourLabels []int
	for k, _ := range readingHourMap {
		hourKeys = append(hourKeys, k)
	}
	sort.Ints(hourKeys)
	unit := ""
	for _, k := range hourKeys {
		switch profile.Unit {
		case "kW":
			hourLabels = append(hourLabels, int(readingHourMap[k]*GetValueForEnergyUnit(profile.Unit)))
			unit = "W"
			break
		default:
			hourLabels = append(hourLabels, int(readingHourMap[k]))
			unit = profile.Unit
			break
		}

	}
	newHourKeys := GetStringSliceFromInt(hourKeys)
	header := fmt.Sprintf("Hourly Consumption for %s in (%s) ", mDay, unit)
	PlotBarChart(hourLabels, newHourKeys, header)
}

func ShowMonthConsumption(profile Profile, month string) {
	mMonth := month
	intervalMultiplier := 60 / profile.Interval
	var newReadingValues []Reading
	for _, reading := range profile.Readings {
		if strings.Contains(reading.Time.String(), mMonth) {
			reading.State = reading.State * intervalMultiplier
			newReadingValues = append(newReadingValues, reading)
		}
	}

	readingMonthMap := make(map[int]float64)
	for _, reading := range newReadingValues {
		readingMonthMap[reading.Time.Day()] = reading.State
	}

	var daysKeys []int
	var daysLabels []float64
	for k, _ := range readingMonthMap {
		daysKeys = append(daysKeys, k)
		daysLabels = append(daysLabels, readingMonthMap[k])
	}

	sort.Ints(daysKeys)

	var normDaysLabels []string
	for _, k := range daysKeys {
		normDaysLabels = append(normDaysLabels, strconv.Itoa(k))
	}

	normDaysKeys := GetIntSliceFromFloat(daysLabels)
	header := fmt.Sprintf("Monthly Consumption for %s in (%s) ", mMonth, profile.Unit)
	PlotBarChart(normDaysKeys, normDaysLabels, header)
}

func ShowYearConsumption(profile Profile, year string) {
	mYear := year
	intervalMultiplier := 60 / profile.Interval
	var newReadingValues []Reading
	for _, reading := range profile.Readings {
		if strings.Contains(reading.Time.String(), mYear) {
			reading.State = reading.State * intervalMultiplier
			newReadingValues = append(newReadingValues, reading)
		}
	}

	readingMonthMap := make(map[time.Month]float64)
	for _, reading := range newReadingValues {
		readingMonthMap[reading.Time.Month()] = reading.State
	}

	var monthKeys []time.Month
	var monthLabels []float64
	for k, _ := range readingMonthMap {
		monthKeys = append(monthKeys, k)
		monthLabels = append(monthLabels, readingMonthMap[k])
	}
	var normMonthLabels []string
	for _, k := range monthKeys {
		normMonthLabels = append(normMonthLabels, k.String())
	}
	normMonthKeys := GetIntSliceFromFloat(monthLabels)
	header := fmt.Sprintf("Monthly Consumption for %s in (%s) ", mYear, profile.Unit)
	PlotBarChart(normMonthKeys, normMonthLabels, header)
}
