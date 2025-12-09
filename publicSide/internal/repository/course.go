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
	coursesTable = "courses"
)

type CourseRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Course, int, error)
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error)
	GetByID(ctx context.Context, id string) (domain.Course, error)
}

type courseRepository struct {
	db *pgxpool.Pool
}

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

func (r *courseRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Course, int, error) {
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE visibility = 'public'", coursesTable)
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	query := fmt.Sprintf("SELECT id, title, description, level, visibility, category_id, created_at, updated_at FROM %s WHERE visibility = 'public' ORDER BY created_at DESC LIMIT $1 OFFSET $2", coursesTable)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all courses: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		course, err := r.scanCourse(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, total, nil
}

func (r *courseRepository) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error) {
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE category_id = $1 AND visibility = 'public'", coursesTable)
	err := r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses by category: %w", err)
	}

	query := fmt.Sprintf("SELECT id, title, description, level, visibility, category_id, created_at, updated_at FROM %s WHERE category_id = $1 AND visibility = 'public' ORDER BY created_at DESC LIMIT $2 OFFSET $3", coursesTable)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get courses by category: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		course, err := r.scanCourse(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	return courses, total, nil
}

func (r *courseRepository) GetByID(ctx context.Context, id string) (domain.Course, error) {
	query := fmt.Sprintf("SELECT id, title, description, level, visibility, category_id, created_at, updated_at FROM %s WHERE id = $1 AND visibility = 'public'", coursesTable)
	row := r.db.QueryRow(ctx, query, id)
	course, err := r.scanCourse(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Course{}, fmt.Errorf("course with id %s not found", id)
		}
		return domain.Course{}, fmt.Errorf("failed to get course by id: %w", err)
	}

	return course, nil
}
