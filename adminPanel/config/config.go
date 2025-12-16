// Package config отвечает за конфигурацию приложения
//
// Пакет предоставляет структуру Settings для хранения конфигурационных
// параметров и функции для их получения из переменных окружения.
//
// Поддерживаемые параметры:
//   - База данных: URL или отдельные параметры (host, port, user, password, name)
//   - Keycloak: URL, аудитория, JWKS
//   - API: адрес, режим отладки
//   - CORS: разрешенные origins, credentials
//   - OpenTelemetry: endpoint, протокол
//
// Пример использования переменных окружения:
//
//	# Вариант 1: Полный URL подключения
//	DATABASE_URL=postgresql://user:pass@host:5432/dbname?sslmode=disable
//
//	# Вариант 2: Отдельные параметры (если DATABASE_URL не задан)
//	DB_HOST=localhost
//	DB_PORT=5432
//	DB_USER=appuser
//	DB_PASSWORD=password
//	DB_NAME=appdb
//	DB_SSLMODE=disable
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DatabaseConfig - конфигурация подключения к базе данных
//
// Позволяет задавать параметры подключения как через DATABASE_URL,
// так и через отдельные переменные окружения.
type DatabaseConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Name        string
	SSLMode     string
	MinPoolSize int
	MaxPoolSize int
}

// URL возвращает строку подключения к базе данных
func (d *DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// OTelConfig - конфигурация OpenTelemetry
type OTelConfig struct {
	Endpoint    string
	ServiceName string
	Protocol    string
	Enabled     bool
}

// KeycloakConfig - конфигурация Keycloak для аутентификации
type KeycloakConfig struct {
	IssuerURL    string
	Audience     string
	JWKSURL      string
	ClientID     string
	ClientSecret string
	AppName      string
}

// CORSConfig - конфигурация CORS
type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

// ServerConfig - конфигурация HTTP сервера
type ServerConfig struct {
	Address  string
	AppName  string
	RootPath string
}

// Settings - структура конфигурации приложения
//
// Содержит все необходимые параметры для работы Admin Panel,
// организованные в логические группы.
type Settings struct {
	Database DatabaseConfig
	OTel     OTelConfig
	Keycloak KeycloakConfig
	CORS     CORSConfig
	Server   ServerConfig
	Debug    bool
}

// Validate проверяет обязательные параметры конфигурации
//
// Возвращает:
//   - error: ошибка валидации (если есть)
func (s *Settings) Validate() error {
	var missingVars []string

	// Проверяем конфигурацию БД
	if s.Database.Host == "" {
		missingVars = append(missingVars, "DB_HOST (or DATABASE_URL)")
	}
	if s.Database.User == "" {
		missingVars = append(missingVars, "DB_USER (or DATABASE_URL)")
	}
	if s.Database.Password == "" {
		missingVars = append(missingVars, "DB_PASSWORD (or DATABASE_URL)")
	}
	if s.Database.Name == "" {
		missingVars = append(missingVars, "DB_NAME (or DATABASE_URL)")
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missingVars, ", "))
	}

	return nil
}

// NewSettings создает и возвращает новую конфигурацию
//
// Функция читает параметры из переменных окружения.
// Приоритет для базы данных: DATABASE_URL > отдельные DB_* переменные
//
// Возвращает:
//   - *Settings: указатель на структуру конфигурации
func NewSettings() *Settings {
	return &Settings{
		Database: loadDatabaseConfig(),
		OTel:     loadOTelConfig(),
		Keycloak: loadKeycloakConfig(),
		CORS:     loadCORSConfig(),
		Server:   loadServerConfig(),
		Debug:    getEnvAsBool("DEBUG", false),
	}
}

// loadDatabaseConfig загружает конфигурацию базы данных
func loadDatabaseConfig() DatabaseConfig {
	cfg := DatabaseConfig{
		MinPoolSize: getEnvAsInt("DATABASE_POOL_MIN_SIZE", 5),
		MaxPoolSize: getEnvAsInt("DATABASE_POOL_MAX_SIZE", 20),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
	}

	// Приоритет: DATABASE_URL > отдельные переменные
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// Парсим DATABASE_URL для заполнения структуры
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode = parseDatabaseURL(databaseURL)
	} else {
		// Используем отдельные переменные
		cfg.Host = os.Getenv("DB_HOST")
		cfg.Port = getEnvAsInt("DB_PORT", 5432)
		cfg.User = os.Getenv("DB_USER")
		cfg.Password = os.Getenv("DB_PASSWORD")
		cfg.Name = os.Getenv("DB_NAME")
	}

	return cfg
}

// parseDatabaseURL парсит DATABASE_URL и возвращает компоненты
func parseDatabaseURL(url string) (host string, port int, user, password, name, sslmode string) {
	// Формат: postgresql://user:password@host:port/dbname?sslmode=disable
	port = 5432
	sslmode = "disable"

	// Убираем протокол
	url = strings.TrimPrefix(url, "postgresql://")
	url = strings.TrimPrefix(url, "postgres://")

	// Разделяем на credentials и host
	if atIdx := strings.LastIndex(url, "@"); atIdx != -1 {
		credentials := url[:atIdx]
		rest := url[atIdx+1:]

		// Парсим credentials
		if colonIdx := strings.Index(credentials, ":"); colonIdx != -1 {
			user = credentials[:colonIdx]
			password = credentials[colonIdx+1:]
		} else {
			user = credentials
		}

		// Парсим host:port/dbname?params
		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			hostPort := rest[:slashIdx]
			dbAndParams := rest[slashIdx+1:]

			// Парсим host:port
			if colonIdx := strings.Index(hostPort, ":"); colonIdx != -1 {
				host = hostPort[:colonIdx]
				port = getEnvAsInt("", 5432)
				if p, err := strconv.Atoi(hostPort[colonIdx+1:]); err == nil {
					port = p
				}
			} else {
				host = hostPort
			}

			// Парсим dbname?params
			if qIdx := strings.Index(dbAndParams, "?"); qIdx != -1 {
				name = dbAndParams[:qIdx]
				params := dbAndParams[qIdx+1:]
				// Парсим sslmode из параметров
				for _, param := range strings.Split(params, "&") {
					if strings.HasPrefix(param, "sslmode=") {
						sslmode = strings.TrimPrefix(param, "sslmode=")
					}
				}
			} else {
				name = dbAndParams
			}
		} else {
			host = rest
		}
	}

	return
}

// loadOTelConfig загружает конфигурацию OpenTelemetry
func loadOTelConfig() OTelConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	return OTelConfig{
		Endpoint:    endpoint,
		ServiceName: getEnv("OTEL_SERVICE_NAME", "admin-panel"),
		Protocol:    getEnv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"),
		Enabled:     endpoint != "",
	}
}

// loadKeycloakConfig загружает конфигурацию Keycloak
func loadKeycloakConfig() KeycloakConfig {
	issuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL")
	if jwksURL == "" && issuer != "" {
		jwksURL = strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
	}

	return KeycloakConfig{
		IssuerURL:    issuer,
		Audience:     os.Getenv("KEYCLOAK_AUDIENCE"),
		JWKSURL:      jwksURL,
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		AppName:      os.Getenv("KEYCLOAK_APP_NAME"),
	}
}

// loadCORSConfig загружает конфигурацию CORS
func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		AllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		AllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
		ExposeHeaders:    getEnv("CORS_EXPOSE_HEADERS", "Content-Length"),
	}
}

// loadServerConfig загружает конфигурацию сервера
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Address:  getEnv("API_ADDRESS", ":4000"),
		AppName:  getEnv("APP_NAME", "Admin Panel API"),
		RootPath: getEnv("ROOT_PATH", "/admin"),
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
	if s.CORS.AllowOrigins == "*" {
		return []string{"*"}
	}
	origins := strings.Split(s.CORS.AllowOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

// DatabaseURL возвращает URL подключения к базе данных
// Deprecated: используйте Settings.Database.URL()
func (s *Settings) DatabaseURL() string {
	return s.Database.URL()
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
