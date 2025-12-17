// Package repository provides the data persistence layer for the application.
package repository

import (
	"context"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CourseRepository defines the interface for course data operations.
type CourseRepository interface {
	GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error)
	GetCategoryByID(ctx context.Context, categoryID string) (*domain.Category, error)
}

type courseRepository struct {
	db *pgxpool.Pool
}

// NewCourseRepository creates a new instance of the course repository.
func NewCourseRepository(db *pgxpool.Pool) CourseRepository {
	return &courseRepository{db: db}
}

// GetCoursesByCategoryID retrieves paginated public courses for a given category.
func (r *courseRepository) GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "courseRepository.GetCoursesByCategoryID")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.Int("page", page),
		attribute.Int("limit", limit),
	)

	// Count total courses
	countQuery := `
		SELECT COUNT(*)
		FROM knowledge_base.course_b
		WHERE category_id = $1 AND visibility = 'public'
	`

	var total int
	err := r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to count courses")
		return nil, 0, apperrors.NewDatabaseError("Failed to count courses", err)
	}

	if total == 0 {
		return []domain.Course{}, 0, nil
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated courses
	query := `
		SELECT id, title, description, level, category_id, visibility, created_at, updated_at
		FROM knowledge_base.course_b
		WHERE category_id = $1 AND visibility = 'public'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, categoryID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to query courses")
		return nil, 0, apperrors.NewDatabaseError("Failed to retrieve courses", err)
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
			&course.CategoryID,
			&course.Visibility,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to scan course")
			return nil, 0, apperrors.NewDatabaseError("Failed to scan course", err)
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error iterating courses")
		return nil, 0, apperrors.NewDatabaseError("Error iterating courses", err)
	}

	span.SetAttributes(attribute.Int("courses_count", len(courses)))
	return courses, total, nil
}

// GetCategoryByID retrieves a category by its ID.
func (r *courseRepository) GetCategoryByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "courseRepository.GetCategoryByID")
	defer span.End()

	span.SetAttributes(attribute.String("category_id", categoryID))

	query := `
		SELECT id, title, created_at, updated_at
		FROM knowledge_base.category_d
		WHERE id = $1
	`

	var category domain.Category
	err := r.db.QueryRow(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Title,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			span.SetStatus(codes.Error, "Category not found")
			return nil, apperrors.NewNotFoundError(fmt.Sprintf("Category with ID %s not found", categoryID))
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to query category")
		return nil, apperrors.NewDatabaseError("Failed to retrieve category", err)
	}

	return &category, nil
}
