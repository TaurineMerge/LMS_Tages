package repositories

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/database"
	"adminPanel/handlers/dto/request"
	"adminPanel/models"

	"github.com/jackc/pgx/v5"
)

// LessonRepository предоставляет методы для работы с уроками.
// Содержит ссылку на базу данных для выполнения запросов.
type LessonRepository struct {
	db *database.Database
}

// NewLessonRepository создает новый экземпляр LessonRepository.
// Принимает соединение с базой данных.
func NewLessonRepository(db *database.Database) *LessonRepository {
	return &LessonRepository{
		db: db,
	}
}

// GetAllByCourseID получает все уроки для заданного курса с пагинацией и сортировкой.
// Принимает courseID, limit, offset, sortBy (title, created_at, updated_at), sortOrder (ASC/DESC).
// Возвращает список уроков.
func (r *LessonRepository) GetAllByCourseID(ctx context.Context, courseID string, limit, offset int, sortBy, sortOrder string) ([]models.Lesson, error) {
	allowedSortBy := map[string]bool{"title": true, "created_at": true, "updated_at": true}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}
	if !(strings.EqualFold(sortOrder, "ASC") || strings.EqualFold(sortOrder, "DESC")) {
		sortOrder = "ASC"
	}

	query := fmt.Sprintf(`
	       SELECT id, title, course_id, content, created_at, updated_at
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
		var content *string
		if err := rows.Scan(&lesson.ID, &lesson.Title, &lesson.CourseID, &content, &lesson.CreatedAt, &lesson.UpdatedAt); err != nil {
			return nil, err
		}
		if content != nil {
			lesson.Content = *content
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// CountByCourseID подсчитывает количество уроков для заданного курса.
// Возвращает количество уроков.
func (r *LessonRepository) CountByCourseID(ctx context.Context, courseID string) (int, error) {
	query := `SELECT COUNT(*) FROM knowledge_base.lesson_d WHERE course_id = $1`
	var count int
	err := r.db.Pool.QueryRow(ctx, query, courseID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByID получает урок по ID.
// Возвращает урок или nil, если не найден.
func (r *LessonRepository) GetByID(ctx context.Context, lessonID string) (*models.Lesson, error) {
	query := `SELECT id, title, course_id, content, created_at, updated_at FROM knowledge_base.lesson_d WHERE id = $1`

	row := r.db.Pool.QueryRow(ctx, query, lessonID)

	var lesson models.Lesson
	var content *string

	err := row.Scan(&lesson.ID, &lesson.Title, &lesson.CourseID, &content, &lesson.CreatedAt, &lesson.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if content != nil {
		lesson.Content = *content
	}

	return &lesson, nil
}

// Create создает новый урок для заданного курса на основе данных из request.LessonCreate.
// Возвращает созданный урок.
func (r *LessonRepository) Create(ctx context.Context, courseID string, lesson request.LessonCreate) (*models.Lesson, error) {
	query := `
	       INSERT INTO knowledge_base.lesson_d (title, course_id, content)
	       VALUES ($1, $2, $3)
	       RETURNING id, title, course_id, content, created_at, updated_at
       `

	row := r.db.Pool.QueryRow(ctx, query, lesson.Title, courseID, lesson.Content)

	var newLesson models.Lesson
	var content *string

	err := row.Scan(&newLesson.ID, &newLesson.Title, &newLesson.CourseID, &content, &newLesson.CreatedAt, &newLesson.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if content != nil {
		newLesson.Content = *content
	}

	return &newLesson, nil
}

// Update обновляет урок по ID на основе данных из request.LessonUpdate.
// Возвращает обновленный урок.
func (r *LessonRepository) Update(ctx context.Context, lessonID string, lesson request.LessonUpdate) (*models.Lesson, error) {
	query := `
	       UPDATE knowledge_base.lesson_d 
	       SET 
		       title = COALESCE(NULLIF($1, ''), title),
		       content = $2,
		       updated_at = NOW()
	       WHERE id = $3
	       RETURNING id, title, course_id, content, created_at, updated_at
       `
	row := r.db.Pool.QueryRow(ctx, query, lesson.Title, lesson.Content, lessonID)

	var updatedLesson models.Lesson
	var content *string

	err := row.Scan(&updatedLesson.ID, &updatedLesson.Title, &updatedLesson.CourseID, &content, &updatedLesson.CreatedAt, &updatedLesson.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if content != nil {
		updatedLesson.Content = *content
	}

	return &updatedLesson, nil
}

// Delete удаляет урок по ID.
// Возвращает true, если урок был удален, false - если не найден.
func (r *LessonRepository) Delete(ctx context.Context, lessonID string) (bool, error) {
	query := `DELETE FROM knowledge_base.lesson_d WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, lessonID)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
