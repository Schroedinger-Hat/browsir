package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName     string            `yaml:"app_name"`
	BrowserName string            `yaml:"browser_name"`
	Profiles    []Profile         `yaml:"profiles"`
	Shortcuts   map[string]string `yaml:"shortcuts"`
}

type Profile struct {
	Name        string `yaml:"name"`
	ProfileDir  string `yaml:"profile_dir"`
	Description string `yaml:"description"`
}

const CONFIG_PATH = "/.config/browsir"

func LoadConfig() (Config, error) {
	configPath, err := FindConfigFile()
	if err != nil {
		return Config{}, errors.New("config file not found")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, errors.New("error reading config file")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	if config.AppName == "" {
		config.AppName = "browsir"
	}
	if config.BrowserName == "" {
		config.BrowserName = "chrome"
	}
	if len(config.Profiles) == 0 {
		config.Profiles = []Profile{
			{Name: "default", ProfileDir: "Default", Description: "Default profile"},
		}
	}
	if config.Shortcuts == nil {
		config.Shortcuts = map[string]string{
			"cal": "calendar.google.com",
		}
	}

	return config, nil
}

func FindConfigFile() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = os.Getenv("HOME") + CONFIG_PATH
	} else {
		configHome = configHome + CONFIG_PATH
	}
	configPath := configHome + "/config.yml"

	_, err := os.Stat(configPath)
	if err != nil {
		return "", fmt.Errorf("config file not found, %v", err)
	}

	return configPath, nil
}

func GetConfigBase() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = os.Getenv("HOME") + CONFIG_PATH
	} else {
		configHome = configHome + CONFIG_PATH
	}

	_, err := os.Stat(configHome)
	if err != nil {
		return "", fmt.Errorf("config home not found, %v", err)
	}

	return configHome, nil
}
