package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const fileName = "gh-pr-links-configuration.json"

const BypassFilter = "-"

type Config struct {
	Style        string `json:"style"`
	UseEmoji     string `json:"useEmoji"`
	Organization string `json:"organization"`
}

func getConfigFileName() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, fileName)
}

func ReadConfig() (*Config, error) {
	data, err := os.ReadFile(getConfigFileName())
	if err != nil {
		if os.IsNotExist(err) {
			CreateDefaultConfig()
			return ReadConfig()
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if config.Organization == "" {
		config.Organization = "-"
		UpdateConfig(&config)
	}
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func UpdateConfig(newConfig *Config) error {
	data, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(getConfigFileName(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CreateDefaultConfig() {
	defaultConfig := &Config{
		Style:        "StyleRounded",
		UseEmoji:     "true",
		Organization: "-",
	}

	UpdateConfig(defaultConfig)
}
