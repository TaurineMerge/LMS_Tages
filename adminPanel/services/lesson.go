package services

import (
	"context"
	"fmt"

	"adminPanel/exceptions"
	"adminPanel/models"
	"adminPanel/repositories"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LessonService - сервис для работы с уроками
type LessonService struct {
	lessonRepo *repositories.LessonRepository
	courseRepo *repositories.CourseRepository
}

var lessonTracer = otel.Tracer("admin-panel/lesson-service")

// NewLessonService создает сервис уроков
func NewLessonService(
	lessonRepo *repositories.LessonRepository,
	courseRepo *repositories.CourseRepository,
) *LessonService {
	return &LessonService{
		lessonRepo: lessonRepo,
		courseRepo: courseRepo,
	}
}

// GetLessons - получение уроков по курсу
func (s *LessonService) GetLessons(ctx context.Context, courseID string) ([]models.LessonResponse, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.GetLessons")
	span.SetAttributes(attribute.String("course.id", courseID))
	defer span.End()

	// Проверяем существование курса
	courseExists, err := s.courseRepo.Exists(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if !courseExists {
		return nil, exceptions.NotFoundError("Course", courseID)
	}

	// Получаем уроки
	data, err := s.lessonRepo.GetByCourse(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get lessons: %v", err))
	}

	lessons := make([]models.LessonResponse, 0, len(data))
	for _, item := range data {
		// Парсим контент
		parsedData, _ := s.lessonRepo.ParseContent(item)
		
		lesson := models.LessonResponse{
			Lesson: models.Lesson{
				BaseModel: models.BaseModel{
					ID:        fmt.Sprintf("%v", parsedData["id"]),
					CreatedAt: parseTime(parsedData["created_at"]),
					UpdatedAt: parseTime(parsedData["updated_at"]),
				},
				Title:    fmt.Sprintf("%v", parsedData["title"]),
				CourseID: fmt.Sprintf("%v", parsedData["course_id"]),
				Content:  parsedData["content"].(map[string]interface{}),
			},
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// GetLesson - получение урока по ID
func (s *LessonService) GetLesson(ctx context.Context, id, courseID string) (*models.LessonResponse, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.GetLesson")
	span.SetAttributes(
		attribute.String("lesson.id", id),
		attribute.String("course.id", courseID),
	)
	defer span.End()

	// Проверяем существование курса
	courseExists, err := s.courseRepo.Exists(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if !courseExists {
		return nil, exceptions.NotFoundError("Course", courseID)
	}

	// Получаем урок
	data, err := s.lessonRepo.GetByIDAndCourse(ctx, id, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get lesson: %v", err))
	}

	if data == nil {
		return nil, exceptions.NotFoundError("Lesson", id)
	}

	// Парсим контент
	parsedData, _ := s.lessonRepo.ParseContent(data)

	lesson := &models.LessonResponse{
		Lesson: models.Lesson{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", parsedData["id"]),
				CreatedAt: parseTime(parsedData["created_at"]),
				UpdatedAt: parseTime(parsedData["updated_at"]),
			},
			Title:    fmt.Sprintf("%v", parsedData["title"]),
			CourseID: fmt.Sprintf("%v", parsedData["course_id"]),
			Content:  parsedData["content"].(map[string]interface{}),
		},
	}

	return lesson, nil
}

// CreateLesson - создание урока
func (s *LessonService) CreateLesson(ctx context.Context, courseID string, input models.LessonCreate) (*models.LessonResponse, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.CreateLesson")
	span.SetAttributes(
		attribute.String("course.id", courseID),
		attribute.String("lesson.title", input.Title),
	)
	defer span.End()

	// Проверяем существование курса
	courseExists, err := s.courseRepo.Exists(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if !courseExists {
		return nil, exceptions.NotFoundError("Course", courseID)
	}

	// Создаем урок
	data, err := s.lessonRepo.Create(ctx, courseID, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create lesson: %v", err))
	}

	// Парсим контент
	parsedData, _ := s.lessonRepo.ParseContent(data)

	lesson := &models.LessonResponse{
		Lesson: models.Lesson{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", parsedData["id"]),
				CreatedAt: parseTime(parsedData["created_at"]),
				UpdatedAt: parseTime(parsedData["updated_at"]),
			},
			Title:    fmt.Sprintf("%v", parsedData["title"]),
			CourseID: fmt.Sprintf("%v", parsedData["course_id"]),
			Content:  parsedData["content"].(map[string]interface{}),
		},
	}

	return lesson, nil
}

// UpdateLesson - обновление урока
func (s *LessonService) UpdateLesson(ctx context.Context, id, courseID string, input models.LessonUpdate) (*models.LessonResponse, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.UpdateLesson")
	span.SetAttributes(
		attribute.String("lesson.id", id),
		attribute.String("course.id", courseID),
		attribute.String("lesson.title", input.Title),
	)
	defer span.End()

	// Проверяем существование урока
	existing, err := s.lessonRepo.GetByIDAndCourse(ctx, id, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check lesson: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Lesson", id)
	}

	// Обновляем урок
	data, err := s.lessonRepo.Update(ctx, id, courseID, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to update lesson: %v", err))
	}

	// Парсим контент
	parsedData, _ := s.lessonRepo.ParseContent(data)

	lesson := &models.LessonResponse{
		Lesson: models.Lesson{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", parsedData["id"]),
				CreatedAt: parseTime(parsedData["created_at"]),
				UpdatedAt: parseTime(parsedData["updated_at"]),
			},
			Title:    fmt.Sprintf("%v", parsedData["title"]),
			CourseID: fmt.Sprintf("%v", parsedData["course_id"]),
			Content:  parsedData["content"].(map[string]interface{}),
		},
	}

	return lesson, nil
}

// DeleteLesson - удаление урока
func (s *LessonService) DeleteLesson(ctx context.Context, id, courseID string) error {
	ctx, span := lessonTracer.Start(ctx, "LessonService.DeleteLesson")
	span.SetAttributes(
		attribute.String("lesson.id", id),
		attribute.String("course.id", courseID),
	)
	defer span.End()

	// Проверяем существование урока
	existing, err := s.lessonRepo.GetByIDAndCourse(ctx, id, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check lesson: %v", err))
	}

	if existing == nil {
		return exceptions.NotFoundError("Lesson", id)
	}

	// Удаляем урок
	deleted, err := s.lessonRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to delete lesson: %v", err))
	}

	if !deleted {
		return exceptions.InternalError("Failed to delete lesson")
	}

	return nil
}