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

// LessonService defines the interface for lesson-related business logic.
type LessonService interface {
	// GetAllByCourseID retrieves a paginated list of lessons for a specific course.
	GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int) ([]response.LessonDTO, response.Pagination, error)
	// GetByID retrieves a single detailed lesson by its ID.
	GetByID(ctx context.Context, categoryID, courseID, lessonID string) (response.LessonDTODetailed, error)
}

type lessonService struct {
	repo repository.LessonRepository
}

// NewLessonService creates a new instance of a lesson service.
func NewLessonService(repo repository.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func toLessonDTO(lesson domain.Lesson) response.LessonDTO {
	return response.LessonDTO{
		ID:        lesson.ID,
		Title:     lesson.Title,
		CourseID:  lesson.CourseID,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}

func toLessonDTODetailed(lesson domain.Lesson) response.LessonDTODetailed {
	return response.LessonDTODetailed{
		ID:        lesson.ID,
		Title:     lesson.Title,
		CourseID:  lesson.CourseID,
		Content:   lesson.Content,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}

func (s *lessonService) GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int) ([]response.LessonDTO, response.Pagination, error) {
	ctx, span := otel.Tracer("lessonService").Start(ctx, "GetAllByCourseID")
	span.SetAttributes(attribute.String("course.id", courseID), attribute.String("category.id", categoryID))
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

	lessons, total, err := s.repo.GetAllByCourseID(ctx, categoryID, courseID, page, limit)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	lessonDTOs := make([]response.LessonDTO, len(lessons))
	for i, lesson := range lessons {
		lessonDTOs[i] = toLessonDTO(lesson)
	}

	pagination := response.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return lessonDTOs, pagination, nil
}

func (s *lessonService) GetByID(ctx context.Context, categoryID, courseID, lessonID string) (response.LessonDTODetailed, error) {
	ctx, span := otel.Tracer("lessonService").Start(ctx, "GetByID")
	span.SetAttributes(attribute.String("lesson.id", lessonID), attribute.String("course.id", courseID), attribute.String("category.id", categoryID))
	defer span.End()

	lesson, err := s.repo.GetByID(ctx, categoryID, courseID, lessonID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.LessonDTODetailed{}, apperrors.NewNotFound("Lesson")
		}
		return response.LessonDTODetailed{}, err
	}
	return toLessonDTODetailed(lesson), nil
}
