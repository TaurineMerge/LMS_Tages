// Package service contains the business logic of the application.
// It orchestrates data from repositories and prepares it for the handler layer.
package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CourseService defines the interface for course-related business logic.
type CourseService interface {
	// GetAllByCategoryID retrieves a paginated list of courses for a specific category, with sorting capabilities.
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int, sort string) ([]response.CourseDTO, response.Pagination, error)
	// GetByID retrieves a single course by its ID.
	GetByID(ctx context.Context, categoryID, courseID string) (response.CourseDTO, error)
}

type courseService struct {
	repo repository.CourseRepository
}

// NewCourseService creates a new instance of a course service.
func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func toCourseDTO(course domain.Course) response.CourseDTO {
	return response.CourseDTO{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Level:       course.Level,
		CategoryID:  course.CategoryID,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}

func (s *courseService) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int, sort string) ([]response.CourseDTO, response.Pagination, error) {
	ctx, span := otel.Tracer("courseService").Start(ctx, "GetAllByCategoryID")
	span.SetAttributes(attribute.String("category.id", categoryID), attribute.String("sort", sort))
	defer span.End()

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	courses, total, err := s.repo.GetAllByCategoryID(ctx, categoryID, page, limit, sort)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	courseDTOs := make([]response.CourseDTO, len(courses))
	for i, course := range courses {
		courseDTOs[i] = toCourseDTO(course)
	}

	pagination := response.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return courseDTOs, pagination, nil
}

func (s *courseService) GetByID(ctx context.Context, categoryID, courseID string) (response.CourseDTO, error) {
	ctx, span := otel.Tracer("courseService").Start(ctx, "GetByID")
	span.SetAttributes(attribute.String("course.id", courseID), attribute.String("category.id", categoryID))
	defer span.End()

	course, err := s.repo.GetByID(ctx, categoryID, courseID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.CourseDTO{}, apperrors.NewNotFound("Course")
		}
		return response.CourseDTO{}, err
	}
	return toCourseDTO(course), nil
}
