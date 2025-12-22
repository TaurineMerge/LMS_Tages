// Package service contains the business logic layer of the application.
package service

import (
	"context"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CourseService defines the interface for course business logic.
type CourseService interface {
	GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int, level, sortBy string) ([]response.CourseDTO, response.Pagination, error)
	GetCourseByID(ctx context.Context, categoryID, courseID string) (response.CourseDTO, error)
}

type courseService struct {
	repo         repository.CourseRepository
	categoryRepo repository.CategoryRepository
}

// NewCourseService creates a new instance of the course service.
func NewCourseService(repo repository.CourseRepository, categoryRepo repository.CategoryRepository) CourseService {
	return &courseService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// GetCoursesByCategoryID retrieves paginated courses for a category with filters and converts them to DTOs.
func (s *courseService) GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int, level, sortBy string) ([]response.CourseDTO, response.Pagination, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "courseService.GetCoursesByCategoryID")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.Int("page", page),
		attribute.Int("limit", limit),
		attribute.String("level", level),
		attribute.String("sort_by", sortBy),
	)

	// Validate that category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, response.Pagination{}, apperrors.NewNotFound("Category")
		}
		return nil, response.Pagination{}, err
	}

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 
	}

	courses, total, err := s.repo.GetCoursesByCategoryID(ctx, categoryID, page, limit, level, sortBy)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	courseDTOs := make([]response.CourseDTO, 0, len(courses))
	for _, course := range courses {
		courseDTOs = append(courseDTOs, s.mapCourseToDTO(course))
	}

	// Calculate total pages
	pages := int(math.Ceil(float64(total) / float64(limit)))

	pagination := response.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: pages,
	}

	return courseDTOs, pagination, nil
}

// mapCourseToDTO converts a domain Course to a CourseDTO.
func (s *courseService) mapCourseToDTO(course domain.Course) response.CourseDTO {
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

// TruncateDescription truncates a description to a specified number of characters,
// ensuring it doesn't cut in the middle of a word and adds ellipsis if truncated.
func TruncateDescription(text string, maxChars int) string {
	if maxChars <= 0 {
		return ""
	}

	text = strings.TrimSpace(text)

	if utf8.RuneCountInString(text) <= maxChars {
		return text
	}

	// Truncate to maxChars
	truncated := ""
	count := 0
	for _, r := range text {
		if count >= maxChars {
			break
		}
		truncated += string(r)
		count++
	}

	// Find the last space to avoid cutting words
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}

	return truncated + "..."
}

// GetCourseByID retrieves a single course by ID and converts it to DTO.
func (s *courseService) GetCourseByID(ctx context.Context, categoryID, courseID string) (response.CourseDTO, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "courseService.GetCourseByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.String("course_id", courseID),
	)

	// Validate that category exists
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.CourseDTO{}, apperrors.NewNotFound("Category")
		}
		return response.CourseDTO{}, err
	}

	course, err := s.repo.GetCourseByID(ctx, categoryID, courseID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return response.CourseDTO{}, apperrors.NewNotFound("Course")
		}
		return response.CourseDTO{}, err
	}

	return s.mapCourseToDTO(course), nil
}
