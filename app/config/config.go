package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type MySqlDatabase struct {
	SqlHost     string `json:"sqlHost"`
	SqlPort     int    `json:"sqlPort"`
	SqlUser     string `json:"sqlUser"`
	SqlPassword string `json:"sqlPassword"`
	SqlDatabase string `json:"sqlDatabase"`
}

type Configuration struct {
	Environment string        `json:"environment"`
	Port        string        `json:"port"`
	Mysql       MySqlDatabase `json:"mysql"`
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
	sqlHost := os.Getenv("MYSQL_HOST")
	sqlPort := os.Getenv("MYSQL_PORT")
	sqlUser := os.Getenv("MYSQL_USER")
	sqlPassword := os.Getenv("MYSQL_PASSWORD")
	sqlDatabase := os.Getenv("MYSQL_DATABASE")
	intSqlPort, err := strconv.Atoi(sqlPort)
	if err != nil {
		log.Println("Error converting string to int", "warn")
		// log.Fatalf("Error converting string to int: %v", err)
	}
	// TODO: add additional errors to make sure all required environment values are provided

	database := MySqlDatabase{
		SqlHost:     sqlHost,
		SqlPort:     intSqlPort,
		SqlUser:     sqlUser,
		SqlPassword: sqlPassword,
		SqlDatabase: sqlDatabase,
	}

	configurationSingleton = &Config{
		Configuration: &Configuration{
			Environment: environment,
			Port:        port,
			Mysql:       database,
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
func GenerateConfiguration(config *Config) *Configuration {
	port := config.Configuration.Port
	environment := config.Configuration.Environment

	return &Configuration{
		Environment: environment,
		Port:        port,
	}
}
