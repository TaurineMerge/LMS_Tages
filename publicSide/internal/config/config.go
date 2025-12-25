// Package config provides flexible configuration management for the application.
// It uses an options-based pattern to load settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
)

type (
	// Config is the root configuration struct that aggregates all other configs.
	Config struct {
		App      AppConfig
		Server   ServerConfig
		Database DatabaseConfig
		CORS     CORSConfig
		Otel     OtelConfig
		Log      LogConfig
		OIDC     OIDCConfig
		Minio    MinioConfig
		TestingService TestingServiceConfig
	}

	// AppConfig holds general application settings.
	AppConfig struct {
		Dev bool
	}

	// ServerConfig holds HTTP server settings.
	ServerConfig struct {
		Port string
	}

	// DatabaseConfig holds database connection settings.
	DatabaseConfig struct {
		URL string
	}

	// CORSConfig holds Cross-Origin Resource Sharing settings.
	CORSConfig struct {
		AllowedOrigins   string
		AllowedMethods   string
		AllowedHeaders   string
		AllowCredentials bool
	}

	// OtelConfig holds OpenTelemetry settings.
	OtelConfig struct {
		ServiceName      string
		CollectorEndpoint string
	}

	// LogConfig holds logging settings.
	LogConfig struct {
		Level string
	}

	// OIDCConfig holds OIDC-specific settings for authentication.
	OIDCConfig struct {
		ClientID     string
		ClientSecret string
		IssuerURL    string
		RedirectURL  string
	}

	// MinioConfig holds MinIO (S3) connection settings.
	MinioConfig struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		UseSSL    bool
		PublicURL string
	}

	// TestingServiceConfig holds settings for the external testing service.
	TestingServiceConfig struct {
		BaseURL string
	}
)

// Option defines a function that configures a Config object.
type Option func(*Config) error

// New returns a new Config struct and an error if configuration is invalid.
func New(opts ...Option) (*Config, error) {
	cfg := &Config{}

	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// WithOIDCFromEnv configures OIDC settings from environment variables.
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


// WithDBFromEnv configures the database settings from environment variables.
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

// WithCORSFromEnv configures CORS settings from environment variables, with sensible defaults.
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

// WithPortFromEnv configures the application port from the APP_PORT environment variable.
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

// WithTracingFromEnv configures OpenTelemetry settings from environment variables.
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

// WithLogLevelFromEnv configures the logging level from the LOG_LEVEL environment variable.
func WithLogLevelFromEnv() Option {
	return func(cfg *Config) error {
		cfg.Log.Level = getOptionalEnv("LOG_LEVEL", "INFO")
		return nil
	}
}

// WithDevFromEnv configures the Dev mode from the DEV environment variable.
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

// WithMinioFromEnv configures MinIO settings from environment variables.
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

// WithTestingFromEnv configures the external testing service settings from environment variables.
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

// getRequiredEnv retrieves a required environment variable.
func getRequiredEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return "", fmt.Errorf("required environment variable '%s' is not set", key)
	}
	return value, nil
}

// getOptionalEnv retrieves an optional environment variable with a default value.
func getOptionalEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// getOptionalEnvAsBool retrieves an optional environment variable as a boolean, with a default.
func getOptionalEnvAsBool(key string, defaultValue bool) (bool, error) {
	valueStr := getOptionalEnv(key, strconv.FormatBool(defaultValue))
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse environment variable '%s' as boolean: %w", key, err)
	}
	return value, nil
}
