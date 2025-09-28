package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	Environment string `json:"environment"`
	Port        string `json:"port"`
}

type Config struct {
	Configuration *Configuration
}

var configurationSingleton *Config

/*
Build and return the environmental configuration.

Returns Configuration or error, if issues occur.
*/
func BuildConfiguration() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	port := os.Getenv("PORT")
	environment := os.Getenv("ENV")
	configurationSingleton = &Config{
		Configuration: &Configuration{
			Environment: environment,
			Port:        port,
		},
	}

	return configurationSingleton, nil
}

func GetConfiguration() (*Config, error) {
	if configurationSingleton == nil {
		build, e := BuildConfiguration()
		if e != nil {
			return nil, errors.New("project configuration error")
		}

		return build, nil
	}

	return configurationSingleton, nil
}

// Generate a custom configuration object based on variables passed.
// This function does not effect the configuration singleton
func GenerateConfiguration(config *Config) *Config {
	var generatedConfiguration *Config
	port := config.Configuration.Port
	environment := config.Configuration.Environment

	generatedConfiguration = &Config{
		Configuration: &Configuration{
			Environment: environment,
			Port:        port,
		},
	}

	return generatedConfiguration
}
