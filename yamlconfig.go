package main

import (
	"io"
	"log"
	"os"

	"asterisk-dialer/asterisk"
	"asterisk-dialer/config"

	"gopkg.in/yaml.v2"
)

// Parser must implement ParseYAML
type Parser interface {
	ParseYAML([]byte) error
}

// Load the YAML config file
func configLoad(configFile string, p Parser) {
	var err error
	var input = io.ReadCloser(os.Stdin)
	if input, err = os.Open(configFile); err != nil {
		log.Fatalln(err)
	}

	// Read the config file
	yamlBytes, err := io.ReadAll(input)
	input.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// Parse the config
	if err := p.ParseYAML(yamlBytes); err != nil {
		//log.Fatalf("Content: %v", yamlBytes)
		log.Fatalf("Could not parse %q: %v", configFile, err)
	}
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

type (
	// configuration contains the application settings
	conf struct {
		Listen string `yaml:"listen"`

		Asterisk asterisk.Ami `yaml:"asterisk"`

		Config *config.Api `yaml:"api"`
	}
)

// ParseYAML unmarshals bytes to structs
func (c *conf) ParseYAML(b []byte) error {
	return yaml.Unmarshal(b, &c)
}

// Make config
func getConfig() {
	// Load the configuration file
	if *config_file == "" {
		*config_file = "config" + string(os.PathSeparator) + "config.yml"
	}
	configLoad(*config_file, cnf)
}
