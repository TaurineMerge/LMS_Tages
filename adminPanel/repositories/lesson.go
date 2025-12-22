package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"adminPanel/database"
	"adminPanel/handlers/dto/request"
	"adminPanel/models"

	"github.com/jackc/pgx/v5"
)

// LessonRepository - репозиторий для работы с уроками в базе данных
type LessonRepository struct {
	db *database.Database
	// Убираем BaseRepository, т.к. методы теперь специфичны
}

// NewLessonRepository создает новый репозиторий для уроков
func NewLessonRepository(db *database.Database) *LessonRepository {
	return &LessonRepository{
		db: db,
	}
}

// GetAllByCourseID получает уроки по идентификатору курса с пагинацией и сортировкой
func (r *LessonRepository) GetAllByCourseID(ctx context.Context, courseID string, limit, offset int, sortBy, sortOrder string) ([]models.Lesson, error) {
	allowedSortBy := map[string]bool{"title": true, "created_at": true, "updated_at": true}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}
	if !(strings.EqualFold(sortOrder, "ASC") || strings.EqualFold(sortOrder, "DESC")) {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT id, title, course_id, html_content, created_at, updated_at
		FROM knowledge_base.lesson_d
		WHERE course_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, sortOrder)

	rows, err := r.db.Pool.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var lesson models.Lesson
		var htmlContent *string
		if err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.CourseID, &htmlContent, &lesson.CreatedAt, &lesson.UpdatedAt); err != nil {
			return nil, err
		}
		if htmlContent != nil {
			lesson.HTMLContent = *htmlContent
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// CountByCourseID подсчитывает количество уроков по идентификатору курса
func (r *LessonRepository) CountByCourseID(ctx context.Context, courseID string) (int, error) {
	query := `SELECT COUNT(*) FROM knowledge_base.lesson_d WHERE course_id = $1`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, courseID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByID получает урок по его уникальному идентификатору
func (r *LessonRepository) GetByID(ctx context.Context, lessonID string) (*models.LessonDetailed, error) {
	query := `SELECT id, title, course_id, html_content, created_at, updated_at, content FROM knowledge_base.lesson_d WHERE id = $1`

	row := r.db.Pool.QueryRow(ctx, query, lessonID)

	var lesson models.LessonDetailed
	var contentBytes []byte
	var htmlContent *string

	err := row.Scan(&lesson.ID, &lesson.Title, &lesson.CourseID, &htmlContent, &lesson.CreatedAt, &lesson.UpdatedAt, &contentBytes)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if htmlContent != nil {
		lesson.HTMLContent = *htmlContent
	}

	if contentBytes != nil {
		if err := json.Unmarshal(contentBytes, &lesson.Content); err != nil {
			return nil, err
		}
	}

	return &lesson, nil
}

// Create создает новый урок в базе данных
func (r *LessonRepository) Create(ctx context.Context, courseID string, lesson request.LessonCreate) (*models.LessonDetailed, error) {
	contentJSON, err := json.Marshal(lesson.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal lesson content: %w", err)
	}

	query := `
		INSERT INTO knowledge_base.lesson_d (title, course_id, html_content, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, course_id, html_content, created_at, updated_at, content
	`

	row := r.db.Pool.QueryRow(ctx, query, lesson.Title, courseID, lesson.HTMLContent, contentJSON)

	var newLesson models.LessonDetailed
	var contentBytes []byte
	var htmlContent *string

	err = row.Scan(&newLesson.ID, &newLesson.Title, &newLesson.CourseID, &htmlContent, &newLesson.CreatedAt, &newLesson.UpdatedAt, &contentBytes)
	if err != nil {
		return nil, err
	}

	if htmlContent != nil {
		newLesson.HTMLContent = *htmlContent
	}

	if contentBytes != nil {
		if err := json.Unmarshal(contentBytes, &newLesson.Content); err != nil {
			return nil, err
		}
	}

	return &newLesson, nil
}

// Update обновляет существующий урок в базе данных
func (r *LessonRepository) Update(ctx context.Context, lessonID string, lesson request.LessonUpdate) (*models.LessonDetailed, error) {
	var contentJSON interface{}
	if lesson.Content != nil {
		marshalled, err := json.Marshal(lesson.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal lesson content: %w", err)
		}
		contentJSON = marshalled
	}

	query := `
		UPDATE knowledge_base.lesson_d 
		SET 
			title = COALESCE(NULLIF($1, ''), title),
			html_content = $2,
			content = COALESCE($3, content),
			updated_at = NOW()
		WHERE id = $4
		RETURNING id, title, course_id, html_content, created_at, updated_at, content
	`
	row := r.db.Pool.QueryRow(ctx, query, lesson.Title, lesson.HTMLContent, contentJSON, lessonID)

	var updatedLesson models.LessonDetailed
	var contentBytes []byte
	var htmlContent *string

	err := row.Scan(&updatedLesson.ID, &updatedLesson.Title, &updatedLesson.CourseID, &htmlContent, &updatedLesson.CreatedAt, &updatedLesson.UpdatedAt, &contentBytes)
	if err != nil {
		return nil, err
	}

	if htmlContent != nil {
		updatedLesson.HTMLContent = *htmlContent
	}

	if contentBytes != nil {
		if err := json.Unmarshal(contentBytes, &updatedLesson.Content); err != nil {
			return nil, err
		}
	}

	return &updatedLesson, nil
}

// Delete удаляет урок по его уникальному идентификатору
func (r *LessonRepository) Delete(ctx context.Context, lessonID string) (bool, error) {
	query := `DELETE FROM knowledge_base.lesson_d WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, lessonID)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
