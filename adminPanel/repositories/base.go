// Package repositories содержит реализации слоя доступа к данным (Data Access Layer).
// Этот пакет предоставляет базовые структуры и методы для работы с базой данных,
// включая CRUD операции и фильтрацию данных.
//
// BaseRepository - базовый репозиторий, предоставляющий общие методы для всех репозиториев:
//   - GetByID - получение записи по ID
//   - GetAll - получение всех записей с пагинацией
//   - Count - подсчет записей
//   - Delete - удаление записи
//   - Exists - проверка существования записи
//   - GetFiltered - получение записей с фильтрацией
//
// Все методы возвращают сырые данные в формате map[string]interface{},
// что позволяет гибко работать с различными типами данных.
package repositories

import (
	"context"
	"fmt"     // Нужен для fmt.Sprintf
	"strings" // Нужен для strings.Join

	"adminPanel/database"
)

// BaseRepository - базовый репозиторий для всех сущностей.
// Предоставляет общие методы для работы с базой данных.
type BaseRepository struct {
	db        *database.Database
	tableName string
	schema    string
}

// NewBaseRepository создает новый базовый репозиторий.
//
// Параметры:
//   - db: указатель на соединение с базой данных
//   - tableName: имя таблицы для операций
//   - schema: имя схемы (может быть пустой строкой)
//
// Возвращает:
//   - *BaseRepository: указатель на новый экземпляр BaseRepository
func NewBaseRepository(db *database.Database, tableName, schema string) *BaseRepository {
	return &BaseRepository{
		db:        db,
		tableName: tableName,
		schema:    schema,
	}
}

// FullTableName возвращает полное имя таблицы в формате "схема.таблица".
//
// Возвращает:
//   - string: полное имя таблицы, включая схему (если указана)
func (r *BaseRepository) FullTableName() string {
	if r.schema == "" {
		return r.tableName
	}
	return fmt.Sprintf("%s.%s", r.schema, r.tableName)
}

// GetByID выполняет запрос SELECT для получения записи по уникальному идентификатору.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - id: уникальный идентификатор записи
//
// Возвращает:
//   - map[string]interface{}: найденная запись в виде map
//   - error: ошибка, если произошла
func (r *BaseRepository) GetByID(ctx context.Context, id string) (map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.FullTableName())
	return r.db.FetchOne(ctx, query, id)
}

// GetAll выполняет запрос SELECT для получения всех записей с пагинацией и сортировкой.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - limit: максимальное количество записей
//   - offset: смещение от начала
//   - orderBy: поле для сортировки (по умолчанию "created_at")
//   - orderDir: направление сортировки (по умолчанию "DESC")
//
// Возвращает:
//   - []map[string]interface{}: список найденных записей
//   - error: ошибка, если произошла
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

// Count выполняет запрос COUNT для подсчета количества записей.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - whereClause: условие WHERE для фильтрации (может быть пустым)
//   - params: параметры для подстановки в условие WHERE
//
// Возвращает:
//   - int: количество найденных записей
//   - error: ошибка, если произошла
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

// Delete выполняет запрос DELETE для удаления записи по уникальному идентификатору.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - id: уникальный идентификатор записи для удаления
//
// Возвращает:
//   - bool: true, если запись была удалена, false - если запись не найдена
//   - error: ошибка, если произошла
func (r *BaseRepository) Delete(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.FullTableName())
	affected, err := r.db.Execute(ctx, query, id)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

// Exists выполняет запрос для проверки существования записи по уникальному идентификатору.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - id: уникальный идентификатор записи
//
// Возвращает:
//   - bool: true, если запись существует, false - если не существует
//   - error: ошибка, если произошла
func (r *BaseRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 LIMIT 1", r.FullTableName())
	result, err := r.db.FetchOne(ctx, query, id)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

// GetFiltered выполняет запрос SELECT для получения записей с фильтрацией.
//
// Параметры:
//   - ctx: контекст выполнения запроса
//   - conditions: список условий WHERE для фильтрации
//   - params: параметры для подстановки в условия
//   - orderBy: поле для сортировки (по умолчанию "created_at")
//   - orderDir: направление сортировки (по умолчанию "DESC")
//
// Возвращает:
//   - []map[string]interface{}: список найденных записей
//   - error: ошибка, если произошла
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
