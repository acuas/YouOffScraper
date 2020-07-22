package config

import (
	"os"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////

// Config contains environment variables used by this prject
type Config struct {
	// YouTube Data API key v3
	YTApiKey             string
	// Minio setting
	MinioEndpoint        string
	MinioAccessKeyID     string
	MinioSecretAccessKey string
	MinioUseSSL          bool
	MinioBucketName      string
	MinioBucketRegion    string
	// Elasticsearch Index
	EsIndex              string
	// fiber settings
	FiberPort            int
	NWorkers             int
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
		FiberPort:            getEnvAsInt("FIBER_PORT", 8000),
		NWorkers:             getEnvAsInt("N_WORKERS", 2),
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

///////////////////////////////////////////////////////////////////////////////
