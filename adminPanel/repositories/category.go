package repositories

import (
	"context"

	"adminPanel/database"
)

// CategoryRepository - репозиторий для работы с категориями в базе данных
//
// Репозиторий предоставляет методы для выполнения CRUD операций
// с категориями курсов в базе данных.
type CategoryRepository struct {
	*BaseRepository
}

// NewCategoryRepository создает новый репозиторий для категорий
//
// Параметры:
//   - db: экземпляр базы данных
//
// Возвращает:
//   - *CategoryRepository: указатель на новый репозиторий
func NewCategoryRepository(db *database.Database) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db, "category_d", "knowledge_base"),
	}
}

// Create создает новую категорию в базе данных
//
// Параметры:
//   - ctx: контекст выполнения
//   - title: название категории
//
// Возвращает:
//   - map[string]interface{}: созданный объект категории
//   - error: ошибка выполнения (если есть)
func (r *CategoryRepository) Create(ctx context.Context, title string) (map[string]interface{}, error) {
	query := `
		INSERT INTO knowledge_base.category_d 
		(id, title, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, NOW(), NOW())
		RETURNING *
	`
	return r.db.ExecuteReturning(ctx, query, title)
}

// Update обновляет существующую категорию в базе данных
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор категории
//   - title: новое название категории
//
// Возвращает:
//   - map[string]interface{}: обновленный объект категории
//   - error: ошибка выполнения (если есть)
func (r *CategoryRepository) Update(ctx context.Context, id, title string) (map[string]interface{}, error) {
	query := `
		UPDATE knowledge_base.category_d 
		SET title = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING *
	`
	return r.db.ExecuteReturning(ctx, query, title, id)
}

// CountCoursesForCategory подсчитывает количество курсов в указанной категории
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//
// Возвращает:
//   - int: количество курсов
//   - error: ошибка выполнения (если есть)
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

// GetByTitle получает категорию по ее названию
//
// Параметры:
//   - ctx: контекст выполнения
//   - title: название категории
//
// Возвращает:
//   - map[string]interface{}: найденная категория или nil
//   - error: ошибка выполнения (если есть)
func (r *CategoryRepository) GetByTitle(ctx context.Context, title string) (map[string]interface{}, error) {
	query := "SELECT * FROM knowledge_base.category_d WHERE title = $1"
	return r.db.FetchOne(ctx, query, title)
}

// GetAllWithCourses получает все категории с количеством курсов в каждой
//
// Параметры:
//   - ctx: контекст выполнения
//
// Возвращает:
//   - []map[string]interface{}: список категорий с количеством курсов
//   - error: ошибка выполнения (если есть)
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
