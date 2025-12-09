package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	lessonsTable = "lessons"
)

type LessonRepository interface {
	GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]domain.Lesson, int, error)
	GetByID(ctx context.Context, id string) (domain.Lesson, error)
}

type lessonRepository struct {
	db *pgxpool.Pool
}

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

func (r *lessonRepository) GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]domain.Lesson, int, error) {
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE course_id = $1", lessonsTable)
	err := r.db.QueryRow(ctx, countQuery, courseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons by course: %w", err)
	}

	query := fmt.Sprintf("SELECT id, title, course_id, content, created_at, updated_at FROM %s WHERE course_id = $1 ORDER BY created_at ASC LIMIT $2 OFFSET $3", lessonsTable)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, courseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get lessons by course: %w", err)
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		// Note: content is not scanned here for the list view, matching LessonDTO
		var lesson domain.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.CourseID,
			&lesson.Content, // Still scan it to avoid scan error, but it won't be used by the DTO
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

func (r *lessonRepository) GetByID(ctx context.Context, id string) (domain.Lesson, error) {
	query := fmt.Sprintf("SELECT id, title, course_id, content, created_at, updated_at FROM %s WHERE id = $1", lessonsTable)
	row := r.db.QueryRow(ctx, query, id)
	lesson, err := r.scanLesson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, fmt.Errorf("lesson with id %s not found", id)
		}
		return domain.Lesson{}, fmt.Errorf("failed to get lesson by id: %w", err)
	}

	return lesson, nil
}
