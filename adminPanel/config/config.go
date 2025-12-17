// Package config предоставляет управление конфигурацией для Admin Panel.
//
// Пакет поддерживает загрузку конфигурации из переменных окружения с разумными
// значениями по умолчанию и валидацией. Поддерживаются два способа конфигурации БД:
//
//   - DATABASE_URL: полная строка подключения (имеет приоритет)
//   - Отдельные DB_* переменные: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
//
// # Переменные окружения
//
// База данных:
//
//	DATABASE_URL           - Полная строка подключения PostgreSQL
//	DB_HOST                - Хост базы данных (обязательно, если DATABASE_URL не задан)
//	DB_PORT                - Порт базы данных (по умолчанию: 5432)
//	DB_USER                - Пользователь БД (обязательно, если DATABASE_URL не задан)
//	DB_PASSWORD            - Пароль БД (обязательно, если DATABASE_URL не задан)
//	DB_NAME                - Имя базы данных (обязательно, если DATABASE_URL не задан)
//	DB_SSLMODE             - Режим SSL (по умолчанию: disable)
//	DATABASE_POOL_MIN_SIZE - Минимальный размер пула соединений (по умолчанию: 5)
//	DATABASE_POOL_MAX_SIZE - Максимальный размер пула соединений (по умолчанию: 20)
//
// Сервер:
//
//	API_ADDRESS - Адрес для прослушивания (по умолчанию: :4000)
//	APP_NAME    - Название приложения (по умолчанию: Admin Panel API)
//	ROOT_PATH   - Корневой путь URL (по умолчанию: /admin)
//	DEBUG       - Режим отладки (по умолчанию: false)
//
// Keycloak:
//
//	KEYCLOAK_ISSUER_URL    - URL издателя токенов для валидации
//	KEYCLOAK_AUDIENCE      - Ожидаемая аудитория JWT
//	KEYCLOAK_JWKS_URL      - URL эндпоинта JWKS (генерируется автоматически из issuer)
//	KEYCLOAK_CLIENT_ID     - OAuth client ID
//	KEYCLOAK_CLIENT_SECRET - OAuth client secret
//	KEYCLOAK_APP_NAME      - Название приложения для OAuth
//
// CORS:
//
//	CORS_ALLOW_ORIGINS     - Разрешённые origins (по умолчанию: *)
//	CORS_ALLOW_METHODS     - Разрешённые методы (по умолчанию: GET,POST,PUT,DELETE,OPTIONS)
//	CORS_ALLOW_HEADERS     - Разрешённые заголовки
//	CORS_ALLOW_CREDENTIALS - Разрешить credentials (по умолчанию: false)
//	CORS_EXPOSE_HEADERS    - Экспонируемые заголовки (по умолчанию: Content-Length)
//
// OpenTelemetry:
//
//	OTEL_EXPORTER_OTLP_ENDPOINT - Эндпоинт OTLP коллектора
//	OTEL_SERVICE_NAME           - Имя сервиса (по умолчанию: admin-panel)
//	OTEL_EXPORTER_OTLP_PROTOCOL - Протокол OTLP (по умолчанию: grpc)
//
// # Использование
//
//	settings := config.NewSettings()
//	if err := settings.Validate(); err != nil {
//	    log.Fatalf("Ошибка конфигурации: %v", err)
//	}
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DatabaseConfig содержит настройки подключения к PostgreSQL.
//
// Конфигурация может быть предоставлена либо через переменную окружения DATABASE_URL,
// либо через отдельные DB_* переменные. Когда DATABASE_URL задан, он имеет приоритет
// над отдельными переменными.
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

// URL возвращает строку подключения к PostgreSQL в стандартном формате.
//
// Формат: postgresql://user:password@host:port/dbname?sslmode=value
func (d *DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// OTelConfig содержит настройки трассировки OpenTelemetry.
type OTelConfig struct {
	Endpoint    string
	ServiceName string
	Protocol    string
	Enabled     bool
}

// KeycloakConfig содержит настройки аутентификации через Keycloak.
type KeycloakConfig struct {
	IssuerURL    string
	Audience     string
	JWKSURL      string
	ClientID     string
	ClientSecret string
	AppName      string
}

// CORSConfig содержит настройки Cross-Origin Resource Sharing.
type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

// ServerConfig содержит настройки HTTP сервера.
type ServerConfig struct {
	Address  string
	AppName  string
	RootPath string
}

// Settings является главным контейнером конфигурации приложения.
//
// Объединяет все подсекции конфигурации и предоставляет методы
// валидации и доступа к значениям конфигурации.
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

// loadDatabaseConfig загружает конфигурацию базы данных из переменных окружения.
//
// Функция читает настройки подключения к PostgreSQL. Поддерживает два способа:
//   - DATABASE_URL: полная строка подключения (приоритет)
//   - Отдельные переменные: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
//
// Дополнительно загружает настройки пула соединений:
//   - DATABASE_POOL_MIN_SIZE: минимальное количество соединений (по умолчанию: 5)
//   - DATABASE_POOL_MAX_SIZE: максимальное количество соединений (по умолчанию: 20)
//
// Возвращает:
//   - DatabaseConfig: структура с настройками подключения к БД
func loadDatabaseConfig() DatabaseConfig {
	cfg := DatabaseConfig{
		MinPoolSize: getEnvAsInt("DATABASE_POOL_MIN_SIZE", 5),
		MaxPoolSize: getEnvAsInt("DATABASE_POOL_MAX_SIZE", 20),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
	}

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode = parseDatabaseURL(databaseURL)
	} else {
		cfg.Host = os.Getenv("DB_HOST")
		cfg.Port = getEnvAsInt("DB_PORT", 5432)
		cfg.User = os.Getenv("DB_USER")
		cfg.Password = os.Getenv("DB_PASSWORD")
		cfg.Name = os.Getenv("DB_NAME")
	}

	return cfg
}

// parseDatabaseURL разбирает строку подключения DATABASE_URL на отдельные компоненты.
//
// Поддерживаемые форматы URL:
//   - postgresql://user:password@host:port/dbname?sslmode=disable
//   - postgres://user:password@host:port/dbname
//
// Параметры:
//   - url: строка подключения к PostgreSQL
//
// Возвращает:
//   - host: имя хоста сервера БД
//   - port: порт сервера БД (по умолчанию 5432)
//   - user: имя пользователя
//   - password: пароль пользователя
//   - name: имя базы данных
//   - sslmode: режим SSL (по умолчанию "disable")
func parseDatabaseURL(url string) (host string, port int, user, password, name, sslmode string) {
	port = 5432
	sslmode = "disable"

	url = strings.TrimPrefix(url, "postgresql://")
	url = strings.TrimPrefix(url, "postgres://")

	if atIdx := strings.LastIndex(url, "@"); atIdx != -1 {
		credentials := url[:atIdx]
		rest := url[atIdx+1:]

		if colonIdx := strings.Index(credentials, ":"); colonIdx != -1 {
			user = credentials[:colonIdx]
			password = credentials[colonIdx+1:]
		} else {
			user = credentials
		}

		if slashIdx := strings.Index(rest, "/"); slashIdx != -1 {
			hostPort := rest[:slashIdx]
			dbAndParams := rest[slashIdx+1:]

			if colonIdx := strings.Index(hostPort, ":"); colonIdx != -1 {
				host = hostPort[:colonIdx]
				port = getEnvAsInt("", 5432)
				if p, err := strconv.Atoi(hostPort[colonIdx+1:]); err == nil {
					port = p
				}
			} else {
				host = hostPort
			}

			if qIdx := strings.Index(dbAndParams, "?"); qIdx != -1 {
				name = dbAndParams[:qIdx]
				params := dbAndParams[qIdx+1:]
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

// loadOTelConfig загружает конфигурацию трассировки OpenTelemetry из переменных окружения.
//
// Читаемые переменные:
//   - OTEL_EXPORTER_OTLP_ENDPOINT: URL эндпоинта OTLP коллектора
//   - OTEL_SERVICE_NAME: имя сервиса для трасс (по умолчанию: "admin-panel")
//   - OTEL_EXPORTER_OTLP_PROTOCOL: протокол экспорта (по умолчанию: "grpc")
//
// Трассировка автоматически включается, если OTEL_EXPORTER_OTLP_ENDPOINT задан.
//
// Возвращает:
//   - OTelConfig: структура с настройками OpenTelemetry
func loadOTelConfig() OTelConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	return OTelConfig{
		Endpoint:    endpoint,
		ServiceName: getEnv("OTEL_SERVICE_NAME", "admin-panel"),
		Protocol:    getEnv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"),
		Enabled:     endpoint != "",
	}
}

// loadKeycloakConfig загружает конфигурацию аутентификации Keycloak из переменных окружения.
//
// Читаемые переменные:
//   - KEYCLOAK_ISSUER_URL: URL издателя токенов для валидации JWT
//   - KEYCLOAK_AUDIENCE: ожидаемый claim audience в JWT-токенах
//   - KEYCLOAK_JWKS_URL: URL эндпоинта JWKS (если не задан, генерируется из issuer)
//   - KEYCLOAK_CLIENT_ID: OAuth2 client ID для авторизации
//   - KEYCLOAK_CLIENT_SECRET: OAuth2 client secret
//   - KEYCLOAK_APP_NAME: название приложения для OAuth2
//
// Если KEYCLOAK_JWKS_URL не задан, он автоматически формируется как:
// {KEYCLOAK_ISSUER_URL}/protocol/openid-connect/certs
//
// Возвращает:
//   - KeycloakConfig: структура с настройками Keycloak
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

// loadCORSConfig загружает конфигурацию CORS из переменных окружения.
//
// Читаемые переменные:
//   - CORS_ALLOW_ORIGINS: разрешённые origins через запятую (по умолчанию: "*")
//   - CORS_ALLOW_METHODS: разрешённые HTTP методы (по умолчанию: "GET,POST,PUT,DELETE,OPTIONS")
//   - CORS_ALLOW_HEADERS: разрешённые заголовки (по умолчанию: "Origin,Content-Type,Accept,Authorization")
//   - CORS_ALLOW_CREDENTIALS: разрешить передачу credentials (по умолчанию: false)
//   - CORS_EXPOSE_HEADERS: экспонируемые заголовки (по умолчанию: "Content-Length")
//
// Возвращает:
//   - CORSConfig: структура с настройками CORS
func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		AllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		AllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
		ExposeHeaders:    getEnv("CORS_EXPOSE_HEADERS", "Content-Length"),
	}
}

// loadServerConfig загружает конфигурацию HTTP сервера из переменных окружения.
//
// Читаемые переменные:
//   - API_ADDRESS: адрес для прослушивания (по умолчанию: ":4000")
//   - APP_NAME: название приложения (по умолчанию: "Admin Panel API")
//   - ROOT_PATH: корневой путь URL (по умолчанию: "/admin")
//
// Возвращает:
//   - ServerConfig: структура с настройками сервера
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
