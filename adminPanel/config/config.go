// Пакет config предоставляет структуры и функции для загрузки и валидации конфигурации приложения adminPanel.
// Он считывает настройки из переменных окружения и предоставляет удобный доступ к ним.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DatabaseConfig содержит настройки подключения к базе данных PostgreSQL.
// Включает параметры хоста, порта, пользователя, пароля, имени базы данных,
// режима SSL и размеров пула соединений.
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

// URL возвращает строку подключения к базе данных в формате PostgreSQL DSN.
// Формирует URL на основе полей структуры, включая параметры sslmode.
func (d *DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// OTelConfig содержит настройки для OpenTelemetry.
// Включает endpoint для экспорта, имя сервиса, протокол и флаг включения.
type OTelConfig struct {
	Endpoint    string
	ServiceName string
	Protocol    string
	Enabled     bool
}

// KeycloakConfig содержит настройки для интеграции с Keycloak.
// Включает URL issuer, audience, JWKS URL, client ID, secret и имя приложения.
type KeycloakConfig struct {
	IssuerURL    string
	Audience     string
	JWKSURL      string
	ClientID     string
	ClientSecret string
	AppName      string
}

// CORSConfig содержит настройки для Cross-Origin Resource Sharing (CORS).
// Определяет разрешенные origins, методы, заголовки, credentials и exposed headers.
type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
}

// ServerConfig содержит настройки сервера.
// Включает адрес прослушивания, имя приложения и корневой путь API.
type ServerConfig struct {
	Address  string
	AppName  string
	RootPath string
}

// MinioConfig содержит настройки для подключения к MinIO (S3-compatible storage).
// Включает endpoint, ключи доступа, имя bucket, флаг SSL и публичный URL.
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	PublicURL string
}

// TestModuleConfig содержит настройки для тестового модуля.
// Включает базовый URL и флаг включения модуля.
type TestModuleConfig struct {
	BaseURL string
	Enabled bool
}

// Settings объединяет все конфигурационные структуры в одну.
// Содержит настройки базы данных, OTel, Keycloak, CORS, сервера, MinIO, тестового модуля и флаг отладки.
type Settings struct {
	Database   DatabaseConfig
	OTel       OTelConfig
	Keycloak   KeycloakConfig
	CORS       CORSConfig
	Server     ServerConfig
	Debug      bool
	Minio      MinioConfig
	TestModule TestModuleConfig
}

// Validate проверяет наличие обязательных переменных окружения для базы данных.
// Возвращает ошибку, если отсутствуют DB_HOST, DB_USER, DB_PASSWORD или DB_NAME (или DATABASE_URL).
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

// NewSettings создает новый экземпляр Settings, загружая конфигурацию из переменных окружения.
// Использует вспомогательные функции для загрузки каждой части конфигурации.
func NewSettings() *Settings {
	return &Settings{
		Database:   loadDatabaseConfig(),
		OTel:       loadOTelConfig(),
		Keycloak:   loadKeycloakConfig(),
		CORS:       loadCORSConfig(),
		Server:     loadServerConfig(),
		Debug:      getEnvAsBool("DEBUG", false),
		Minio:      loadMinioConfig(),
		TestModule: loadTestModuleConfig(),
	}
}

// loadDatabaseConfig загружает конфигурацию базы данных из переменных окружения.
// Если задана DATABASE_URL, парсит её; иначе использует отдельные переменные DB_HOST, DB_PORT и т.д.
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

// parseDatabaseURL разбирает строку DATABASE_URL в формате PostgreSQL DSN.
// Возвращает host, port, user, password, name базы данных и sslmode.
// Если порт не указан, использует 5432; sslmode по умолчанию "disable".
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

// loadOTelConfig загружает конфигурацию OpenTelemetry из переменных окружения.
// Включает endpoint, service name, protocol и определяет, включен ли OTel (по наличию endpoint).
func loadOTelConfig() OTelConfig {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	return OTelConfig{
		Endpoint:    endpoint,
		ServiceName: getEnv("OTEL_SERVICE_NAME", "admin-panel"),
		Protocol:    getEnv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc"),
		Enabled:     endpoint != "",
	}
}

// loadKeycloakConfig загружает конфигурацию Keycloak из переменных окружения.
// Если JWKS_URL не задан, формирует его на основе issuer URL.
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

// loadCORSConfig загружает настройки CORS из переменных окружения.
// Определяет разрешенные origins, методы, заголовки и т.д.
func loadCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     getEnv("CORS_ALLOW_ORIGINS", "*"),
		AllowMethods:     getEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		AllowHeaders:     getEnv("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false),
		ExposeHeaders:    getEnv("CORS_EXPOSE_HEADERS", "Content-Length"),
	}
}

// loadServerConfig загружает настройки сервера из переменных окружения.
// Включает адрес, имя приложения и корневой путь.
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Address:  getEnv("API_ADDRESS", ":4000"),
		AppName:  getEnv("APP_NAME", "Admin Panel API"),
		RootPath: getEnv("ROOT_PATH", "/admin"),
	}
}

// loadMinioConfig загружает настройки MinIO из переменных окружения.
// Включает endpoint, ключи, bucket, SSL и публичный URL.
func loadMinioConfig() MinioConfig {
	return MinioConfig{
		Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		Bucket:    getEnv("MINIO_BUCKET", "snapshots"),
		UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		PublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
	}
}

// loadTestModuleConfig загружает настройки тестового модуля из переменных окружения.
// Включает базовый URL и флаг включения.
func loadTestModuleConfig() TestModuleConfig {
	return TestModuleConfig{
		BaseURL: getEnv("TEST_MODULE_BASE_URL", "http://localhost:8080"),
		Enabled: getEnvAsBool("TEST_MODULE_ENABLED", false),
	}
}

// GetCORSOrigins возвращает список разрешенных origins для CORS.
// Если AllowOrigins равно "*", возвращает ["*"]; иначе разбивает строку по запятым и удаляет пробелы.
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

// DatabaseURL возвращает строку подключения к базе данных, используя метод URL() структуры DatabaseConfig.
// Это удобный метод для получения полного DSN.
func (s *Settings) DatabaseURL() string {
	return s.Database.URL()
}

// getEnv получает значение переменной окружения по ключу, возвращая defaultValue, если переменная не установлена.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int, возвращая defaultValue при ошибке или отсутствии.
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool получает значение переменной окружения как bool, возвращая defaultValue при ошибке или отсутствии.
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
