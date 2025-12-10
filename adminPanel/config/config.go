// Package config handles application configuration.
package config

import (
	"os"
	"strconv"
	"strings"
)

// Settings - конфигурация
type Settings struct {
	DatabaseURL          string
	DatabaseMinPoolSize  int
	DatabaseMaxPoolSize  int
	KeycloakIssuerURL    string
	KeycloakAudience     string
	KeycloakJWKSURL      string
	APIAddress           string
	Debug                bool
	CORSAllowOrigins     string
	CORSAllowCredentials bool
	ClientID             string
	ClientSecret         string
	AppName              string
}

// NewSettings создает конфигурацию
func NewSettings() *Settings {
	// Получаем CORS настройки
	corsOrigins := getEnv("CORS_ALLOW_ORIGINS", "http://localhost,http://localhost:4000,http://localhost:8080")
	corsCredentials := getEnvAsBool("CORS_ALLOW_CREDENTIALS", true)

	return &Settings{
		DatabaseURL:          getEnv("DATABASE_URL", "postgresql://appuser:password@app-db:5432/appdb?sslmode=disable"),
		DatabaseMinPoolSize:  getEnvAsInt("DATABASE_POOL_MIN_SIZE", 5),
		DatabaseMaxPoolSize:  getEnvAsInt("DATABASE_POOL_MAX_SIZE", 20),
		KeycloakIssuerURL:    os.Getenv("KEYCLOAK_ISSUER_URL"),
		KeycloakAudience:     os.Getenv("KEYCLOAK_AUDIENCE"),
		KeycloakJWKSURL:      os.Getenv("KEYCLOAK_JWKS_URL"),
		APIAddress:           getEnv("API_ADDRESS", ":4000"),
		Debug:                getEnvAsBool("DEBUG", false),
		CORSAllowOrigins:     corsOrigins,
		CORSAllowCredentials: corsCredentials,
		ClientID:             os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret:         os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		AppName:              os.Getenv("KEYCLOAK_APP_NAME"),
	}
}

// GetCORSOrigins возвращает origins как слайс
func (s *Settings) GetCORSOrigins() []string {
	if s.CORSAllowOrigins == "*" {
		return []string{"*"}
	}
	// Разделяем по запятой и убираем пробелы
	origins := strings.Split(s.CORSAllowOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

// Вспомогательные функции
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
