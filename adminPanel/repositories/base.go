// Package repositories предоставляет базовые репозитории для работы с данными.
// Включает BaseRepository для общих операций CRUD с базой данных.
package repositories

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/database"
)

// BaseRepository предоставляет базовые методы для работы с таблицей базы данных.
// Содержит ссылку на базу данных, имя таблицы и схему.
type BaseRepository struct {
	db        *database.Database
	tableName string
	schema    string
}

// NewBaseRepository создает новый экземпляр BaseRepository.
// Принимает соединение с БД, имя таблицы и схему.
func NewBaseRepository(db *database.Database, tableName, schema string) *BaseRepository {
	return &BaseRepository{
		db:        db,
		tableName: tableName,
		schema:    schema,
	}
}

// FullTableName возвращает полное имя таблицы с учетом схемы.
// Если схема не указана, возвращает только имя таблицы.
func (r *BaseRepository) FullTableName() string {
	if r.schema == "" {
		return r.tableName
	}
	return fmt.Sprintf("%s.%s", r.schema, r.tableName)
}

// GetByID получает запись по ID.
// Возвращает map[string]interface{} с данными или nil, если не найдено.
func (r *BaseRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.FullTableName())
	return r.db.FetchOne(ctx, query, id)
}

// GetAll получает все записи с пагинацией и сортировкой.
// Принимает limit, offset, orderBy и orderDir. По умолчанию сортирует по created_at DESC.
func (r *BaseRepository) GetAll(ctx context.Context, limit, offset int, orderBy, orderDir string) ([]map[string]interface{}, error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if orderDir == "" {
		orderDir = "DESC"
	}

	query := fmt.Sprintf(`
		SELECT * FROM %s
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, r.FullTableName(), orderBy, orderDir)

	return r.db.FetchAll(ctx, query, limit, offset)
}

// Count подсчитывает количество записей с опциональным условием WHERE.
// Принимает whereClause и параметры для него.
func (r *BaseRepository) Count(ctx context.Context, whereClause string, params ...interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM %s", r.FullTableName())

	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	result, err := r.db.FetchOne(ctx, query, params...)
	if err != nil {
		return 0, err
	}

	if count, ok := result["count"].(int64); ok {
		return int(count), nil
	}
	return 0, nil
}

// Delete удаляет запись по ID.
// Возвращает true, если запись была удалена, false - если не найдена.
func (r *BaseRepository) Delete(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.FullTableName())
	affected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

// Exists проверяет существование записи по ID.
// Возвращает true, если запись существует.
func (r *BaseRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 LIMIT 1", r.FullTableName())
	result, err := r.db.FetchOne(ctx, query, id)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

// GetFiltered получает записи с фильтрами и сортировкой.
// Принимает условия WHERE, параметры, orderBy и orderDir.
func (r *BaseRepository) GetFiltered(ctx context.Context, conditions []string, params []interface{}, orderBy, orderDir string) ([]map[string]interface{}, error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if orderDir == "" {
		orderDir = "DESC"
	}

	query := fmt.Sprintf("SELECT * FROM %s", r.FullTableName())

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)

	return r.db.FetchAll(ctx, query, params...)
}
