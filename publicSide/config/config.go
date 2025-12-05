package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Загружаем .env файл если он существует (не критично если его нет)
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://appuser:password@app-db:5432/appdb?sslmode=disable"),
	}

	return cfg
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

