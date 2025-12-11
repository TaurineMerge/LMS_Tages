// Package repository provides the data persistence layer for the application.
// It abstracts the database interactions for domain models.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LessonRepository defines the interface for database operations on lessons.
type LessonRepository interface {
	// GetAllByCourseID retrieves a paginated list of lessons for a specific course and category.
	GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int) ([]domain.Lesson, int, error)
	// GetByID retrieves a single lesson by its ID, scoped to a course and category.
	GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error)
}

type lessonRepository struct {
	db *pgxpool.Pool
}

// NewLessonRepository creates a new instance of a lesson repository.
func NewLessonRepository(db *pgxpool.Pool) LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) scanLesson(row pgx.Row) (domain.Lesson, error) {
	var lesson domain.Lesson
	err := row.Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&lesson.Content,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	return lesson, err
}

func (r *lessonRepository) GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int) ([]domain.Lesson, int, error) {
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(l.*) FROM %s l
		JOIN %s c ON l.course_id = c.id
		WHERE c.category_id = $1 AND l.course_id = $2 AND c.visibility = 'public'`, lessonsTable, courseTable)
	err := r.db.QueryRow(ctx, countQuery, categoryID, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons by course: %w", err)
	}

	query := fmt.Sprintf(`SELECT l.id, l.title, l.course_id, l.content, l.created_at, l.updated_at FROM %s l
		JOIN %s c ON l.course_id = c.id
		WHERE c.category_id = $1 AND l.course_id = $2 AND c.visibility = 'public'
		ORDER BY l.created_at ASC LIMIT $3 OFFSET $4`, lessonsTable, courseTable)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, categoryID, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get lessons by course: %w", err)
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		var lesson domain.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.CourseID,
			&lesson.Content,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, total, nil
}

func (r *lessonRepository) GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error) {
	query := fmt.Sprintf(`SELECT l.id, l.title, l.course_id, l.content, l.created_at, l.updated_at FROM %s l
		JOIN %s c ON l.course_id = c.id
		WHERE c.category_id = $1 AND l.course_id = $2 AND l.id = $3`, lessonsTable, courseTable)
	row := r.db.QueryRow(ctx, query, categoryID, courseID, lessonID)
	lesson, err := r.scanLesson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, fmt.Errorf("lesson with id %s not found in course %s", lessonID, courseID)
		}
		return domain.Lesson{}, fmt.Errorf("failed to get lesson by id: %w", err)
	}

	return lesson, nil
}
