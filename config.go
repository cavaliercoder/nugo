package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	FilePath     string        `yaml:"-"`
	ListenPort   string        `yaml:"listenPort"`
	BaseURL      string        `yaml:"baseUrl"`
	Repositories []*Repository `yaml:"repositories"`
}

var configPaths []string = []string{
	"./nugo.yaml",
	"/etc/nugo/nugo.yaml",
}

// config is a singleton cache for configuration loaded at start up
var config *Configuration = nil

// GetConfig returns runtime configuration for the Nugo server.
func GetConfig() *Configuration {
	if config == nil {
		// create default configuration
		config = &Configuration{
			ListenPort:   ":1105",
			Repositories: make([]*Repository, 0),
		}

		// check default paths for a configuration file
		var f *os.File = nil
		var err error
		for _, path := range configPaths {
			f, err = os.Open(path)
			if err == nil {
				config.FilePath = path
				break
			}
		}

		if f == nil {
			panic("No configuration file found")
		}

		defer f.Close()

		// read content
		b, err := ioutil.ReadAll(f)
		PanicOn(err)

		// decode yaml
		PanicOn(yaml.Unmarshal(b, config))

		// validate repos
		for _, r := range config.Repositories {
			if !strings.HasPrefix(r.RemotePath, "/") {
				r.RemotePath = "/" + r.RemotePath
			}

			if !strings.HasSuffix(r.RemotePath, "/") {
				r.RemotePath += "/"
			}
		}

		LogInfof("Loaded configuration file: %s", config.FilePath)
	}

	return config
}
