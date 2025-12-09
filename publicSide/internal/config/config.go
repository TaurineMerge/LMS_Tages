package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config stores all configuration of the application.
type Config struct {
	DatabaseURL          string
	CORSAllowedOrigins   string
	CORSAllowedMethods   string
	CORSAllowedHeaders   string
	CORSAllowCredentials bool
	Port                 string // New: Application port, defaults to 3000
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

// WithCORSFromEnv configures CORS settings from environment variables.
func WithCORSFromEnv() Option {
	return func(cfg *Config) error {
		var err error

		stringEnvs := map[string]*string{
			"CORS_ALLOWED_ORIGINS": &cfg.CORSAllowedOrigins,
			"CORS_ALLOWED_METHODS": &cfg.CORSAllowedMethods,
			"CORS_ALLOWED_HEADERS": &cfg.CORSAllowedHeaders,
		}

		for key, target := range stringEnvs {
			*target, err = getRequiredEnv(key)
			if err != nil {
				return err
			}
		}

		cfg.CORSAllowCredentials, err = getRequiredEnvAsBool("CORS_ALLOW_CREDENTIALS")
		if err != nil {
			return err
		}

		return nil
	}
}

// WithPortFromEnv configures the application port from environment variables.
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
