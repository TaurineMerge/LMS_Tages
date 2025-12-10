package repositories

import (
	"context"

	"adminPanel/database"
)

// CategoryRepository - репозиторий для категорий
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository создает репозиторий категорий
func NewCategoryRepository(db *database.Database) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db, "category_d", "knowledge_base"),
	}
}

// Create - создание категории
func (r *CategoryRepository) Create(ctx context.Context, title string) (map[string]interface{}, error) {
	query := `
		INSERT INTO knowledge_base.category_d 
		(id, title, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, NOW(), NOW())
		RETURNING id, title, created_at, updated_at
	`
	return r.db.ExecuteReturning(ctx, query, title)
}

// Update - обновление категории
func (r *CategoryRepository) Update(ctx context.Context, id, title string) (map[string]interface{}, error) {
	query := `
		UPDATE knowledge_base.category_d 
		SET title = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, created_at, updated_at
	`
	return r.db.ExecuteReturning(ctx, query, title, id)
}

// CountCoursesForCategory - подсчет курсов в категории
func (r *CategoryRepository) CountCoursesForCategory(ctx context.Context, categoryID string) (int, error) {
	query := `
		SELECT COUNT(*) as count 
		FROM knowledge_base.course_b 
		WHERE category_id = $1
	`
	result, err := r.db.FetchOne(ctx, query, categoryID)
	if err != nil {
		return 0, err
	}

	if count, ok := result["count"].(int64); ok {
		return int(count), nil
	}
	return 0, nil
}

// GetByTitle - получение категории по названию
func (r *CategoryRepository) GetByTitle(ctx context.Context, title string) (map[string]interface{}, error) {
	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d 
		WHERE title = $1
	`
	return r.db.FetchOne(ctx, query, title)
}

// GetAllWithPagination - получение категорий с пагинацией
func (r *CategoryRepository) GetAllWithPagination(ctx context.Context, limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d 
		ORDER BY title ASC
		LIMIT $1 OFFSET $2
	`
	return r.db.FetchAll(ctx, query, limit, offset)
}

// CountAll - общее количество категорий
func (r *CategoryRepository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) as count FROM knowledge_base.category_d`
	result, err := r.db.FetchOne(ctx, query)
	if err != nil {
		return 0, err
	}

	if count, ok := result["count"].(int64); ok {
		return int(count), nil
	}
	return 0, nil
}
