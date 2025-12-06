package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ============ ПУЛ ПОДКЛЮЧЕНИЙ ============

var dbPool *pgxpool.Pool

func initDB() error {
	config := getConfig()

	poolConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга конфигурации БД: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5

	ctx := context.Background()
	dbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	log.Println("✅ Успешно подключено к PostgreSQL")
	return nil
}
