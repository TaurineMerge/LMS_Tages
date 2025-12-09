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
		RETURNING *
	`
	return r.db.ExecuteReturning(ctx, query, title)
}

// Update - обновление категории
func (r *CategoryRepository) Update(ctx context.Context, id, title string) (map[string]interface{}, error) {
	query := `
		UPDATE knowledge_base.category_d 
		SET title = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING *
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
	query := "SELECT * FROM knowledge_base.category_d WHERE title = $1"
	return r.db.FetchOne(ctx, query, title)
}

// GetAllWithCourses - получение всех категорий с количеством курсов
func (r *CategoryRepository) GetAllWithCourses(ctx context.Context) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			c.*,
			COUNT(cb.id) as course_count
		FROM knowledge_base.category_d c
		LEFT JOIN knowledge_base.course_b cb ON c.id = cb.category_id
		GROUP BY c.id
		ORDER BY c.title
	`
	return r.db.FetchAll(ctx, query)
}