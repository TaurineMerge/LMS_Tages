package config

import "os"

// Config stores all configuration of the application.
type Config struct {
	DatabaseURL string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://appuser:password@localhost:5432/appdb?sslmode=disable"),
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
