package main

var helpVersion = "v0.1"

var helpMsg = `andy config
Andy generates profiles for demonstration accounts
	andy config version								Show andy version
	andy config profile								Path to demonstration profile file
	andy config generate   							Generate a new profile configuration with the default name
	andy config generate   "my_new_config.json"		Generate a new profile configuration with the name
	andy config send   	   "readings_file.json"		Send the readings in the file to the API server
	andy config show year  "my_new_config.json"		Show the year set for the config
	andy config show month "my_new_config.json"		Show the month set for the config
	andy config show week  "my_new_config.json"		Show the week set for the config
	andy config show day   "my_new_config.json"		Show the day set for the config
 	andy config validate 							Validates all the configuration files in the profiles folder
	andy config validate   "my_new_config.json"		Validate the provided configuration file
	andy config init 								Create a default profile file `
