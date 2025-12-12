// Package repository provides the data persistence layer for the application.
// It abstracts the database interactions for domain models.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LessonRepository defines the interface for database operations on lessons.
type LessonRepository interface {
	// GetAllByCourseID retrieves a paginated list of lessons for a specific course and category, with sorting capabilities.
	GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]domain.Lesson, int, error)
	// GetByID retrieves a single lesson by its ID, scoped to a course and category.
	GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error)
}

type lessonRepository struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

// NewLessonRepository creates a new instance of a lesson repository.
func NewLessonRepository(db *pgxpool.Pool) LessonRepository {
	return &lessonRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// pgx.Row and *pgx.Rows both implement this interface
type scanner interface {
	Scan(dest ...any) error
}

func (r *lessonRepository) scanLesson(row scanner) (domain.Lesson, error) {
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

func (r *lessonRepository) GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]domain.Lesson, int, error) {
	countBuilder := r.psql.Select("COUNT(l.id)").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"c.visibility":  "public",
		})

	countQuery, args, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query for lessons: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons by course: %w", err)
	}

	if total == 0 {
		return []domain.Lesson{}, 0, nil
	}

	queryBuilder := r.psql.Select("l.id", "l.title", "l.course_id", "l.content", "l.created_at", "l.updated_at").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"c.visibility":  "public",
		}).
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))

	queryBuilder = r.applySorting(queryBuilder, sort)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build get all lessons query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get lessons by course: %w", err)
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		lesson, err := r.scanLesson(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, lesson)
	}

	return lessons, total, nil
}

func (r *lessonRepository) GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error) {
	queryBuilder := r.psql.Select("l.id", "l.title", "l.course_id", "l.content", "l.created_at", "l.updated_at").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"l.id":          lessonID,
			"c.visibility":  "public",
		})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("failed to build get lesson by id query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	lesson, err := r.scanLesson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, fmt.Errorf("lesson with id %s not found in course %s", lessonID, courseID)
		}
		return domain.Lesson{}, fmt.Errorf("failed to get lesson by id: %w", err)
	}

	return lesson, nil
}

func (r *lessonRepository) applySorting(builder squirrel.SelectBuilder, sort string) squirrel.SelectBuilder {
	if sort == "" {
		return builder.OrderBy("l.created_at ASC")
	}

	allowedFields := map[string]string{
		"title":      "l.title",
		"created_at": "l.created_at",
		"updated_at": "l.updated_at",
	}

	direction := "ASC"
	if strings.HasPrefix(sort, "-") {
		direction = "DESC"
		sort = strings.TrimPrefix(sort, "-")
	}

	dbColumn, ok := allowedFields[sort]
	if !ok {
		return builder.OrderBy("l.created_at ASC")
	}

	return builder.OrderBy(fmt.Sprintf("%s %s", dbColumn, direction))
}
