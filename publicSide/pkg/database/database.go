// Package database provides utilities for creating and managing database connections.
package database

import (
	"context"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new database connection pool using the provided configuration.
// It pings the database to verify the connection before returning the pool.
func NewConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
