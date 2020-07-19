package config

import (
	"os"
	"strconv"
)

type Config struct {
	YTApiKey             string
	MinioEndpoint        string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	MinioUseSSL          bool
	MinioBucketName      string
	MinioBucketRegion    string
	EsIndex              string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		YTApiKey:             getEnv("YTAPI_KEY", ""),
		MinioEndpoint:        getEnv("MINIO_ENDPOINT", ""),
		MinioAccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", ""),
		MinioSecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", ""),
		MinioUseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
		MinioBucketName:      getEnv("MINIO_BUCKET_NAME", "youtube"),
		MinioBucketRegion:    getEnv("MINIO_BUCKET_REGION", ""),
		EsIndex:              getEnv("ELASTICSEACH_INDEX", "youtube-video"),
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

// Simple helper function to read an environment variable into boolean or return a default value
func getEnvAsBool(name string, defaultValue bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultValue
}
