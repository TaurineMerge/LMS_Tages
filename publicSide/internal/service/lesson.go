package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
)

type LessonService interface {
	GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]dto.LessonDTO, domain.Pagination, error)
	GetByID(ctx context.Context, id string) (dto.LessonDTODetailed, error)
}

type lessonService struct {
	repo repository.LessonRepository
}

func NewLessonService(repo repository.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func toLessonDTO(lesson domain.Lesson) dto.LessonDTO {
	return dto.LessonDTO{
		ID:        lesson.ID,
		Title:     lesson.Title,
		CourseID:  lesson.CourseID,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}

func toLessonDTODetailed(lesson domain.Lesson) dto.LessonDTODetailed {
	return dto.LessonDTODetailed{
		ID:        lesson.ID,
		Title:     lesson.Title,
		CourseID:  lesson.CourseID,
		Content:   lesson.Content,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}

func (s *lessonService) GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]dto.LessonDTO, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	lessons, total, err := s.repo.GetAllByCourseID(ctx, courseID, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, apperrors.NewInternal()
	}

	lessonDTOs := make([]dto.LessonDTO, len(lessons))
	for i, lesson := range lessons {
		lessonDTOs[i] = toLessonDTO(lesson)
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return lessonDTOs, pagination, nil
}

func (s *lessonService) GetByID(ctx context.Context, id string) (dto.LessonDTODetailed, error) {
	lesson, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return dto.LessonDTODetailed{}, apperrors.NewNotFound("Lesson")
		}
		return dto.LessonDTODetailed{}, apperrors.NewInternal()
	}
	return toLessonDTODetailed(lesson), nil
}
