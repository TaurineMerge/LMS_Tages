package services

import (
	"context"
	"fmt"
	"strconv"

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

// GetLessons - получение уроков по курсу с пагинацией
func (s *LessonService) GetLessons(ctx context.Context, categoryID, courseID, pageStr, limitStr string) ([]models.Lesson, models.Pagination, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.GetLessons")
	span.SetAttributes(
		attribute.String("course.id", courseID),
		attribute.String("category.id", categoryID),
	)
	defer span.End()

	// Проверяем существование курса и его принадлежность категории
	courseData, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, models.Pagination{}, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}
	if courseData == nil || toString(courseData["category_id"]) != categoryID {
		return nil, models.Pagination{}, exceptions.NotFoundError("Course", courseID)
	}

	// Параметры пагинации
	page := parsePositiveInt(pageStr, 1)
	limit := parsePositiveInt(limitStr, 20)
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Подсчитываем количество
	total, err := s.lessonRepo.CountByCourse(ctx, categoryID, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, models.Pagination{}, exceptions.InternalError(fmt.Sprintf("Failed to count lessons: %v", err))
	}

	data, err := s.lessonRepo.GetByCourse(ctx, categoryID, courseID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, models.Pagination{}, exceptions.InternalError(fmt.Sprintf("Failed to get lessons: %v", err))
	}

	lessons := make([]models.Lesson, 0, len(data))
	for _, item := range data {
		lesson := models.Lesson{
			BaseModel: models.BaseModel{
				ID:        toString(item["id"]),
				CreatedAt: parseTime(item["created_at"]),
				UpdatedAt: parseTime(item["updated_at"]),
			},
			Title:      toString(item["title"]),
			CategoryID: toString(item["category_id"]),
			CourseID:   toString(item["course_id"]),
		}
		lessons = append(lessons, lesson)
	}

	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	pagination := models.Pagination{
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}

	return lessons, pagination, nil
}

// GetLesson - получение урока по ID
func (s *LessonService) GetLesson(ctx context.Context, id, courseID, categoryID string) (*models.LessonResponse, error) {
	ctx, span := lessonTracer.Start(ctx, "LessonService.GetLesson")
	span.SetAttributes(
		attribute.String("lesson.id", id),
		attribute.String("course.id", courseID),
		attribute.String("category.id", categoryID),
	)
	defer span.End()

	// Проверяем существование курса и его принадлежность категории
	courseData, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}
	if courseData == nil || toString(courseData["category_id"]) != categoryID {
		return nil, exceptions.NotFoundError("Course", courseID)
	}

	// Получаем урок
	data, err := s.lessonRepo.GetByIDAndCourse(ctx, id, categoryID, courseID)
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
		Data: models.LessonDetailed{
			Lesson: models.Lesson{
				BaseModel: models.BaseModel{
					ID:        toString(parsedData["id"]),
					CreatedAt: parseTime(parsedData["created_at"]),
					UpdatedAt: parseTime(parsedData["updated_at"]),
				},
				Title:      toString(parsedData["title"]),
				CategoryID: toString(parsedData["category_id"]),
				CourseID:   toString(parsedData["course_id"]),
			},
			Content: parsedData["content"].(map[string]interface{}),
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
		attribute.String("category.id", input.CategoryID),
	)
	defer span.End()

	// Проверяем существование курса и принадлежность категории
	courseData, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}
	if courseData == nil {
		return nil, exceptions.NotFoundError("Course", courseID)
	}
	if input.CategoryID == "" {
		input.CategoryID = toString(courseData["category_id"])
	}
	if toString(courseData["category_id"]) != input.CategoryID {
		return nil, exceptions.ValidationError("Category ID does not match course")
	}

	// Создаем урок
	data, err := s.lessonRepo.Create(ctx, courseID, input.CategoryID, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create lesson: %v", err))
	}

	// Парсим контент
	parsedData, _ := s.lessonRepo.ParseContent(data)

	lesson := &models.LessonResponse{
		Data: models.LessonDetailed{
			Lesson: models.Lesson{
				BaseModel: models.BaseModel{
					ID:        toString(parsedData["id"]),
					CreatedAt: parseTime(parsedData["created_at"]),
					UpdatedAt: parseTime(parsedData["updated_at"]),
				},
				Title:      toString(parsedData["title"]),
				CategoryID: toString(parsedData["category_id"]),
				CourseID:   toString(parsedData["course_id"]),
			},
			Content: parsedData["content"].(map[string]interface{}),
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
		attribute.String("category.id", input.CategoryID),
		attribute.String("lesson.title", input.Title),
	)
	defer span.End()

	// Проверяем существование урока
	existing, err := s.lessonRepo.GetByIDAndCourse(ctx, id, input.CategoryID, courseID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check lesson: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Lesson", id)
	}

	// Проверяем категорию/курс
	if existingCat := toString(existing["category_id"]); input.CategoryID != "" && input.CategoryID != existingCat {
		return nil, exceptions.ValidationError("Category ID does not match lesson")
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
		Data: models.LessonDetailed{
			Lesson: models.Lesson{
				BaseModel: models.BaseModel{
					ID:        toString(parsedData["id"]),
					CreatedAt: parseTime(parsedData["created_at"]),
					UpdatedAt: parseTime(parsedData["updated_at"]),
				},
				Title:      toString(parsedData["title"]),
				CategoryID: toString(parsedData["category_id"]),
				CourseID:   toString(parsedData["course_id"]),
			},
			Content: parsedData["content"].(map[string]interface{}),
		},
	}

	return lesson, nil
}

// DeleteLesson - удаление урока
func (s *LessonService) DeleteLesson(ctx context.Context, id, courseID, categoryID string) error {
	ctx, span := lessonTracer.Start(ctx, "LessonService.DeleteLesson")
	span.SetAttributes(
		attribute.String("lesson.id", id),
		attribute.String("course.id", courseID),
		attribute.String("category.id", categoryID),
	)
	defer span.End()

	// Проверяем существование урока
	existing, err := s.lessonRepo.GetByIDAndCourse(ctx, id, categoryID, courseID)
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

// parsePositiveInt parses int with default and minimum 1.
func parsePositiveInt(value string, def int) int {
	if value == "" {
		return def
	}
	if v, err := strconv.Atoi(value); err == nil && v > 0 {
		return v
	}
	return def
}
