// Package service предоставляет бизнес-логику приложения.
package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// LessonService определяет интерфейс для бизнес-логики, связанной с уроками.
type LessonService interface {
	// GetAllByCourseID получает все уроки для данного курса с пагинацией и сортировкой.
	GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]response.LessonDTO, response.Pagination, error)
	// GetByID получает один урок по его ID.
	GetByID(ctx context.Context, categoryID, courseID, lessonID string) (response.LessonDTODetailed, error)
	// GetNeighboringLessons находит предыдущий и следующий уроки относительно текущего.
	GetNeighboringLessons(ctx context.Context, categoryID, courseID, lessonID string) (prevLesson, nextLesson response.LessonDTO, err error)
}

// lessonService является реализацией LessonService.
type lessonService struct {
	repo repository.LessonRepository
}

// NewLessonService создает новый экземпляр lessonService.
func NewLessonService(repo repository.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

// toLessonDTO преобразует доменную модель Lesson в краткую DTO LessonDTO.
func toLessonDTO(lesson domain.Lesson) response.LessonDTO {
	return response.LessonDTO{
		ID:        lesson.ID,
		Title:     lesson.Title,
		CourseID:  lesson.CourseID,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}

// toLessonDTODetailed преобразует доменную модель Lesson в детальную DTO LessonDTODetailed.
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

// GetAllByCourseID обрабатывает запрос на получение уроков, валидирует параметры,
// вызывает репозиторий и преобразует результат в DTO.
func (s *lessonService) GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]response.LessonDTO, response.Pagination, error) {
	ctx, span := otel.Tracer("lessonService").Start(ctx, "GetAllByCourseID")
	span.SetAttributes(attribute.String("course.id", courseID), attribute.String("category.id", categoryID), attribute.String("sort", sort))
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

	lessons, total, err := s.repo.GetAllByCourseID(ctx, categoryID, courseID, page, limit, sort)
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

// GetByID находит урок по ID. Если урок не найден,
// возвращает стандартизированную ошибку `apperrors.NewNotFound`.
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

// GetNeighboringLessons находит предыдущий и следующий уроки для навигации.
// Сначала получает текущий урок, чтобы использовать его `created_at` как опорную точку,
// затем делает два запроса к репозиторию для получения соседних уроков.
func (s *lessonService) GetNeighboringLessons(ctx context.Context, categoryID, courseID, lessonID string) (response.LessonDTO, response.LessonDTO, error) {
	ctx, span := otel.Tracer("lessonService").Start(ctx, "GetNeighboringLessons")
	span.SetAttributes(attribute.String("lesson.id", lessonID), attribute.String("course.id", courseID))
	defer span.End()

	currentLesson, err := s.repo.GetByID(ctx, categoryID, courseID, lessonID)
	if err != nil {
		return response.LessonDTO{}, response.LessonDTO{}, err
	}

	orderBy := "created_at"

	// Ищем один урок до текущего
	prevLessons, err := s.repo.GetLessonsChunk(ctx, courseID, repository.LessonChunkOptions{
		PivotValue: currentLesson.CreatedAt,
		OrderBy:    orderBy,
		Direction:  repository.DirectionPrevious,
		Limit:      1,
	})
	if err != nil {
		return response.LessonDTO{}, response.LessonDTO{}, err
	}

	// Ищем один урок после текущего
	nextLessons, err := s.repo.GetLessonsChunk(ctx, courseID, repository.LessonChunkOptions{
		PivotValue: currentLesson.CreatedAt,
		OrderBy:    orderBy,
		Direction:  repository.DirectionNext,
		Limit:      1,
	})
	if err != nil {
		return response.LessonDTO{}, response.LessonDTO{}, err
	}

	var prevLessonDTO, nextLessonDTO response.LessonDTO
	if len(prevLessons) > 0 {
		prevLessonDTO = toLessonDTO(prevLessons[0])
	}
	if len(nextLessons) > 0 {
		nextLessonDTO = toLessonDTO(nextLessons[0])
	}

	return prevLessonDTO, nextLessonDTO, nil
}
