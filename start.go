package main

import (
	binary "encoding/binary"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	parseableFileFullPath   = filepath.Join(parseFolderPath, parseableFileName)
	unparseableFileFullPath = filepath.Join(parseFolderPath, unparseableFileName)
)

const maxInt32 = 1<<(32-1) - 1

func RunAppCommands() (string, error) {
	fmt.Println(os.Args[:])
	enteredCommand := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch enteredCommand {
	case andyVersion.FullCommand():
		return helpVersion, nil
	case andyConfigStart.FullCommand():
		return InitGenerator()
	case andyConfigGenerate.FullCommand():
		return CmdGenerate(*andyConfigGenerateArg)
	case andyConfigInit.FullCommand():
		return CmdInit()
	case andyConfigPreview.FullCommand():
		if *andyConfigPreviewArg != "" {
			CmdPreviewAction(*andyConfigPreviewArg)
		}
		CmdPreviewAction(defaultProfileName)
		return "", nil
	case andyConfigSend.FullCommand():
		if *andyConfigSendArg != "" {
			CmdSendReadingsToServer(*andyConfigSendArg)
			return "", nil
		}
		return helpMsg, nil
	case andyConfigValidate.FullCommand():
		if *andyConfigValidateArg != "" {
			CmdValidateAction(*andyConfigValidateArg, true)
		}
		CmdValidateAction(*andyConfigValidateArg, false)
		return "", nil
	case andyConfigProfile.FullCommand():
		if *andyConfigProfileArg != "" {
			CmdProfileAction(*andyConfigProfileArg)
			return "", nil
		}
		return helpMsg, nil
	case andyConfigShowYear.FullCommand():
		if *andyConfigShowYearArg != "" {
			return "No year config values to show for now :) ", nil
		}
		return helpMsg, nil
	case andyConfigShowMonth.FullCommand():
		if *andyConfigShowMonthArg != "" {
			CmdPreviewMonthAction(*andyConfigShowMonthArg)
			return "", nil
		}
		return helpMsg, nil
	case andyConfigShowWeek.FullCommand():
		if *andyConfigShowWeekArg != "" {
			CmdPreviewWeekAction(*andyConfigShowWeekArg)
			return "", nil
		}
		return helpMsg, nil
	case andyConfigShowDay.FullCommand():
		if *andyConfigShowDayArg != "" {
			CmdPreviewHourAction(*andyConfigShowDayArg)
			return "", nil
		}
		return helpMsg, nil
	case andyConfigClear.FullCommand():
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
	GenerateReadings(profile)
}

func CmdPreviewAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotReadingsChart(profile)
}

func CmdPreviewMonthAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotMonthlyProfilesChart(profile)
}

func CmdPreviewWeekAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotWeeklyProfilesChart(profile)
}

func CmdPreviewHourAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotHourlyProfilesChart(profile)
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
	parseableFileNamesList, _ := GetAppendedParsedFileNames()
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
			if resp.StatusCode == 200 || resp.StatusCode == 201 {
				SaveReadings(profile, defaultReadingsPath)
			} else {
				log.Println("Sent reading responded with a status other than 200 OK success", resp.StatusCode)
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
		profile, _ := GetProfileFromJson(readingsFile)
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

func GetAppendedParsedFileNames() (validFileNames []string, invalidFileNames []string) {
	parseableNamesList, unparseableNamesList := SplitParseableConfigFiles()
	os.MkdirAll(parseFolderPath, os.ModePerm)
	if _, err := os.Stat(parseableFileFullPath); os.IsNotExist(err) {
		_, err = os.Create(parseableFileFullPath)
	}
	parseableNamesBytes := EncodeFileNamesAsBytes(parseableNamesList)
	err := ioutil.WriteFile(parseableFileFullPath, parseableNamesBytes, 0644)
	checkWithPanic(err)

	if _, err := os.Stat(unparseableFileFullPath); os.IsNotExist(err) {
		_, err = os.Create(unparseableFileFullPath)
	}
	unparseableNamesBytes := EncodeFileNamesAsBytes(unparseableNamesList)
	err = ioutil.WriteFile(unparseableFileFullPath, unparseableNamesBytes, 0644)
	checkWithPanic(err)

	parseableFileBytes, err := ioutil.ReadFile(parseableFileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	parseableFileNamesList := DecodeBytesAsFileNames(parseableFileBytes)

	unparseableFileBytes, err := ioutil.ReadFile(unparseableFileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	unparseableFileNamesList := DecodeBytesAsFileNames(unparseableFileBytes)

	return parseableFileNamesList, unparseableFileNamesList
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
