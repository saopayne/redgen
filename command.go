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

	// andy config init
	andyConfigStart = andyConfig.Command("start", "Start running the application")

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
	andyConfigSend    = andyConfig.Command("send", "A readings file should be sent along with this command")
	andyConfigSendArg = andyConfigSend.Arg("file_to_send.json", "Send the readings in the file specified to the server").String()

	// andy config show
	andyConfigShow = andyConfig.Command("show", "")

	andyConfigShowFileName  = andyConfigShow.Arg("filename", "Add the readings filename to show").String()
	andyConfigShowYearFlag  = andyConfigShow.Flag("year", "Add the year you wish to display consumption for").String()
	andyConfigShowMonthFlag = andyConfigShow.Flag("month", "Add the year-month you wish to display consumption for").String()
	andyConfigShowDayFlag   = andyConfigShow.Flag("day", "Add the year-month-day you wish to display consumption for").String()
)
