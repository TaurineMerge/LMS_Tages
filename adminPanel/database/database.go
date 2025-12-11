// Пакет database предоставляет функции для работы с PostgreSQL
//
// Пакет реализует:
//   - Инициализацию пула соединений
//   - Выполнение SQL-запросов
//   - Преобразование результатов в map
//   - Управление жизненным циклом соединений
//
// Основные функции:
//   - InitDB: инициализация пула соединений
//   - Close: закрытие пула
//   - FetchOne: выполнение запроса с одним результатом
//   - FetchAll: выполнение запроса с несколькими результатами
//   - Execute: выполнение запроса без результата
//   - ExecuteReturning: выполнение запроса с возвратом результата
package database

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"adminPanel/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database - обертка над пулом соединений PostgreSQL
//
// Структура предоставляет удобный интерфейс для выполнения
// SQL-запросов и управления соединениями.
type Database struct {
	Pool *pgxpool.Pool
}

var (
	dbInstance *Database
)

// InitDB инициализирует пул соединений с PostgreSQL
//
// Функция создает и настраивает пул соединений на основе
// конфигурационных параметров.
//
// Параметры:
//   - settings: конфигурация приложения
//
// Возвращает:
//   - *Database: указатель на экземпляр базы данных
//   - error: ошибка инициализации (если есть)
func InitDB(settings *config.Settings) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(settings.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if settings.DatabaseMinPoolSize < 0 || settings.DatabaseMinPoolSize > math.MaxInt32 {
		return nil, fmt.Errorf("invalid DatabaseMinPoolSize: %d", settings.DatabaseMinPoolSize)
	}
	if settings.DatabaseMaxPoolSize < 0 || settings.DatabaseMaxPoolSize > math.MaxInt32 {
		return nil, fmt.Errorf("invalid DatabaseMaxPoolSize: %d", settings.DatabaseMaxPoolSize)
	}
	poolConfig.MinConns = int32(settings.DatabaseMinPoolSize) //nolint:gosec // validated above
	poolConfig.MaxConns = int32(settings.DatabaseMaxPoolSize) //nolint:gosec // validated above

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

// Close закрывает пул соединений с базой данных
//
// Функция освобождает все ресурсы, связанные с пулом соединений.
// Должна вызываться при завершении работы приложения.
func Close() {
	if dbInstance != nil && dbInstance.Pool != nil {
		dbInstance.Pool.Close()
		log.Println("Database connection pool closed")
	}
}

// GetDB возвращает текущий экземпляр базы данных
//
// Возвращает:
//   - *Database: указатель на экземпляр базы данных
func GetDB() *Database {
	return dbInstance
}

// FetchOne выполняет SQL-запрос и возвращает одну строку результата
//
// Функция аналогична fetch_one из Python. Возвращает первую
// строку результата или nil, если строк нет.
//
// Параметры:
//   - ctx: контекст выполнения
//   - query: SQL-запрос
//   - args: аргументы для запроса
//
// Возвращает:
//   - map[string]interface{}: строка результата или nil
//   - error: ошибка выполнения (если есть)
func (db *Database) FetchOne(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowToMap(rows)
}

// FetchAll выполняет SQL-запрос и возвращает все строки результата
//
// Функция аналогична fetch_all из Python. Возвращает все
// строки результата в виде слайса map.
//
// Параметры:
//   - ctx: контекст выполнения
//   - query: SQL-запрос
//   - args: аргументы для запроса
//
// Возвращает:
//   - []map[string]interface{}: слайс строк результата
//   - error: ошибка выполнения (если есть)
func (db *Database) FetchAll(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowsToMap(rows)
}

// Execute выполняет SQL-запрос без возврата результата
//
// Функция аналогична execute из Python. Используется для
// INSERT, UPDATE, DELETE и других запросов, не возвращающих данные.
//
// Параметры:
//   - ctx: контекст выполнения
//   - query: SQL-запрос
//   - args: аргументы для запроса
//
// Возвращает:
//   - int64: количество затронутых строк
//   - error: ошибка выполнения (если есть)
func (db *Database) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// ExecuteReturning выполняет SQL-запрос с возвратом результата
//
// Функция аналогична execute_returning из Python. Используется
// для запросов, которые возвращают данные (например, INSERT ... RETURNING).
//
// Параметры:
//   - ctx: контекст выполнения
//   - query: SQL-запрос
//   - args: аргументы для запроса
//
// Возвращает:
//   - map[string]interface{}: строка результата
//   - error: ошибка выполнения (если есть)
func (db *Database) ExecuteReturning(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRowToMap(rows)
}

// scanRowToMap сканирует одну строку результата в map
//
// Функция преобразует строку результата запроса в map,
// где ключи - это имена столбцов, а значения - данные.
//
// Параметры:
//   - rows: результат выполнения запроса
//
// Возвращает:
//   - map[string]interface{}: строка результата
//   - error: ошибка сканирования (если есть)
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
		result[fd.Name] = convertValue(values[i])
	}

	return result, nil
}

// scanRowsToMap сканирует все строки результата в слайс map
//
// Функция преобразует все строки результата запроса в слайс map.
//
// Параметры:
//   - rows: результат выполнения запроса
//
// Возвращает:
//   - []map[string]interface{}: слайс строк результата
//   - error: ошибка сканирования (если есть)
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
			result[fd.Name] = convertValue(values[i])
		}

		results = append(results, result)
	}

	return results, nil
}

// convertValue преобразует значения PostgreSQL в Go типы
//
// Функция обрабатывает специфические типы данных PostgreSQL,
// такие как UUID, и преобразует их в соответствующие Go типы.
//
// Параметры:
//   - value: значение из PostgreSQL
//
// Возвращает:
//   - interface{}: преобразованное значение
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
