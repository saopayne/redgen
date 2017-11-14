package main

import "gopkg.in/alecthomas/kingpin.v2"

var (
	app        = kingpin.New("andy", "Andy random readings generator")
	andyConfig = app.Command("config", "Work with a configuration")

	// config clear
	andyConfigClear = andyConfig.Command("clear", "Clears the profiles folder")

	// config version
	andyVersion = andyConfig.Command("version", "Display the version of andy")

	// config init
	andyConfigInit = andyConfig.Command("init", "Create a default profile.")

	// config start
	andyConfigStart = andyConfig.Command("start", "Start running the application")

	// config generate "sample_file.json"
	andyConfigGenerate    = andyConfig.Command("generate", "Create a default profile.")
	andyConfigGenerateArg = andyConfigGenerate.Arg("file_to_generate.json", "Create a default profile into the provided file").String()

	// config preview "sample_file.json"
	andyConfigPreview    = andyConfig.Command("preview", "Preview the default profile.")
	andyConfigPreviewArg = andyConfigPreview.Arg("file_to_preview.json", "Preview the given profile.").String()

	// config preview "sample_file.json"
	andyConfigValidate    = andyConfig.Command("validate", "Validates all configurations")
	andyConfigValidateArg = andyConfigValidate.Arg("file_to_preview.json", "Validates the given configuration").String()

	// config profile "sample_file.json"
	andyConfigProfile    = andyConfig.Command("profile", "")
	andyConfigProfileArg = andyConfigProfile.Arg("profile.json", "Validates the given configuration").String()

	// config send
	andyConfigSend    = andyConfig.Command("send", "A readings file should be sent along with this command")
	andyConfigSendArg = andyConfigSend.Arg("file_to_send.json", "Send the readings in the file specified to the server").String()

	// config show
	andyConfigShow = andyConfig.Command("show", "")

	andyConfigShowFileName = andyConfigShow.Arg("filename", "Add the readings filename to show").String()
	andyConfigShowDateFlag = andyConfigShow.Flag("date", "Add the year-month-day you wish to display consumption for").String()
)
