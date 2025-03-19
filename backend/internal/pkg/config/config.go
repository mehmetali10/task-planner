package config

import (
	"fmt"
	"os"
	"strconv"
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

	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
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

	// Load database configuration
	appConf.DBHost = os.Getenv("DB_HOST")
	appConf.DBUser = os.Getenv("DB_USER")
	appConf.DBPassword = os.Getenv("DB_PASSWORD")
	appConf.DBName = os.Getenv("DB_NAME")

	// Parse DBPort as an integer, with a default value if not set or invalid
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		appConf.DBPort = 5432 // Default port for PostgreSQL
	} else {
		dbPort, err := strconv.Atoi(dbPortStr)
		if err != nil {
			return fmt.Errorf("invalid DB_PORT value: %v", err)
		}
		appConf.DBPort = dbPort
	}

	return nil
}

// parseCSV takes a comma-separated string from environment variables and returns a slice of strings.
func parseCSV(value, defaultValue string) []string {
	if value == "" {
		return strings.Split(defaultValue, ",")
	}
	return strings.Split(value, ",")
}
