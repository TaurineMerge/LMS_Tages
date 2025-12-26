package repositories

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/database"
	"adminPanel/handlers/dto/request"
)

// CourseRepository предоставляет методы для работы с курсами.
// Встраивает BaseRepository для общих операций.
type CourseRepository struct {
	*BaseRepository
}

// NewCourseRepository создает новый экземпляр CourseRepository.
// Использует таблицу "course_b" в схеме "knowledge_base".
func NewCourseRepository(db *database.Database) *CourseRepository {
	return &CourseRepository{
		BaseRepository: NewBaseRepository(db, "course_b", "knowledge_base"),
	}
}

// Create создает новый курс на основе данных из request.CourseCreate.
// Генерирует UUID и устанавливает время создания и обновления.
// Возвращает созданный курс.
func (r *CourseRepository) Create(ctx context.Context, course request.CourseCreate) (map[string]interface{}, error) {
	query := `
		INSERT INTO knowledge_base.course_b 
		(id, title, description, level, category_id, visibility, image_key, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		course.Title,
		course.Description,
		course.Level,
		course.CategoryID,
		course.Visibility,
		course.ImageKey,
	)
}

// Update обновляет курс по ID на основе данных из request.CourseUpdate.
// Использует COALESCE для обновления только переданных полей.
// Возвращает обновленный курс.
func (r *CourseRepository) Update(ctx context.Context, id string, course request.CourseUpdate) (map[string]interface{}, error) {
	query := `
		UPDATE knowledge_base.course_b 
		SET title = COALESCE($1, title),
			description = COALESCE($2, description),
			level = COALESCE($3, level),
			category_id = COALESCE($4, category_id),
			visibility = COALESCE($5, visibility),
			image_key = COALESCE($6, image_key),
			updated_at = NOW()
		WHERE id = $7
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		course.Title,
		course.Description,
		course.Level,
		course.CategoryID,
		course.Visibility,
		course.ImageKey,
		id,
	)
}

// GetFiltered получает курсы с фильтрами из request.CourseFilter.
// Возвращает список курсов, общее количество и ошибку.
func (r *CourseRepository) GetFiltered(ctx context.Context, filter request.CourseFilter) ([]map[string]interface{}, int, error) {
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

// GetByCategory получает все курсы для заданной категории.
// Сортирует по времени создания в порядке убывания.
func (r *CourseRepository) GetByCategory(ctx context.Context, categoryID string) ([]map[string]interface{}, error) {
	query := `
		SELECT * FROM knowledge_base.course_b
		WHERE category_id = $1
		ORDER BY created_at DESC
	`

	return r.db.FetchAll(ctx, query, categoryID)
}

// ExistsByCategory проверяет существование категории по ID.
// Возвращает true, если категория существует.
func (r *CourseRepository) ExistsByCategory(ctx context.Context, categoryID string) (bool, error) {
	query := "SELECT 1 FROM knowledge_base.category_d WHERE id = $1 LIMIT 1"
	result, err := r.db.FetchOne(ctx, query, categoryID)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}
