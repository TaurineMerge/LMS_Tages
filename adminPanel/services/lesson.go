package services

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/handlers/dto/request"
	"adminPanel/handlers/dto/response"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/repositories"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// LessonService предоставляет бизнес-логику для работы с уроками.
// Содержит репозитории для уроков и курсов, методы для CRUD операций.
type LessonService struct {
	lessonRepo   *repositories.LessonRepository
	courseRepo   *repositories.CourseRepository
	lessonTracer trace.Tracer
}

// NewLessonService создает новый экземпляр LessonService.
// Принимает репозитории для уроков и курсов, инициализирует трассировщик.
func NewLessonService(
	lessonRepo *repositories.LessonRepository,
	courseRepo *repositories.CourseRepository,
) *LessonService {
	return &LessonService{
		lessonRepo:   lessonRepo,
		courseRepo:   courseRepo,
		lessonTracer: otel.Tracer("admin-panel/lesson-service"),
	}
}

// GetLessons получает уроки для заданного курса с пагинацией и сортировкой из models.QueryList.
// Проверяет существование курса и возвращает пагинированный ответ с уроками.
func (s *LessonService) GetLessons(ctx context.Context, courseID string, queryParams models.QueryList) (*response.LessonListResponse, error) {
	ctx, span := s.lessonTracer.Start(ctx, "LessonService.GetLessons")
	defer span.End()

	if queryParams.Page < 1 {
		queryParams.Page = 1
	}
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		queryParams.Limit = 20
	}

	courseExists, err := s.courseRepo.Exists(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check course existence: %v", err))
	}
	if !courseExists {
		return nil, middleware.NotFoundError("Course", courseID)
	}

	sortBy, sortOrder := parseSortParameter(queryParams.Sort)
	offset := (queryParams.Page - 1) * queryParams.Limit

	total, err := s.lessonRepo.CountByCourseID(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to count lessons: %v", err))
	}

	lessons, err := s.lessonRepo.GetAllByCourseID(ctx, courseID, queryParams.Limit, offset, sortBy, sortOrder)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get lessons: %v", err))
	}

	pages := 0
	if queryParams.Limit > 0 {
		pages = (total + queryParams.Limit - 1) / queryParams.Limit
	}
	if pages == 0 {
		pages = 1
	}

	return &response.LessonListResponse{
		Status: "success",
		Data: models.ResponsePaginationLessonsList{
			Items: lessons,
			Pagination: models.Pagination{
				Total: total,
				Page:  queryParams.Page,
				Limit: queryParams.Limit,
				Pages: pages,
			},
		},
	}, nil
}

// GetLesson получает урок по ID в заданном курсе.
// Возвращает ответ с уроком или ошибку, если не найден.
func (s *LessonService) GetLesson(ctx context.Context, lessonID, courseID string) (*response.LessonResponse, error) {
	ctx, span := s.lessonTracer.Start(ctx, "LessonService.GetLesson")
	defer span.End()

	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		span.RecordError(err)
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get lesson: %v", err))
	}

	if lesson == nil || lesson.CourseID != courseID {
		return nil, middleware.NotFoundError("Lesson", lessonID)
	}

	fmt.Printf("[DEBUG] GetLesson: lessonID=%s, courseID=%s, title=%s, content length=%d\n",
		lessonID, courseID, lesson.Title, len(lesson.Content))

	return &response.LessonResponse{
		Status: "success",
		Data:   *lesson,
	}, nil
}

// CreateLesson создает новый урок для заданного курса на основе данных из request.LessonCreate.
// Проверяет существование курса и возвращает ответ с созданным уроком.
func (s *LessonService) CreateLesson(ctx context.Context, courseID string, input request.LessonCreate) (*response.LessonResponse, error) {
	ctx, span := s.lessonTracer.Start(ctx, "LessonService.CreateLesson")
	defer span.End()

	courseExists, err := s.courseRepo.Exists(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check course existence: %v", err))
	}
	if !courseExists {
		return nil, middleware.NotFoundError("Course", courseID)
	}

	lesson, err := s.lessonRepo.Create(ctx, courseID, input)
	if err != nil {
		span.RecordError(err)
		return nil, middleware.InternalError(fmt.Sprintf("Failed to create lesson: %v", err))
	}

	return &response.LessonResponse{
		Status: "success",
		Data:   *lesson,
	}, nil
}

// UpdateLesson обновляет урок по ID в курсе на основе данных из request.LessonUpdate.
// Проверяет существование и возвращает ответ с обновленным уроком.
func (s *LessonService) UpdateLesson(ctx context.Context, lessonID, courseID string, input request.LessonUpdate) (*response.LessonResponse, error) {
	ctx, span := s.lessonTracer.Start(ctx, "LessonService.UpdateLesson")
	defer span.End()

	existing, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		span.RecordError(err)
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check lesson: %v", err))
	}
	if existing == nil || existing.CourseID != courseID {
		return nil, middleware.NotFoundError("Lesson", lessonID)
	}

	lesson, err := s.lessonRepo.Update(ctx, lessonID, input)
	if err != nil {
		span.RecordError(err)
		return nil, middleware.InternalError(fmt.Sprintf("Failed to update lesson: %v", err))
	}

	return &response.LessonResponse{
		Status: "success",
		Data:   *lesson,
	}, nil
}

// DeleteLesson удаляет урок по ID в заданном курсе.
// Проверяет существование перед удалением.
func (s *LessonService) DeleteLesson(ctx context.Context, lessonID, courseID string) error {
	ctx, span := s.lessonTracer.Start(ctx, "LessonService.DeleteLesson")
	defer span.End()

	existing, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		span.RecordError(err)
		return middleware.InternalError(fmt.Sprintf("Failed to check lesson: %v", err))
	}
	if existing == nil || existing.CourseID != courseID {
		return middleware.NotFoundError("Lesson", lessonID)
	}

	deleted, err := s.lessonRepo.Delete(ctx, lessonID)
	if err != nil {
		span.RecordError(err)
		return middleware.InternalError(fmt.Sprintf("Failed to delete lesson: %v", err))
	}
	if !deleted {
		return middleware.InternalError("Failed to delete lesson for an unknown reason")
	}

	return nil
}

// parseSortParameter разбирает параметр сортировки.
// Если начинается с "-", то DESC, иначе ASC. По умолчанию "created_at ASC".
func parseSortParameter(sort string) (sortBy, sortOrder string) {
	if sort == "" {
		return "created_at", "ASC"
	}
	if strings.HasPrefix(sort, "-") {
		return strings.TrimPrefix(sort, "-"), "DESC"
	}
	return sort, "ASC"
}
