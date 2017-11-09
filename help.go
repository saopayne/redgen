package main

var helpVersion = "v0.1"

var helpMsg = `andy config
Andy generates profiles for demonstration accounts
	andy config start								Starts running the app, validate and send readings
	andy config version								Show andy version
	andy config profile								Path to demonstration profile file
	andy config generate   							Generate a new profile configuration with the default name
	andy config generate   "my_new_config.json"		Generate a new profile configuration with the name
	andy config send   	   "readings_file.json"		Send the readings in the file to the API server
	andy config show year  "readings_file.json"		Show the year consumption for the readings
	andy config show month "readings_file.json"		Show the monthly consumption for the config
	andy config show day   "readings_file.json"		Show the consumption for the particular day
 	andy config validate 							Validates all the configuration files in the profiles folder
	andy config validate   "my_new_config.json"		Validate the provided configuration file
	andy config init 								Create a default profile file `
