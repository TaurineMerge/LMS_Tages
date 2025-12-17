// Package service contains the business logic layer of the application.
package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CourseService defines the interface for course business logic.
type CourseService interface {
	GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int, level, sortBy string) ([]response.CourseDTO, *response.Pagination, error)
	GetCategoryByID(ctx context.Context, categoryID string) (*response.CategoryDTO, error)
}

type courseService struct {
	repo repository.CourseRepository
}

// NewCourseService creates a new instance of the course service.
func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

// GetCoursesByCategoryID retrieves paginated courses for a category with filters and converts them to DTOs.
func (s *courseService) GetCoursesByCategoryID(ctx context.Context, categoryID string, page, limit int, level, sortBy string) ([]response.CourseDTO, *response.Pagination, error) {
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

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 28 // Default limit
	}

	courses, total, err := s.repo.GetCoursesByCategoryID(ctx, categoryID, page, limit, level, sortBy)
	if err != nil {
		return nil, nil, err
	}

	courseDTOs := make([]response.CourseDTO, 0, len(courses))
	for _, course := range courses {
		courseDTOs = append(courseDTOs, s.mapCourseToDTO(course))
	}

	// Calculate total pages
	pages := total / limit
	if total%limit > 0 {
		pages++
	}

	pagination := &response.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: pages,
	}

	return courseDTOs, pagination, nil
}

// GetCategoryByID retrieves a category by ID and converts it to a DTO.
func (s *courseService) GetCategoryByID(ctx context.Context, categoryID string) (*response.CategoryDTO, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "courseService.GetCategoryByID")
	defer span.End()

	span.SetAttributes(attribute.String("category_id", categoryID))

	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return &response.CategoryDTO{
		ID:        category.ID,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}, nil
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
