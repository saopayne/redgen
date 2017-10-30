package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

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
	// cmd: andy config generate "file.json"
	case andyConfigGenerate.FullCommand():
		if *andyConfigGenerateArg != "" {
			err := WriteProfileToFile(CreateDefaultProfile(), *andyConfigGenerateArg)
			if err != nil {
				log.Fatal(err.Error())
			}
			return fmt.Sprintf("%s file created into ./profiles", *andyConfigGenerateArg), nil
		} else {
			return helpMsg, nil
		}
	// cmd: andy config init
	case andyConfigInit.FullCommand():
		err := WriteProfileToFile(CreateDefaultProfile(), defaultProfileName)
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
	profileService := new(ProfileService)
	// send the readings to the api at this point
	resp, err := profileService.sendReadingsAction(profile)
	fmt.Println("SendReadingsAction() httpResp", resp)

	ping := func() { fmt.Println("#") }
	stop := sendReadingsSchedule(ping, 5*time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	stop <- true
	time.Sleep(25 * time.Millisecond)

}