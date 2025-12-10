package repositories

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/database"
	"adminPanel/models"
)

// CourseRepository - репозиторий для курсов
type CourseRepository struct {
	*BaseRepository
}

// NewCourseRepository создает репозиторий курсов
func NewCourseRepository(db *database.Database) *CourseRepository {
	return &CourseRepository{
		BaseRepository: NewBaseRepository(db, "course_b", "knowledge_base"),
	}
}

// Create - создание курса
func (r *CourseRepository) Create(ctx context.Context, course models.CourseCreate) (map[string]interface{}, error) {
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

// Update - обновление курса
func (r *CourseRepository) Update(ctx context.Context, id string, course models.CourseUpdate) (map[string]interface{}, error) {
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

// GetFiltered - получение курсов с фильтрацией и пагинацией
func (r *CourseRepository) GetFiltered(ctx context.Context, filter models.CourseFilter) ([]map[string]interface{}, int, error) {
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
	
	// Запрос для подсчета общего количества
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
	
	// Запрос для получения данных с пагинацией
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

// GetByCategory - получение курсов по категории
func (r *CourseRepository) GetByCategory(ctx context.Context, categoryID string) ([]map[string]interface{}, error) {
	query := `
		SELECT * FROM knowledge_base.course_b
		WHERE category_id = $1
		ORDER BY created_at DESC
	`
	
	return r.db.FetchAll(ctx, query, categoryID)
}

// ExistsByCategory - проверка существования категории
func (r *CourseRepository) ExistsByCategory(ctx context.Context, categoryID string) (bool, error) {
	query := "SELECT 1 FROM knowledge_base.category_d WHERE id = $1 LIMIT 1"
	result, err := r.db.FetchOne(ctx, query, categoryID)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}