// Package database provides utilities for creating and managing database connections.
package database

import (
	"context"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new, instrumented database connection pool.
// It injects an OpenTelemetry tracer into the pgxpool config to automatically
// trace all SQL queries.
func NewConnection(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Parse the config from the database URL
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Set the tracer on the pool config
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	// Create a new pool with the instrumented config
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
