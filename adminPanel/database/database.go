package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"adminPanel/config"
)

// Database - обертка над пулом соединений
type Database struct {
	Pool *pgxpool.Pool
}

var (
	dbInstance *Database
)

// InitDB инициализирует пул соединений (аналог init_db_pool)
func InitDB(settings *config.Settings) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(settings.DatabaseURL)
	if err != nil {
		return nil, err
	}

	poolConfig.MinConns = int32(settings.DatabaseMinPoolSize)
	poolConfig.MaxConns = int32(settings.DatabaseMaxPoolSize)
	
	// Настройки здоровья соединений
	poolConfig.HealthCheckPeriod = 1 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	dbInstance = &Database{Pool: pool}
	log.Println("✅ Database connection pool initialized")
	return dbInstance, nil
}

// Close закрывает пул соединений (аналог close_db_pool)
func Close() {
	if dbInstance != nil && dbInstance.Pool != nil {
		dbInstance.Pool.Close()
		log.Println("Database connection pool closed")
	}
}

// GetDB возвращает инстанс базы данных
func GetDB() *Database {
	return dbInstance
}

// FetchOne - аналог fetch_one из Python
func (db *Database) FetchOne(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowToMap(rows)
}

// FetchAll - аналог fetch_all из Python
func (db *Database) FetchAll(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMap(rows)
}

// Execute - аналог execute из Python
func (db *Database) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// ExecuteReturning - аналог execute_returning из Python
func (db *Database) ExecuteReturning(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowToMap(rows)
}

// Вспомогательная функция для сканирования строки в map
func scanRowToMap(rows pgx.Rows) (map[string]interface{}, error) {
	if !rows.Next() {
		return nil, nil
	}

	fieldDescriptions := rows.FieldDescriptions()
	values := make([]interface{}, len(fieldDescriptions))
	valuePtrs := make([]interface{}, len(fieldDescriptions))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, fd := range fieldDescriptions {
		// Преобразуем значение
		result[string(fd.Name)] = convertValue(values[i])
	}

	return result, nil
}

// Вспомогательная функция для сканирования всех строк
func scanRowsToMap(rows pgx.Rows) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for rows.Next() {
		fieldDescriptions := rows.FieldDescriptions()
		values := make([]interface{}, len(fieldDescriptions))
		valuePtrs := make([]interface{}, len(fieldDescriptions))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		for i, fd := range fieldDescriptions {
			result[string(fd.Name)] = convertValue(values[i])
		}

		results = append(results, result)
	}

	return results, nil
}

// Преобразование значений PostgreSQL в Go типы
func convertValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		// UUID и другие binary типы
		if len(v) == 16 {
			// Это UUID
			return string(v)
		}
		return string(v)
	default:
		return value
	}
}