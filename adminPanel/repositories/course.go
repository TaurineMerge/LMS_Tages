package repositories

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/database"
	"adminPanel/handlers/dto/request"
)

// CourseRepository - репозиторий для работы с курсами в базе данных
//
// Репозиторий предоставляет методы для выполнения CRUD операций
// с курсами в базе данных, включая фильтрацию и пагинацию.
type CourseRepository struct {
	*BaseRepository
}

// NewCourseRepository создает новый репозиторий для курсов
//
// Параметры:
//   - db: экземпляр базы данных
//
// Возвращает:
//   - *CourseRepository: указатель на новый репозиторий
func NewCourseRepository(db *database.Database) *CourseRepository {
	return &CourseRepository{
		BaseRepository: NewBaseRepository(db, "course_b", "knowledge_base"),
	}
}

// Create создает новый курс в базе данных
//
// Параметры:
//   - ctx: контекст выполнения
//   - course: данные для создания курса
//
// Возвращает:
//   - map[string]interface{}: созданный объект курса
//   - error: ошибка выполнения (если есть)
func (r *CourseRepository) Create(ctx context.Context, course request.CourseCreate) (map[string]interface{}, error) {
	query := `
		INSERT INTO knowledge_base.course_b 
		(id, title, description, level, category_id, visibility, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		course.Title,
		course.Description,
		course.Level,
		course.CategoryID,
		course.Visibility,
	)
}

// Update обновляет существующий курс в базе данных
//
// Использует COALESCE для частичного обновления - если значение
// равно nil, оно не изменяется.
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор курса
//   - course: данные для обновления курса
//
// Возвращает:
//   - map[string]interface{}: обновленный объект курса
//   - error: ошибка выполнения (если есть)
func (r *CourseRepository) Update(ctx context.Context, id string, course request.CourseUpdate) (map[string]interface{}, error) {
	query := `
		UPDATE knowledge_base.course_b 
		SET title = COALESCE($1, title),
			description = COALESCE($2, description),
			level = COALESCE($3, level),
			category_id = COALESCE($4, category_id),
			visibility = COALESCE($5, visibility),
			updated_at = NOW()
		WHERE id = $6
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		course.Title,
		course.Description,
		course.Level,
		course.CategoryID,
		course.Visibility,
		id,
	)
}

// GetFiltered получает курсы с фильтрацией и пагинацией
//
// Поддерживает фильтрацию по уровню сложности, видимости
// и категории. Возвращает отсортированные по дате создания курсы.
//
// Параметры:
//   - ctx: контекст выполнения
//   - filter: фильтр для поиска курсов
//
// Возвращает:
//   - []map[string]interface{}: список курсов
//   - int: общее количество курсов
//   - error: ошибка выполнения (если есть)
func (r *CourseRepository) GetFiltered(ctx context.Context, filter request.CourseFilter) ([]map[string]interface{}, int, error) {
	// Строим WHERE условия
	var conditions []string
	var params []interface{}
	paramCounter := 1

	if filter.Level != "" {
		conditions = append(conditions, fmt.Sprintf("level = $%d", paramCounter))
		params = append(params, filter.Level)
		paramCounter++
	}

	if filter.Visibility != "" {
		conditions = append(conditions, fmt.Sprintf("visibility = $%d", paramCounter))
		params = append(params, filter.Visibility)
		paramCounter++
	}

	if filter.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", paramCounter))
		params = append(params, filter.CategoryID)
		paramCounter++
	}

	countQuery := "SELECT COUNT(*) as count FROM knowledge_base.course_b"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	countResult, err := r.db.FetchOne(ctx, countQuery, params...)
	if err != nil {
		return nil, 0, err
	}

	total := 0
	if count, ok := countResult["count"].(int64); ok {
		total = int(count)
	}

	query := "SELECT * FROM knowledge_base.course_b"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)

	params = append(params, filter.Limit, (filter.Page-1)*filter.Limit)

	data, err := r.db.FetchAll(ctx, query, params...)
	if err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// GetByCategory получает все курсы указанной категории
//
// Возвращает курсы, отсортированные по дате создания.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//
// Возвращает:
//   - []map[string]interface{}: список курсов категории
//   - error: ошибка выполнения (если есть)
func (r *CourseRepository) GetByCategory(ctx context.Context, categoryID string) ([]map[string]interface{}, error) {
	query := `
		SELECT * FROM knowledge_base.course_b
		WHERE category_id = $1
		ORDER BY created_at DESC
	`

	return r.db.FetchAll(ctx, query, categoryID)
}

// ExistsByCategory проверяет существование категории
//
// Используется для валидации перед операциями с курсами.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//
// Возвращает:
//   - bool: true, если категория существует
//   - error: ошибка выполнения (если есть)
func (r *CourseRepository) ExistsByCategory(ctx context.Context, categoryID string) (bool, error) {
	query := "SELECT 1 FROM knowledge_base.category_d WHERE id = $1 LIMIT 1"
	result, err := r.db.FetchOne(ctx, query, categoryID)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}
