package main

type Configuration struct {
	PackagePath string
	ListenPort  string
}

// config is a singleton cache for configuration loaded at start up
var config *Configuration = nil

// GetConfig returns runtime configuration for the Nugo server.
func GetConfig() *Configuration {
	if config == nil {
		// TODO: Load configuration from file
		config = &Configuration{
			PackagePath: "packages",
			ListenPort:  ":1105",
		}
	}

	return config
}
