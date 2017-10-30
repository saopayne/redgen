package main

import "gopkg.in/alecthomas/kingpin.v2"

var (
	app        = kingpin.New("andy", "Andy random readings generator")
	andy       = app.Command("andy", "A random readings generator")
	andyConfig = andy.Command("config", "Work with a configuration")

	// andy config clear
	andyConfigClear = andyConfig.Command("clear", "Clears the profiles folder")

	// andy config version
	andyVersion = andyConfig.Command("version", "Display the version of andy")

	// andy config init
	andyConfigInit = andyConfig.Command("init", "Create a default profile.")

	// andy config generate "sample_file.json"
	andyConfigGenerate    = andyConfig.Command("generate", "Create a default profile.")
	andyConfigGenerateArg = andyConfigGenerate.Arg("file_to_generate.json", "Create a default profile into the provided file").String()

	// andy config preview "sample_file.json"
	andyConfigPreview    = andyConfig.Command("preview", "Preview the default profile.")
	andyConfigPreviewArg = andyConfigPreview.Arg("file_to_preview.json", "Preview the given profile.").String()

	// andy config preview "sample_file.json"
	andyConfigValidate    = andyConfig.Command("validate", "Validates all configurations")
	andyConfigValidateArg = andyConfigValidate.Arg("file_to_preview.json", "Validates the given configuration").String()

	// andy config profile "sample_file.json"
	andyConfigProfile    = andyConfig.Command("profile", "")
	andyConfigProfileArg = andyConfigProfile.Arg("profile.json", "Validates the given configuration").String()

	//andy config send
	andyConfigSend   = andyConfig.Command("send", "A readings file should be sent along with this command")
	andyConfigSendArg = andyConfigSend.Arg("file_to_send.json", "Send the readings in the file specified to the server").String()

	// andy config show
	andyConfigShow = andyConfig.Command("show", "")

	// andy config show year "sample_file.json"
	andyConfigShowYear    = andyConfigShow.Command("year", "")
	andyConfigShowYearArg = andyConfigShowYear.Arg("profile.json", "Displays the year for the provided configuration").String()

	// andy config show month "sample_file.json"
	andyConfigShowMonth    = andyConfigShow.Command("month", "")
	andyConfigShowMonthArg = andyConfigShowMonth.Arg("profile.json", "Displays the month for the provided configuration").String()

	// andy config show week "sample_file.json"
	andyConfigShowWeek    = andyConfigShow.Command("week", "")
	andyConfigShowWeekArg = andyConfigShowWeek.Arg("profile.json", "Displays the week for the provided configuration").String()

	// andy config show day "sample_file.json"
	andyConfigShowDay    = andyConfigShow.Command("day", "")
	andyConfigShowDayArg = andyConfigShowDay.Arg("profile.json", "Displays the day for the provided configuration").String()
)
