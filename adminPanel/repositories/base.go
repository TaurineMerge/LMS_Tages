package repositories

import (
	"context"
	"fmt"     // Нужен для fmt.Sprintf
	"strings" // Нужен для strings.Join

	"adminPanel/database"
)

// BaseRepository - аналог Python BaseRepository
type BaseRepository struct {
	db        *database.Database
	tableName string
	schema    string
}

// NewBaseRepository создает новый базовый репозиторий
func NewBaseRepository(db *database.Database, tableName, schema string) *BaseRepository {
	return &BaseRepository{
		db:        db,
		tableName: tableName,
		schema:    schema,
	}
}

// FullTableName возвращает полное имя таблицы
func (r *BaseRepository) FullTableName() string {
	if r.schema == "" {
		return r.tableName
	}
	return fmt.Sprintf("%s.%s", r.schema, r.tableName)
}

// GetByID - аналог get_by_id из Python
func (r *BaseRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.FullTableName())
	return r.db.FetchOne(ctx, query, id)
}

// GetAll - аналог get_all из Python
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

// Count - аналог count из Python
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

// Delete - аналог delete из Python
func (r *BaseRepository) Delete(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.FullTableName())
	affected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

// Exists - аналог exists из Python
func (r *BaseRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 LIMIT 1", r.FullTableName())
	result, err := r.db.FetchOne(ctx, query, id)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

// GetFiltered - общий метод для фильтрации (аналог get_filtered)
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
