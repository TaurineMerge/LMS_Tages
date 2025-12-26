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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CourseService предоставляет бизнес-логику для работы с курсами.
// Содержит репозитории для курсов и категорий, методы для CRUD операций.
type CourseService struct {
	courseRepo   *repositories.CourseRepository
	categoryRepo *repositories.CategoryRepository
}

// courseTracer трассировщик для сервиса курсов.
// Используется для отслеживания операций с курсами.
var courseTracer = otel.Tracer("admin-panel/course-service")

// NewCourseService создает новый экземпляр CourseService.
// Принимает репозитории для курсов и категорий.
func NewCourseService(
	courseRepo *repositories.CourseRepository,
	categoryRepo *repositories.CategoryRepository,
) *CourseService {
	return &CourseService{
		courseRepo:   courseRepo,
		categoryRepo: categoryRepo,
	}
}

// GetCourses получает курсы с фильтрами и пагинацией из request.CourseFilter.
// Возвращает пагинированный ответ с курсами.
func (s *CourseService) GetCourses(ctx context.Context, filter request.CourseFilter) (*response.PaginatedCoursesResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCourses")
	span.SetAttributes(
		attribute.String("filter.level", filter.Level),
		attribute.String("filter.visibility", filter.Visibility),
		attribute.String("filter.category_id", filter.CategoryID),
		attribute.Int("filter.page", filter.Page),
		attribute.Int("filter.limit", filter.Limit),
	)
	defer span.End()

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	categoryExists, err := s.categoryRepo.Exists(ctx, filter.CategoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}
	if !categoryExists {
		return nil, middleware.NotFoundError("Category", filter.CategoryID)
	}

	data, total, err := s.courseRepo.GetFiltered(ctx, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
	}

	courses := make([]models.Course, 0, len(data))
	for _, item := range data {
		course := models.Course{
			BaseModel: models.BaseModel{
				ID:        toString(item["id"]),
				CreatedAt: parseTime(item["created_at"]),
				UpdatedAt: parseTime(item["updated_at"]),
			},
			Title:       toString(item["title"]),
			Description: toString(item["description"]),
			Level:       toString(item["level"]),
			CategoryID:  toString(item["category_id"]),
			Visibility:  toString(item["visibility"]),
			ImageKey:    toString(item["image_key"]),
		}
		courses = append(courses, course)
	}

	pages := (total + filter.Limit - 1) / filter.Limit
	if pages == 0 {
		pages = 1
	}

	return &response.PaginatedCoursesResponse{
		Status: "success",
		Data: struct {
			Items      []models.Course   `json:"items"`
			Pagination models.Pagination `json:"pagination"`
		}{
			Items: courses,
			Pagination: models.Pagination{
				Total: total,
				Page:  filter.Page,
				Limit: filter.Limit,
				Pages: pages,
			},
		},
	}, nil
}

// GetCourse получает курс по ID в заданной категории.
// Возвращает ответ с курсом или ошибку, если не найден.
func (s *CourseService) GetCourse(ctx context.Context, categoryID, id string) (*response.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCourse")
	span.SetAttributes(attribute.String("course.id", id))
	defer span.End()

	data, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get course: %v", err))
	}

	if data == nil {
		return nil, middleware.NotFoundError("Course", id)
	}

	if toString(data["category_id"]) != categoryID {
		return nil, middleware.NotFoundError("Course", id)
	}

	course := &response.CourseResponse{
		Status: "success",
		Data: models.Course{
			BaseModel: models.BaseModel{
				ID:        toString(data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       toString(data["title"]),
			Description: toString(data["description"]),
			Level:       toString(data["level"]),
			CategoryID:  toString(data["category_id"]),
			Visibility:  toString(data["visibility"]),
			ImageKey:    toString(data["image_key"]),
		},
	}

	return course, nil
}

// CreateCourse создает новый курс на основе данных из request.CourseCreate.
// Проверяет существование категории и устанавливает значения по умолчанию.
// Возвращает ответ с созданным курсом.
func (s *CourseService) CreateCourse(ctx context.Context, input request.CourseCreate) (*response.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.CreateCourse")
	span.SetAttributes(
		attribute.String("course.category_id", input.CategoryID),
		attribute.String("course.level", input.Level),
		attribute.String("course.visibility", input.Visibility),
		attribute.String("course.title", input.Title),
	)
	defer span.End()

	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, input.CategoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, middleware.NotFoundError("Category", input.CategoryID)
	}

	if strings.TrimSpace(input.Level) == "" {
		input.Level = "medium"
	}
	if strings.TrimSpace(input.Visibility) == "" {
		input.Visibility = "draft"
	}

	data, err := s.courseRepo.Create(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to create course: %v", err))
	}

	course := &response.CourseResponse{
		Status: "success",
		Data: models.Course{
			BaseModel: models.BaseModel{
				ID:        toString(data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       toString(data["title"]),
			Description: toString(data["description"]),
			Level:       toString(data["level"]),
			CategoryID:  toString(data["category_id"]),
			Visibility:  toString(data["visibility"]),
			ImageKey:    toString(data["image_key"]),
		},
	}

	return course, nil
}

// UpdateCourse обновляет курс по ID в категории на основе данных из request.CourseUpdate.
// Проверяет существование и возвращает ответ с обновленным курсом.
func (s *CourseService) UpdateCourse(ctx context.Context, categoryID, id string, input request.CourseUpdate) (*response.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.UpdateCourse")
	span.SetAttributes(
		attribute.String("course.id", id),
		attribute.String("course.category_id", input.CategoryID),
		attribute.String("course.level", input.Level),
		attribute.String("course.visibility", input.Visibility),
		attribute.String("course.title", input.Title),
	)
	defer span.End()

	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil || toString(existing["category_id"]) != categoryID {
		return nil, middleware.NotFoundError("Course", id)
	}

	input.CategoryID = categoryID

	data, err := s.courseRepo.Update(ctx, id, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to update course: %v", err))
	}

	course := &response.CourseResponse{
		Status: "success",
		Data: models.Course{
			BaseModel: models.BaseModel{
				ID:        toString(data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       toString(data["title"]),
			Description: toString(data["description"]),
			Level:       toString(data["level"]),
			CategoryID:  toString(data["category_id"]),
			Visibility:  toString(data["visibility"]),
			ImageKey:    toString(data["image_key"]),
		},
	}

	return course, nil
}

// DeleteCourse удаляет курс по ID в заданной категории.
// Проверяет существование перед удалением.
func (s *CourseService) DeleteCourse(ctx context.Context, categoryID, id string) error {
	ctx, span := courseTracer.Start(ctx, "CourseService.DeleteCourse")
	span.SetAttributes(attribute.String("course.id", id))
	defer span.End()

	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return middleware.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil || toString(existing["category_id"]) != categoryID {
		return middleware.NotFoundError("Course", id)
	}

	deleted, err := s.courseRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return middleware.InternalError(fmt.Sprintf("Failed to delete course: %v", err))
	}

	if !deleted {
		return middleware.InternalError("Failed to delete course")
	}

	return nil
}

// GetCategoryCourses получает все курсы для заданной категории.
// Возвращает список курсов.
func (s *CourseService) GetCategoryCourses(ctx context.Context, categoryID string) ([]models.Course, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCategoryCourses")
	span.SetAttributes(attribute.String("category.id", categoryID))
	defer span.End()

	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, middleware.NotFoundError("Category", categoryID)
	}

	data, err := s.courseRepo.GetByCategory(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
	}

	courses := make([]models.Course, 0, len(data))
	for _, item := range data {
		course := models.Course{
			BaseModel: models.BaseModel{
				ID:        toString(item["id"]),
				CreatedAt: parseTime(item["created_at"]),
				UpdatedAt: parseTime(item["updated_at"]),
			},
			Title:       toString(item["title"]),
			Description: toString(item["description"]),
			Level:       toString(item["level"]),
			CategoryID:  toString(item["category_id"]),
			Visibility:  toString(item["visibility"]),
		}
		courses = append(courses, course)
	}

	return courses, nil
}
