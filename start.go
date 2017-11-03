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

// Parse the commands passed to the generator from the command line
func ParseCLICommands() (string, error) {
	enteredCommand := kingpin.MustParse(app.Parse(os.Args[1:]))
	if len(enteredCommand) <= 2 {
		return helpMsg, nil
	}
	switch enteredCommand {
	// cmd: andy config version
	case andyVersion.FullCommand():
		return helpVersion, nil
	case andyConfigStart.FullCommand():
		InitGenerator()
		return "", nil
	// cmd: andy config generate "file.json"
	case andyConfigGenerate.FullCommand():
		if *andyConfigGenerateArg != "" {
			namesArr := strings.Split(*andyConfigGenerateArg, ".")
			extractedName := SanitizeName(namesArr[0])
			err := WriteProfileToFile(CreateDefaultProfile(extractedName), *andyConfigGenerateArg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return fmt.Sprintf("%s file created into ./profiles", *andyConfigGenerateArg), nil
		} else {
			return helpMsg, nil
		}
	// cmd: andy config init
	case andyConfigInit.FullCommand():
		err := WriteProfileToFile(CreateDefaultProfile(""), defaultProfileName)
		if err != nil {
			log.Fatal(err.Error())
		}
		return "Default profile file created into ./profiles", nil
	// cmd: andy config preview "file.json"
	// "q" to quit the preview
	case andyConfigPreview.FullCommand():
		fmt.Println("An ASCII art demonstration of the profile is starting")
		if *andyConfigPreviewArg != "" {
			CmdPreviewAction(*andyConfigPreviewArg)
		}
		CmdPreviewAction(defaultProfileName)
		return "", nil
	// cmdL andy config send "file.json"
	case andyConfigSend.FullCommand():
		if *andyConfigSendArg != "" {
			CmdSendReadingsToServer(*andyConfigSendArg)
			return "", nil
		}
		return helpMsg, nil
	// cmd: andy config validate "file.json"
	case andyConfigValidate.FullCommand():
		if *andyConfigValidateArg != "" {
			CmdValidateAction(*andyConfigValidateArg, true)
		}
		CmdValidateAction(*andyConfigValidateArg, false)
		return "", nil
	// cmd: andy config profile "file.json"
	case andyConfigProfile.FullCommand():
		if *andyConfigProfileArg != "" {
			CmdProfileAction(*andyConfigProfileArg)
			return "", nil
		}
		return helpMsg, nil
	// cmd: andy config show year "file.json"3=
	case andyConfigShowYear.FullCommand():
		if *andyConfigShowYearArg != "" {
			return "No year config values to show for now :) ", nil
		}
		return helpMsg, nil
	// cmd: andy config show month "file.json"
	case andyConfigShowMonth.FullCommand():
		if *andyConfigShowMonthArg != "" {
			CmdPreviewMonthAction(*andyConfigShowMonthArg)
			return "", nil
		}
		return helpMsg, nil
	// cmd: andy config show week "file.json"
	case andyConfigShowWeek.FullCommand():
		if *andyConfigShowWeekArg != "" {
			CmdPreviewWeekAction(*andyConfigShowWeekArg)
			return "", nil
		}
		return helpMsg, nil
	// cmd: andy config show day "file.json"
	case andyConfigShowDay.FullCommand():
		if *andyConfigShowDayArg != "" {
			CmdPreviewHourAction(*andyConfigShowDayArg)
			return "", nil
		}
		return helpMsg, nil
	// cmd: andy config clear
	case andyConfigClear.FullCommand():
		fmt.Println("Do you really want to clear the all saved configurations")
		return "", nil
	default:
		return helpMsg, fmt.Errorf("unknown command: %s", enteredCommand)
	}
	return "", nil
}

func CmdProfileAction(filename string) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	fmt.Println("reading json file")
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	GenerateReadings(profile)
}

// plot all the readings for the given file
func CmdPreviewAction(filename string) {
	// implement ASCII art readings here
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotReadingsChart(profile)
}

// plot all the readings for a given month given a file
func CmdPreviewMonthAction(filename string) {
	// implement ASCII art readings here
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotMonthlyProfilesChart(profile)
}

// plot all the readings for a given week given a file
func CmdPreviewWeekAction(filename string) {
	// implement ASCII art readings here
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotWeeklyProfilesChart(profile)
}

// plot all the readings for a given hour given a file
func CmdPreviewHourAction(filename string) {
	// implement ASCII art readings here
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	PlotHourlyProfilesChart(profile)
}

func CmdValidateSingleFile(filename string) {
	fmt.Printf("attempting to validate profile with name: %s", filename)
	fmt.Println()
	fileBytes, err := ioutil.ReadFile(filepath.Join(defaultProfilePath, filename))
	if err != nil {
		log.Fatal(err.Error())
	}
	profile, err := NewProfileFromJson(fileBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = profile.ValidateProfile()
	if err != nil {
		fmt.Println("The profile is not a valid profile")
		log.Fatal(err)
	}
	fmt.Println("####################################################")
	fmt.Println("------The profile configuration is valid------------")
	fmt.Println("####################################################")
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
	err = profile.ValidateProfile()
	if err != nil {
		fmt.Println("The profile is not a valid profile")
		log.Fatal(err)
	}
	librarianService := new(LibrarianService)
	// send the readings to the api at this point
	resp, err := librarianService.sendReadingsAction(profile)
	fmt.Println("SendReadingsAction() httpResp", resp)

	ping := func() { fmt.Println("#") }
	stop := sendReadingsSchedule(ping, 5*time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	stop <- true
	time.Sleep(25 * time.Millisecond)
}

// check all configurations and log the unparseable ones
// maintain two list of parseable and unparseable (with names)
// check if profile name exists in both files, if not, validate the file and
// append to parseable | unparseable list
func InitGenerator() {
	parseableFileNamesList, unparseableFileNamesList := GetAppendedParsedFileNames()
	fmt.Println(parseableFileNamesList)
	fmt.Println(unparseableFileNamesList)
}

// GetAppendedParsedFileNames gets the 2 slices each containing valid and invalid names of files
// It writes the list to a persistent file instead of keeping in memory
// Returns list of the saved invalid and valid config file names
func GetAppendedParsedFileNames() (validFileNames []string, invalidFileNames []string) {
	parseableNamesList, unparseableNamesList := SplitParseableConfigFiles()
	os.MkdirAll(parseFolderPath, os.ModePerm)
	if _, err := os.Stat(parseableFileFullPath); os.IsNotExist(err) {
		// does not exist
		_, err = os.Create(parseableFileFullPath)
	}
	parseableNamesBytes := EncodeFileNamesAsBytes(parseableNamesList)
	err := ioutil.WriteFile(parseableFileFullPath, parseableNamesBytes, 0644)
	checkWithPanic(err)

	if _, err := os.Stat(unparseableFileFullPath); os.IsNotExist(err) {
		// does not exist
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

// SplitParseableConfigFiles scans the profiles directory and for each json configuration file,
// it creates two slices
// > Slice 1 to store a list of the names of the config files which are valid
// > Slice 2 to store a list of the names of the config files which are invalid
func SplitParseableConfigFiles() (parseableList []string, unparseableList []string) {
	// return list of parseable files
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
		err = profile.ValidateProfile()
		if err != nil {
			// add file to unparseable list
			unparseableNamesList = append(unparseableNamesList, f.Name())
		} else {
			// add the filename to parseable
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
