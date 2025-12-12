// Package repository provides the data persistence layer for the application.
// It abstracts the database interactions for domain models.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CourseRepository defines the interface for database operations on courses.
type CourseRepository interface {
	// GetAllByCategoryID retrieves a paginated list of courses for a specific category, with sorting capabilities.
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int, sort string) ([]domain.Course, int, error)
	// GetByID retrieves a single course by its ID, scoped to a category.
	GetByID(ctx context.Context, categoryID, courseID string) (domain.Course, error)
}

type courseRepository struct {
	db *pgxpool.Pool
}

// NewCourseRepository creates a new instance of a course repository.
func NewCourseRepository(db *pgxpool.Pool) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) scanCourse(row pgx.Row) (domain.Course, error) {
	var course domain.Course
	err := row.Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Level,
		&course.Visibility,
		&course.CategoryID,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	return course, err
}

func (r *courseRepository) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int, sort string) ([]domain.Course, int, error) {
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE category_id = $1 AND visibility = 'public'`, courseTable)
	err := r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses by category: %w", err)
	}

	// Determine ORDER BY clause
	orderByClause := buildOrderByClauseForCourses(sort)

	query := fmt.Sprintf(`SELECT id, title, description, level, visibility, category_id, created_at, updated_at FROM %s
		WHERE category_id = $1 AND visibility = 'public'
		%s LIMIT $2 OFFSET $3`, courseTable, orderByClause)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get courses by category: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var course domain.Course
		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Level,
			&course.Visibility,
			&course.CategoryID,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, total, nil
}

func (r *courseRepository) GetByID(ctx context.Context, categoryID, courseID string) (domain.Course, error) {
	query := fmt.Sprintf(`SELECT id, title, description, level, visibility, category_id, created_at, updated_at FROM %s
		WHERE category_id = $1 AND id = $2 AND visibility = 'public'`, courseTable)
	row := r.db.QueryRow(ctx, query, categoryID, courseID)
	course, err := r.scanCourse(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Course{}, fmt.Errorf("course with id %s not found in category %s", courseID, categoryID)
		}
		return domain.Course{}, fmt.Errorf("failed to get course by id: %w", err)
	}

	return course, nil
}

// buildOrderByClauseForCourses constructs the ORDER BY clause based on the sort parameter.
// It uses a whitelist of allowed fields to prevent SQL injection.
// Example sort strings: "title", "-created_at"
func buildOrderByClauseForCourses(sort string) string {
	// Default sort order
	if sort == "" {
		return "ORDER BY created_at ASC"
	}

	// Map of allowed sortable fields from API param to database column
	allowedFields := map[string]string{
		"title":      "title",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	// Determine sort direction and field
	direction := "ASC"
	field := sort
	if strings.HasPrefix(sort, "-") {
		direction = "DESC"
		field = strings.TrimPrefix(sort, "-")
	}

	dbColumn, ok := allowedFields[field]
	if !ok {
		// If an invalid field is provided, fall back to default sort to prevent errors.
		// Alternatively, an error could be returned, but falling back is often more user-friendly for optional sorts.
		return "ORDER BY created_at ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s", dbColumn, direction)
}
