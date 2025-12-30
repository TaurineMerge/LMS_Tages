// Package database предоставляет функциональность для установки и управления
// соединением с базой данных PostgreSQL.
package database

import (
	"context"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection создает, настраивает и проверяет новый пул соединений с базой данных.
// Он использует предоставленную конфигурацию, настраивает трассировку OpenTelemetry
// с помощью otelpgx и выполняет ping для проверки доступности базы данных.
// Возвращает инициализированный *pgxpool.Pool или ошибку в случае сбоя.
func NewConnection(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Парсинг URL для подключения к базе данных из конфигурации.
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Интеграция трассировщика OpenTelemetry для сбора данных о запросах к БД.
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	// Создание нового пула соединений с использованием настроенной конфигурации.
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Проверка соединения с базой данных путем отправки ping-запроса.
	// Если ping не удался, пул соединений закрывается, и возвращается ошибка.
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
