// Package config отвечает за конфигурацию приложения
//
// Пакет предоставляет структуру Settings для хранения конфигурационных
// параметров и функции для их получения из переменных окружения.
//
// Поддерживаемые параметры:
//   - База данных: URL, размеры пулов
//   - Keycloak: URL, аудитория, JWKS
//   - API: адрес, режим отладки
//   - CORS: разрешенные origins, credentials
//   - OpenTelemetry: endpoint, протокол
package config

import (
	"os"
	"strconv"
	"strings"
)

// Settings - структура конфигурации приложения
//
// Содержит все необходимые параметры для работы Admin Panel:
//
// DatabaseURL - URL подключения к PostgreSQL
// DatabaseMinPoolSize - минимальный размер пула соединений
// DatabaseMaxPoolSize - максимальный размер пула соединений
// KeycloakIssuerURL - URL Keycloak для валидации токенов
// KeycloakAudience - аудитория для JWT-токенов
// KeycloakJWKSURL - URL для получения JWKS ключей
// APIAddress - адрес для запуска HTTP-сервера
// Debug - флаг режима отладки
// CORSAllowOrigins - разрешенные origins для CORS
// CORSAllowCredentials - разрешать ли credentials в CORS
// ClientID - ID клиента для OAuth
// ClientSecret - секрет клиента для OAuth
// AppName - название приложения для OAuth
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

// NewSettings создает и возвращает новую конфигурацию
//
// Функция читает параметры из переменных окружения, используя
// значения по умолчанию для неопределенных переменных.
//
// Возвращает:
//   - *Settings: указатель на структуру конфигурации
func NewSettings() *Settings {
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

// GetCORSOrigins возвращает список разрешенных origins для CORS
//
// Преобразует строку origins в слайс, разделяя по запятым
// и удаляя пробелы. Поддерживает wildcard "*".
//
// Возвращает:
//   - []string: слайс с origins
func (s *Settings) GetCORSOrigins() []string {
	if s.CORSAllowOrigins == "*" {
		return []string{"*"}
	}
	origins := strings.Split(s.CORSAllowOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

// Вспомогательные функции
// getEnv возвращает значение переменной окружения или значение по умолчанию
//
// Параметры:
//   - key: имя переменной окружения
//   - defaultValue: значение по умолчанию
//
// Возвращает:
//   - string: значение переменной окружения или defaultValue
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt возвращает значение переменной окружения как int или значение по умолчанию
//
// Параметры:
//   - key: имя переменной окружения
//   - defaultValue: значение по умолчанию
//
// Возвращает:
//   - int: значение переменной окружения или defaultValue
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool возвращает значение переменной окружения как bool или значение по умолчанию
//
// Параметры:
//   - key: имя переменной окружения
//   - defaultValue: значение по умолчанию
//
// Возвращает:
//   - bool: значение переменной окружения или defaultValue
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
