package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type app struct {
	AppID              string
	HTTPAddr           string
	HTTPServerLogLevel string
	ServiceLogLevel    string
	RepositoryLogLevel string
	HTTPAllowedOrigins []string
	HTTPAllowedMethods []string
	HTTPAllowedHeaders []string
}

var appConf *app

func GetApp() *app {
	return appConf
}

func LoadConfig() error {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	// Set up application configuration from environment variables
	appConf = &app{
		HTTPAddr:           os.Getenv("HTTP_ADDR"),
		HTTPServerLogLevel: os.Getenv("HTTP_SERVER_LOG_LEVEL"),
		ServiceLogLevel:    os.Getenv("SERVICE_LOG_LEVEL"),
		RepositoryLogLevel: os.Getenv("REPOSITORY_LOG_LEVEL"),
	}

	// Handle CORS settings, with default values if not set
	appConf.HTTPAllowedOrigins = parseCSV(os.Getenv("HTTP_ALLOWED_ORIGINS"), "*")
	appConf.HTTPAllowedMethods = parseCSV(os.Getenv("HTTP_ALLOWED_METHODS"), "GET,POST,PUT,DELETE,OPTIONS")
	appConf.HTTPAllowedHeaders = parseCSV(os.Getenv("HTTP_ALLOWED_HEADERS"), "*")

	return nil
}

// parseCSV takes a comma-separated string from environment variables and returns a slice of strings.
func parseCSV(value, defaultValue string) []string {
	if value == "" {
		return strings.Split(defaultValue, ",")
	}
	return strings.Split(value, ",")
}
