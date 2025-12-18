// Package repository provides the data persistence layer for the application.
package repository

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CourseRepository defines the interface for course data operations.
type CourseRepository interface {
	GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error)
	GetCourseByID(ctx context.Context, categoryID, courseID string) (domain.Course, error)
}

type courseRepository struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

// NewCourseRepository creates a new instance of the course repository.
func NewCourseRepository(db *pgxpool.Pool) CourseRepository {
	return &courseRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *courseRepository) scanCourse(row scanner) (domain.Course, error) {
	var course domain.Course
	err := row.Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Level,
		&course.CategoryID,
		&course.Visibility,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	return course, err
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

	// Count total public courses
	countQuery := r.psql.Select("COUNT(*)").
		From(courseTable).
		Where(squirrel.Eq{
			"category_id": categoryID,
			"visibility":  "public",
		})

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to build count query")
		return nil, 0, fmt.Errorf("failed to build count query for courses: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, countSql, countArgs...).Scan(&total)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to count courses")
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	if total == 0 {
		return []domain.Course{}, 0, nil
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated courses ordered from oldest to newest
	queryBuilder := r.psql.Select("id", "title", "description", "level", "category_id", "visibility", "created_at", "updated_at").
		From(courseTable).
		Where(squirrel.Eq{
			"category_id": categoryID,
			"visibility":  "public",
		}).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to build query")
		return nil, 0, fmt.Errorf("failed to build get courses query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to query courses")
		return nil, 0, fmt.Errorf("failed to retrieve courses: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		course, err := r.scanCourse(rows)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to scan course")
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Error iterating courses")
		return nil, 0, fmt.Errorf("error iterating courses: %w", err)
	}

	span.SetAttributes(attribute.Int("courses_count", len(courses)))
	return courses, total, nil
}

// GetCourseByID retrieves a single public course by its ID within a specific category.
func (r *courseRepository) GetCourseByID(ctx context.Context, categoryID, courseID string) (domain.Course, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "courseRepository.GetCourseByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.String("course_id", courseID),
	)

	queryBuilder := r.psql.Select("id", "title", "description", "level", "category_id", "visibility", "created_at", "updated_at").
		From(courseTable).
		Where(squirrel.Eq{
			"id":          courseID,
			"category_id": categoryID,
			"visibility":  "public",
		})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to build query")
		return domain.Course{}, fmt.Errorf("failed to build get course by id query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	course, err := r.scanCourse(row)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to scan course")
		return domain.Course{}, fmt.Errorf("failed to get course by id: %w", err)
	}

	return course, nil
}
