package main

import "gopkg.in/alecthomas/kingpin.v2"

var (
	app        = kingpin.New("redgen", "RedGen random readings generator")
	redgenConfig = app.Command("config", "Work with a configuration")

	// config clear
	redgenConfigClear = redgenConfig.Command("clear", "Clears the profiles folder")

	// config version
	redgenVersion = redgenConfig.Command("version", "Display the version of redgen")

	// config init
	redgenConfigInit = redgenConfig.Command("init", "Create a default profile.")

	// config start
	redgenConfigStart = redgenConfig.Command("start", "Start running the application")

	// config generate "sample_file.json"
	redgenConfigGenerate    = redgenConfig.Command("generate", "Create a default profile.")
	redgenConfigGenerateArg = redgenConfigGenerate.Arg("file_to_generate.json", "Create a default profile into the provided file").String()

	// config preview "sample_file.json"
	redgenConfigPreview        = redgenConfig.Command("preview", "Preview the default profile.")
	redgenConfigPreviewArg     = redgenConfigPreview.Arg("file_to_preview.json", "Preview the given profile.").String()
	redgenConfigPreviewArgTime = redgenConfigPreview.Flag("time", "Preview the given profile for a year|month|day.").String()

	// config validate "sample_file.json"
	redgenConfigValidate    = redgenConfig.Command("validate", "Validates all configurations")
	redgenConfigValidateArg = redgenConfigValidate.Arg("file_to_preview.json", "Validates the given configuration").String()

	// config profile "sample_file.json"
	redgenConfigProfile    = redgenConfig.Command("profile", "")
	redgenConfigProfileArg = redgenConfigProfile.Arg("profile.json", "Validates the given configuration").String()

	// config send
	redgenConfigSend    = redgenConfig.Command("send", "A readings file should be sent along with this command")
	redgenConfigSendArg = redgenConfigSend.Arg("file_to_send.json", "Send the readings in the file specified to the server").String()

	// config show
	redgenConfigShow = redgenConfig.Command("show", "")

	redgenConfigShowFileName = redgenConfigShow.Arg("filename", "Add the readings filename to show").String()
	redgenConfigShowDateFlag = redgenConfigShow.Flag("date", "Add the year-month-day you wish to display consumption for").String()
)
