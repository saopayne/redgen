# RedGen

RedGen is a random readings generator and profiler

##### To Run

`sh build.run`

&&

```
RedGen generates profiles for demonstration accounts

	config start							                Starts running the app, validate and send readings
	config version				                            Show andy version
	config preview "filename.json"			                Show sample consumptions for the date
	config profile			                         		Path to demonstration profile file
	config generate   				                        Generate a new profile configuration with the default name
	config generate "my_new_config.json"			   	  	Generate a new profile configuration with the name
	config send "readings_file.json"	               		Send the readings in the file to the API server
	config show "readings_file.json" --date=2017-03-01		Show the consumption for the particular day
    config validate 				                            Validates all the configuration files in the profiles folder
	config validate "my_new_config.json"		           	Validate the provided configuration file
	config init 			                             	Create a default profile file 

```



