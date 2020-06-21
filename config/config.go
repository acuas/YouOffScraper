package config

import (
	"os"
	"strconv"
)

type Config struct {
	YTApiKey            string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		YTApiKey:            getEnv("YTAPI_KEY", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}