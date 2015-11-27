package main

type Configuration struct {
	PackagePath  string
	ListenPort   string
	BaseURL      string
	Repositories []*Repository
}

// config is a singleton cache for configuration loaded at start up
var config *Configuration = nil

// GetConfig returns runtime configuration for the Nugo server.
func GetConfig() *Configuration {
	if config == nil {
		// TODO: Load configuration from file
		config = &Configuration{
			PackagePath:  "packages",
			ListenPort:   ":1105",
			BaseURL:      "http://10.25.64.224:1105",
			Repositories: make([]*Repository, 0),
		}

		config.Repositories = append(config.Repositories, NewRepository(config.PackagePath))
		config.Repositories[0].GetPackages()
	}

	return config
}
