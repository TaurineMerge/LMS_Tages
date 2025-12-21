// Package config provides flexible configuration management for the application.
// It uses an options-based pattern to load settings from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config stores all configuration of the application.
type Config struct {
	// DatabaseURL is the connection string for the PostgreSQL database.
	DatabaseURL string
	// CORSAllowedOrigins is a comma-separated list of origins that are allowed to make cross-origin requests.
	CORSAllowedOrigins string
	// CORSAllowedMethods is a comma-separated list of methods that are allowed for cross-origin requests.
	CORSAllowedMethods string
	// CORSAllowedHeaders is a comma-separated list of headers that are allowed for cross-origin requests.
	CORSAllowedHeaders string
	// CORSAllowCredentials indicates whether credentials can be shared in cross-origin requests.
	CORSAllowCredentials bool
	// Port is the network port on which the application server will listen.
	Port string
	// OTELServiceName is the name of the service for OpenTelemetry.
	OTELServiceName string
	// OTELCollectorEndpoint is the endpoint of the OpenTelemetry collector.
	OTELCollectorEndpoint string
	// LogLevel is the level for application logging (e.g., DEBUG, INFO, WARN, ERROR).
	LogLevel string
	// Dev indicates whether the application is running in development mode.
	Dev bool
	// OIDC-specific
	OIDCClientID         string
	OIDCClientSecret     string
	OIDCIssuerURL        string
	OIDCRedirectURL      string
}

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
		cfg.OIDCClientID, err = getRequiredEnv("OIDC_CLIENT_ID")
		if err != nil {
			return err
		}
		cfg.OIDCClientSecret, err = getRequiredEnv("OIDC_CLIENT_SECRET")
		if err != nil {
			return err
		}
		cfg.OIDCIssuerURL, err = getRequiredEnv("OIDC_ISSUER_URL")
		if err != nil {
			return err
		}
		cfg.OIDCRedirectURL, err = getRequiredEnv("OIDC_REDIRECT_URL")
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
			cfg.DatabaseURL = databaseURL
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

		cfg.DatabaseURL = fmt.Sprintf(
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
		cfg.CORSAllowedOrigins = getOptionalEnv("CORS_ALLOWED_ORIGINS", "*")
		cfg.CORSAllowedMethods = getOptionalEnv("CORS_ALLOWED_METHODS", "GET")
		cfg.CORSAllowedHeaders = getOptionalEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization")

		var err error
		cfg.CORSAllowCredentials, err = getOptionalEnvAsBool("CORS_ALLOW_CREDENTIALS", false)
		if err != nil {
			return err
		}

		return nil
	}
}

// WithPortFromEnv configures the application port from the APP_PORT environment variable.
// It defaults to "3000" if the variable is not set.
func WithPortFromEnv() Option {
	return func(cfg *Config) error {
		portStr := getOptionalEnv("APP_PORT", "3000") // Get port as string, default empty

		_, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("failed to parse APP_PORT environment variable as integer: %w", err)
		}
		cfg.Port = portStr
		return nil
	}
}

// WithTracingFromEnv configures OpenTelemetry settings from environment variables.
func WithTracingFromEnv() Option {
	return func(cfg *Config) error {
		cfg.OTELServiceName = getOptionalEnv("OTEL_SERVICE_NAME", "publicSide")
		var err error
		cfg.OTELCollectorEndpoint, err = getRequiredEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
		if err != nil {
			return err
		}
		return nil
	}
}

// WithLogLevelFromEnv configures the logging level from the LOG_LEVEL environment variable.
func WithLogLevelFromEnv() Option {
	return func(cfg *Config) error {
		cfg.LogLevel = getOptionalEnv("LOG_LEVEL", "INFO")
		return nil
	}
}

// WithDevFromEnv configures the Dev mode from the DEV environment variable.
func WithDevFromEnv() Option {
	return func(cfg *Config) error {
		var err error
		cfg.Dev, err = getOptionalEnvAsBool("DEV", false)
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

// getRequiredEnvAsBool retrieves a required environment variable as a boolean.
func getRequiredEnvAsBool(key string) (bool, error) {
	valueStr, err := getRequiredEnv(key)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse environment variable '%s' as boolean: %w", key, err)
	}
	return value, nil
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
