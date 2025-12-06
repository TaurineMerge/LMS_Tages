package main

import "os"

// ============ КОНФИГУРАЦИЯ ============

type Config struct {
	DatabaseURL string
}

func getConfig() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://appuser:password@app-db:5432/appdb?sslmode=disable"
	}

	return Config{
		DatabaseURL: dbURL,
	}
}
