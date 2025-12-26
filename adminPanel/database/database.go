// Пакет database предоставляет интерфейс для работы с базой данных PostgreSQL.
// Включает инициализацию пула соединений, выполнение запросов и сканирование результатов.
// Поддерживает трассировку с помощью OpenTelemetry.
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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Database представляет соединение с базой данных.
// Содержит пул соединений pgxpool.Pool для выполнения запросов.
type Database struct {
	Pool *pgxpool.Pool
}

// dbInstance глобальная переменная, хранящая единственный экземпляр Database.
// Используется для паттерна singleton.
var (
	dbInstance *Database
)

// InitDB инициализирует пул соединений с базой данных на основе настроек.
// Создает pgxpool.Pool, проверяет подключение и сохраняет экземпляр в dbInstance.
// Возвращает ошибку при неудаче парсинга URL, создания пула или пинга базы данных.
func InitDB(settings *config.Settings) (*Database, error) {
	dbURL := settings.Database.URL()
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	if settings.Database.MinPoolSize < 0 || settings.Database.MinPoolSize > math.MaxInt32 {
		return nil, fmt.Errorf("invalid DatabaseMinPoolSize: %d", settings.Database.MinPoolSize)
	}
	if settings.Database.MaxPoolSize < 0 || settings.Database.MaxPoolSize > math.MaxInt32 {
		return nil, fmt.Errorf("invalid DatabaseMaxPoolSize: %d", settings.Database.MaxPoolSize)
	}
	poolConfig.MinConns = int32(settings.Database.MinPoolSize)
	poolConfig.MaxConns = int32(settings.Database.MaxPoolSize)

	poolConfig.HealthCheckPeriod = 1 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance = &Database{Pool: pool}
	log.Printf("✅ Database connection pool initialized (host=%s, db=%s)",
		settings.Database.Host, settings.Database.Name)
	return dbInstance, nil
}

// Close закрывает пул соединений с базой данных.
// Вызывается для корректного завершения работы с БД.
func Close() {
	if dbInstance != nil && dbInstance.Pool != nil {
		dbInstance.Pool.Close()
		log.Println("Database connection pool closed")
	}
}

// GetDB возвращает глобальный экземпляр Database.
// Используется для получения доступа к пулу соединений из других частей приложения.
func GetDB() *Database {
	return dbInstance
}

// executeQueryReturning выполняет запрос, возвращающий одну строку, с трассировкой.
// Принимает контекст, SQL-запрос, операцию (для трассировки) и аргументы.
// Возвращает результат как map[string]interface{} или nil, если строк нет.
func (db *Database) executeQueryReturning(ctx context.Context, query string, operation string, args ...interface{}) (map[string]interface{}, error) {
	tr := otel.Tracer("admin-panel/database")
	ctx, span := tr.Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", query),
			attribute.String("db.operation", operation),
		),
	)
	defer span.End()

	span.AddEvent("db.query.start", trace.WithAttributes(
		attribute.String("db.query", query),
		attribute.Int("db.args.count", len(args)),
	))

	if len(args) > 0 {
		argsStr := fmt.Sprintf("%v", args)
		span.AddEvent("db.query.params", trace.WithAttributes(
			attribute.String("db.args", argsStr),
		))
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()

	result, err := scanRowToMap(rows)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if result != nil {
		span.SetAttributes(attribute.Int("db.rows_affected", 1))
	} else {
		span.SetAttributes(attribute.Int("db.rows_affected", 0))
	}

	span.AddEvent("db.query.end", trace.WithAttributes(
		attribute.Int("db.rows_returned", len(result)),
	))

	return result, nil
}

// FetchOne выполняет SELECT-запрос, ожидающий одну строку.
// Возвращает первую строку результата как map[string]interface{} или nil, если строк нет.
func (db *Database) FetchOne(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	return db.executeQueryReturning(ctx, query, "SELECT", args...)
}

// FetchAll выполняет SELECT-запрос, возвращающий несколько строк.
// Возвращает все строки результата как []map[string]interface{}.
func (db *Database) FetchAll(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	tr := otel.Tracer("admin-panel/database")
	ctx, span := tr.Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", query),
			attribute.String("db.operation", "SELECT"),
		),
	)
	defer span.End()

	span.AddEvent("db.query.start", trace.WithAttributes(
		attribute.String("db.query", query),
		attribute.Int("db.args.count", len(args)),
	))

	if len(args) > 0 {
		argsStr := fmt.Sprintf("%v", args)
		span.AddEvent("db.query.params", trace.WithAttributes(
			attribute.String("db.args", argsStr),
		))
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()

	results, err := scanRowsToMap(rows)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("db.rows_affected", len(results)))
	span.AddEvent("db.query.end", trace.WithAttributes(
		attribute.Int("db.rows_returned", len(results)),
	))

	return results, nil
}

// Execute выполняет запрос, не возвращающий данные (INSERT, UPDATE, DELETE).
// Возвращает количество затронутых строк.
func (db *Database) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	tr := otel.Tracer("admin-panel/database")
	ctx, span := tr.Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.statement", query),
			attribute.String("db.operation", "EXECUTE"),
		),
	)
	defer span.End()

	span.AddEvent("db.query.start", trace.WithAttributes(
		attribute.String("db.query", query),
		attribute.Int("db.args.count", len(args)),
	))

	if len(args) > 0 {
		argsStr := fmt.Sprintf("%v", args)
		span.AddEvent("db.query.params", trace.WithAttributes(
			attribute.String("db.args", argsStr),
		))
	}

	result, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	rowsAffected := result.RowsAffected()
	span.SetAttributes(attribute.Int("db.rows_affected", int(rowsAffected)))
	span.AddEvent("db.query.end", trace.WithAttributes(
		attribute.Int("db.rows_affected", int(rowsAffected)),
	))

	return rowsAffected, nil
}

// ExecuteReturning выполняет запрос, возвращающий одну строку (например, INSERT ... RETURNING).
// Возвращает результат как map[string]interface{}.
func (db *Database) ExecuteReturning(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	return db.executeQueryReturning(ctx, query, "EXECUTE_RETURNING", args...)
}

// scanRowToMap сканирует одну строку из pgx.Rows в map[string]interface{}.
// Преобразует значения полей в подходящие типы.
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
		result[fd.Name] = convertValue(values[i])
	}

	return result, nil
}

// scanRowsToMap сканирует все строки из pgx.Rows в []map[string]interface{}.
// Преобразует значения полей в подходящие типы.
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

// convertValue преобразует значение из базы данных в подходящий тип.
// Например, []byte в string для UUID или других полей.
func convertValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		if len(v) == 16 {
			return string(v)
		}
		return string(v)
	default:
		return value
	}
}
