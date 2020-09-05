package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	databaseURL string
	port        string
	initialized = false
)

// Init itializes the config variables of the API. It must be called
// before using any config variable.
func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	databaseURL = os.Getenv("DATABASE_URL")
	port = os.Getenv("PORT")
	initialized = true

	return nil
}

// DatabaseURL gets the url of the current database
func DatabaseURL() (string, error) {
	if !initialized {
		return "", errors.New("calling config variable without initialization")
	}

	return databaseURL, nil
}

// Port gets the port to listen to
func Port() (string, error) {
	if !initialized {
		return "", errors.New("calling config variable without initialization")
	}

	return ":" + port, nil
}
