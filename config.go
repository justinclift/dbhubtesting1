package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
)

var (
	// Conf holds our configuration info
	Conf TomlConfig
)

func ReadConfig() error {
	// Override config file location via environment variable
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		userHome, err := homedir.Dir()
		if err != nil {
			log.Fatalf("User home directory couldn't be determined: %s", "\n")
		}
		configFile = filepath.Join(userHome, ".dbhub", "apiclient.toml")
	}

	// Read the server configuration from disk
	if _, err := toml.DecodeFile(configFile, &Conf); err != nil {
		return fmt.Errorf("Config file couldn't be parsed: %v\n", err)
	}

	// Verify we have the needed configuration information
	var missingConfig []string
	if Conf.Api.APIKey == "" {
		missingConfig = append(missingConfig, "API key")
	}
	if len(missingConfig) > 0 {
		// Some config is missing
		returnMessage := fmt.Sprintf("The configuration file (%s) is missing these configuration "+
			"value(s):", configFile)
		for _, value := range missingConfig {
			returnMessage += fmt.Sprintf("\n\t* %v", value)
		}
		return fmt.Errorf(returnMessage)
	}

	// The configuration file seems good
	return nil
}
