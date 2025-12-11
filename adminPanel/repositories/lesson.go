package repositories

import (
	"context"
	"encoding/json"

	"adminPanel/database"
	"adminPanel/models"
)

// LessonRepository - репозиторий для работы с уроками в базе данных
//
// Репозиторий предоставляет методы для выполнения CRUD операций
// с уроками, включая работу с JSON-контентом.
type LessonRepository struct {
	*BaseRepository
}

// NewLessonRepository создает новый репозиторий для уроков
//
// Параметры:
//   - db: экземпляр базы данных
//
// Возвращает:
//   - *LessonRepository: указатель на новый репозиторий
func NewLessonRepository(db *database.Database) *LessonRepository {
	return &LessonRepository{
		BaseRepository: NewBaseRepository(db, "lesson_d", "knowledge_base"),
	}
}

// Create создает новый урок в базе данных
//
// Преобразует контент урока в JSON формат для хранения.
//
// Параметры:
//   - ctx: контекст выполнения
//   - courseID: уникальный идентификатор курса
//   - categoryID: уникальный идентификатор категории
//   - lesson: данные для создания урока
//
// Возвращает:
//   - map[string]interface{}: созданный объект урока
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) Create(ctx context.Context, courseID, categoryID string, lesson models.LessonCreate) (map[string]interface{}, error) {
	contentJSON, _ := json.Marshal(lesson.Content)

	query := `
		INSERT INTO knowledge_base.lesson_d 
		(id, title, category_id, course_id, content, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4::jsonb, NOW(), NOW())
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		lesson.Title,
		categoryID,
		courseID,
		string(contentJSON),
	)
}

// Update обновляет существующий урок в базе данных
//
// Преобразует контент урока в JSON формат для хранения.
// Использует COALESCE для частичного обновления.
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор урока
//   - courseID: уникальный идентификатор курса
//   - lesson: данные для обновления урока
//
// Возвращает:
//   - map[string]interface{}: обновленный объект урока
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) Update(ctx context.Context, id, courseID string, lesson models.LessonUpdate) (map[string]interface{}, error) {
	contentJSON, _ := json.Marshal(lesson.Content)

	query := `
		UPDATE knowledge_base.lesson_d 
		SET title = COALESCE($1, title),
			category_id = COALESCE($2, category_id),
			content = COALESCE($4::jsonb, content),
			updated_at = NOW()
		WHERE id = $3 AND course_id = $5
		RETURNING *
	`

	return r.db.ExecuteReturning(ctx, query,
		lesson.Title,
		lesson.CategoryID,
		id,
		string(contentJSON),
		courseID,
	)
}

// GetByCourse получает уроки указанного курса
//
// Возвращает уроки, отсортированные по дате создания.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//   - courseID: уникальный идентификатор курса
//   - limit: количество элементов
//   - offset: смещение
//
// Возвращает:
//   - []map[string]interface{}: список уроков
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) GetByCourse(ctx context.Context, categoryID, courseID string, limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT * FROM knowledge_base.lesson_d
		WHERE category_id = $1 AND course_id = $2
		ORDER BY created_at
		LIMIT $3 OFFSET $4
	`

	return r.db.FetchAll(ctx, query, categoryID, courseID, limit, offset)
}

// GetByIDAndCourse получает урок по уникальному идентификатору и курсу
//
// Проверяет принадлежность урока указанной категории и курсу.
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор урока
//   - categoryID: уникальный идентификатор категории
//   - courseID: уникальный идентификатор курса
//
// Возвращает:
//   - map[string]interface{}: найденный урок или nil
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) GetByIDAndCourse(ctx context.Context, id, categoryID, courseID string) (map[string]interface{}, error) {
	query := `
		SELECT * FROM knowledge_base.lesson_d
		WHERE id = $1 AND category_id = $2 AND course_id = $3
	`

	return r.db.FetchOne(ctx, query, id, categoryID, courseID)
}

// CountByCourse подсчитывает количество уроков в указанном курсе
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//   - courseID: уникальный идентификатор курса
//
// Возвращает:
//   - int: количество уроков
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) CountByCourse(ctx context.Context, categoryID, courseID string) (int, error) {
	query := `
		SELECT COUNT(*) as count FROM knowledge_base.lesson_d
		WHERE category_id = $1 AND course_id = $2
	`
	res, err := r.db.FetchOne(ctx, query, categoryID, courseID)
	if err != nil {
		return 0, err
	}
	if res == nil {
		return 0, nil
	}
	if count, ok := res["count"].(int64); ok {
		return int(count), nil
	}
	return 0, nil
}

// ParseContent парсит JSON контент урока
//
// Преобразует JSON-данные из базы данных в map[string]interface{}.
// Если парсинг не удается, возвращает пустой map.
//
// Параметры:
//   - data: данные урока с JSON-контентом
//
// Возвращает:
//   - map[string]interface{}: данные с распарсенным контентом
//   - error: ошибка выполнения (если есть)
func (r *LessonRepository) ParseContent(data map[string]interface{}) (map[string]interface{}, error) {
	if content, ok := data["content"].([]byte); ok {
		var parsedContent map[string]interface{}
		if err := json.Unmarshal(content, &parsedContent); err == nil {
			data["content"] = parsedContent
		} else {
			data["content"] = make(map[string]interface{})
		}
	} else if contentStr, ok := data["content"].(string); ok {
		var parsedContent map[string]interface{}
		if err := json.Unmarshal([]byte(contentStr), &parsedContent); err == nil {
			data["content"] = parsedContent
		} else {
			data["content"] = make(map[string]interface{})
		}
	} else {
		data["content"] = make(map[string]interface{})
	}

	return data, nil
}
