package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func init() {
	LoadConfig()
}

func GetConfig() (*Config, error) {
	return LoadConfig()
}

// TODO Load config should set Config object using Concurrent object creator
// Enhance Load config to Set Env variables along with Config reader
// Load config Should be part of GETSETGO where we can simply call
// function and pass struct no worry about how it's goona set the context
func LoadConfig() (*Config, error) {
	// Load the configuration here
	// =========================================================================
	// Configuration from file
	// Set the file name of the configurations file
	viper.SetConfigName(func() string {
		if os.Getenv("CONFIG_FILE") != "" {
			return os.Getenv("CONFIG_FILE")
		}
		return "dev" // Default to dev.yaml
	}())
	// Set the path to look for the configurations file
	viper.AddConfigPath(func() string {
		if os.Getenv("CONFIG_PATH") != "" {
			return os.Getenv("CONFIG_PATH")
		}
		return "." // Default to the current directory
	}())
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	cfg := Config{}
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Error reading config file")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "Unable to decode into struct")
	}
	return &cfg, nil
}
