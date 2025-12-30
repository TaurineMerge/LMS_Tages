// Package config предоставляет гибкое управление конфигурацией для приложения.
// Он использует паттерн "функциональные опции" для загрузки настроек из переменных окружения.
package config

import (
	"fmt"
	"os"
	"strconv"
)

type (
	// Config является корневой структурой конфигурации, объединяющей все остальные.
	Config struct {
		App            AppConfig
		Server         ServerConfig
		Database       DatabaseConfig
		CORS           CORSConfig
		Otel           OtelConfig
		Log            LogConfig
		OIDC           OIDCConfig
		Minio          MinioConfig
		TestingService TestingServiceConfig
	}

	// AppConfig содержит общие настройки приложения.
	AppConfig struct {
		Dev bool // Dev режим (true/false) - включает горячую перезагрузку шаблонов и заголовки no-cache.
	}

	// ServerConfig содержит настройки HTTP-сервера.
	ServerConfig struct {
		Port string // Порт, на котором будет запущен веб-сервер.
	}

	// DatabaseConfig содержит настройки подключения к базе данных.
	DatabaseConfig struct {
		URL string // Полная строка подключения к PostgreSQL.
	}

	// CORSConfig содержит настройки Cross-Origin Resource Sharing.
	CORSConfig struct {
		AllowedOrigins   string // Разрешенные источники (через запятую).
		AllowedMethods   string // Разрешенные методы (через запятую).
		AllowedHeaders   string // Разрешенные заголовки (через запятую).
		AllowCredentials bool   // Разрешает передачу credentials.
	}

	// OtelConfig содержит настройки OpenTelemetry.
	OtelConfig struct {
		ServiceName       string // Имя сервиса для отображения в Jaeger/Uptrace.
		CollectorEndpoint string // Адрес OTel-коллектора для отправки трассировок.
	}

	// LogConfig содержит настройки логирования.
	LogConfig struct {
		Level string // Уровень логирования (DEBUG, INFO, WARN, ERROR).
	}

	// OIDCConfig содержит настройки OpenID Connect для аутентификации.
	OIDCConfig struct {
		ClientID     string // ID клиента OIDC.
		ClientSecret string // Секрет клиента OIDC.
		IssuerURL    string // URL издателя токенов OIDC.
		RedirectURL  string // URL для перенаправления после аутентификации.
	}

	// MinioConfig содержит настройки подключения к MinIO (S3-совместимое хранилище).
	MinioConfig struct {
		Endpoint  string // Адрес сервера MinIO.
		AccessKey string // Ключ доступа.
		SecretKey string // Секретный ключ.
		Bucket    string // Название бакета.
		UseSSL    bool   // Использовать ли SSL.
		PublicURL string // Публичный URL для доступа к файлам.
	}

	// TestingServiceConfig содержит настройки для внешнего сервиса тестирования.
	TestingServiceConfig struct {
		BaseURL string // Базовый URL сервиса тестирования.
	}
)

// Option определяет тип функции, которая конфигурирует объект *Config.
type Option func(*Config) error

// New создает новый экземпляр Config, применяя переданные опции.
// Возвращает ошибку, если какая-либо из опций не смогла быть применена.
func New(opts ...Option) (*Config, error) {
	cfg := &Config{}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// WithOIDCFromEnv возвращает Option для конфигурации OIDC из переменных окружения.
func WithOIDCFromEnv() Option {
	return func(cfg *Config) error {
		var err error
		cfg.OIDC.ClientID, err = getRequiredEnv("OIDC_CLIENT_ID")
		if err != nil {
			return err
		}
		cfg.OIDC.ClientSecret, err = getRequiredEnv("OIDC_CLIENT_SECRET")
		if err != nil {
			return err
		}
		cfg.OIDC.IssuerURL, err = getRequiredEnv("OIDC_ISSUER_URL")
		if err != nil {
			return err
		}
		cfg.OIDC.RedirectURL, err = getRequiredEnv("OIDC_REDIRECT_URL")
		if err != nil {
			return err
		}
		return nil
	}
}

// WithDBFromEnv возвращает Option для конфигурации базы данных из переменных окружения.
// Поддерживает как полную строку `DATABASE_URL`, так и отдельные `DB_*` переменные.
func WithDBFromEnv() Option {
	return func(cfg *Config) error {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL != "" {
			cfg.Database.URL = databaseURL
			return nil
		}

		requiredEnvs := map[string]string{
			"DB_HOST":     "",
			"DB_PORT":     "",
			"DB_USER":     "",
			"DB_PASSWORD": "",
			"DB_NAME":     "",
		}

		for key := range requiredEnvs {
			value, err := getRequiredEnv(key)
			if err != nil {
				return err
			}
			requiredEnvs[key] = value
		}

		sslMode := getOptionalEnv("DB_SSLMODE", "disable")

		cfg.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			requiredEnvs["DB_USER"],
			requiredEnvs["DB_PASSWORD"],
			requiredEnvs["DB_HOST"],
			requiredEnvs["DB_PORT"],
			requiredEnvs["DB_NAME"],
			sslMode,
		)
		return nil
	}
}

// WithCORSFromEnv возвращает Option для конфигурации CORS из переменных окружения.
// Использует разумные значения по умолчанию, если переменные не установлены.
func WithCORSFromEnv() Option {
	return func(cfg *Config) error {
		cfg.CORS.AllowedOrigins = getOptionalEnv("CORS_ALLOWED_ORIGINS", "*")
		cfg.CORS.AllowedMethods = getOptionalEnv("CORS_ALLOWED_METHODS", "GET")
		cfg.CORS.AllowedHeaders = getOptionalEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")

		var err error
		cfg.CORS.AllowCredentials, err = getOptionalEnvAsBool("CORS_ALLOW_CREDENTIALS", false)
		if err != nil {
			return err
		}

		return nil
	}
}

// WithPortFromEnv возвращает Option для конфигурации порта приложения из переменной `APP_PORT`.
func WithPortFromEnv() Option {
	return func(cfg *Config) error {
		portStr := getOptionalEnv("APP_PORT", "3000")

		_, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("failed to parse APP_PORT environment variable as integer: %w", err)
		}
		cfg.Server.Port = portStr
		return nil
	}
}

// WithTracingFromEnv возвращает Option для конфигурации OpenTelemetry из переменных окружения.
func WithTracingFromEnv() Option {
	return func(cfg *Config) error {
		cfg.Otel.ServiceName = getOptionalEnv("OTEL_SERVICE_NAME", "publicSide")
		var err error
		cfg.Otel.CollectorEndpoint, err = getRequiredEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if err != nil {
			return err
		}
		return nil
	}
}

// WithLogLevelFromEnv возвращает Option для конфигурации уровня логирования из переменной `LOG_LEVEL`.
func WithLogLevelFromEnv() Option {
	return func(cfg *Config) error {
		cfg.Log.Level = getOptionalEnv("LOG_LEVEL", "INFO")
		return nil
	}
}

// WithDevFromEnv возвращает Option для конфигурации режима разработки из переменной `DEV`.
func WithDevFromEnv() Option {
	return func(cfg *Config) error {
		var err error
		cfg.App.Dev, err = getOptionalEnvAsBool("DEV", false)
		if err != nil {
			return err
		}
		return nil
	}
}

// WithMinioFromEnv возвращает Option для конфигурации MinIO из переменных окружения.
func WithMinioFromEnv() Option {
	return func(cfg *Config) error {
		cfg.Minio.Endpoint = getOptionalEnv("MINIO_ENDPOINT", "minio:9000")
		cfg.Minio.AccessKey = getOptionalEnv("MINIO_ACCESS_KEY", "minioadmin")
		cfg.Minio.SecretKey = getOptionalEnv("MINIO_SECRET_KEY", "minioadmin")
		cfg.Minio.Bucket = getOptionalEnv("MINIO_BUCKET", "images")
		var err error
		cfg.Minio.UseSSL, err = getOptionalEnvAsBool("MINIO_USE_SSL", false)
		if err != nil {
			return err
		}
		cfg.Minio.PublicURL = getOptionalEnv("MINIO_PUBLIC_URL", "http://localhost:9000")
		return nil
	}
}

// WithTestingFromEnv возвращает Option для конфигурации внешнего сервиса тестирования.
func WithTestingFromEnv() Option {
	return func(cfg *Config) error {
		var err error
		cfg.TestingService.BaseURL, err = getRequiredEnv("TESTING_SERVICE_BASE_URL")
		if err != nil {
			return err
		}
		return nil
	}
}

// getRequiredEnv извлекает обязательную переменную окружения.
// Возвращает ошибку, если переменная не установлена или пуста.
func getRequiredEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("required environment variable '%s' is not set", key)
	}
	return value, nil
}

// getOptionalEnv извлекает необязательную переменную окружения с значением по умолчанию.
func getOptionalEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getOptionalEnvAsBool извлекает необязательную переменную окружения как boolean.
func getOptionalEnvAsBool(key string, defaultValue bool) (bool, error) {
	valueStr := getOptionalEnv(key, strconv.FormatBool(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse environment variable '%s' as boolean: %w", key, err)
	}
	return value, nil
}
