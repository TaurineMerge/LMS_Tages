package services

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/exceptions"
	"adminPanel/models"
	"adminPanel/repositories"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CourseService - сервис для работы с курсами
//
// Сервис предоставляет методы для управления курсами:
//   - Получение курсов с фильтрацией и пагинацией
//   - Получение курса по ID
//   - Создание нового курса
//   - Обновление курса
//   - Удаление курса
//   - Получение курсов категории
//
// Особенности:
//   - Проверка существования категорий
//   - Фильтрация по уровню сложности и видимости
//   - Интеграция с OpenTelemetry для трассировки
//   - Валидация данных
type CourseService struct {
	courseRepo   *repositories.CourseRepository
	categoryRepo *repositories.CategoryRepository
}

var courseTracer = otel.Tracer("admin-panel/course-service")

// NewCourseService создает новый сервис для работы с курсами
//
// Параметры:
//   - courseRepo: репозиторий для работы с курсами
//   - categoryRepo: репозиторий для работы с категориями
//
// Возвращает:
//   - *CourseService: указатель на новый сервис
func NewCourseService(
	courseRepo *repositories.CourseRepository,
	categoryRepo *repositories.CategoryRepository,
) *CourseService {
	return &CourseService{
		courseRepo:   courseRepo,
		categoryRepo: categoryRepo,
	}
}

// GetCourses получает курсы с фильтрацией и пагинацией
//
// Метод возвращает курсы указанной категории с возможностью
// фильтрации по уровню сложности и видимости.
//
// Параметры:
//   - ctx: контекст выполнения
//   - filter: фильтр для поиска курсов
//
// Возвращает:
//   - *models.PaginatedCoursesResponse: ответ с курсами и пагинацией
//   - error: ошибка выполнения (если есть)
func (s *CourseService) GetCourses(ctx context.Context, filter models.CourseFilter) (*models.PaginatedCoursesResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCourses")
	span.SetAttributes(
		attribute.String("filter.level", filter.Level),
		attribute.String("filter.visibility", filter.Visibility),
		attribute.String("filter.category_id", filter.CategoryID),
		attribute.Int("filter.page", filter.Page),
		attribute.Int("filter.limit", filter.Limit),
	)
	defer span.End()

	// Устанавливаем значения по умолчанию
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	// Проверяем категорию
	categoryExists, err := s.categoryRepo.Exists(ctx, filter.CategoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}
	if !categoryExists {
		return nil, exceptions.NotFoundError("Category", filter.CategoryID)
	}

	data, total, err := s.courseRepo.GetFiltered(ctx, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
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

	pages := (total + filter.Limit - 1) / filter.Limit
	if pages == 0 {
		pages = 1
	}

	return &models.PaginatedCoursesResponse{
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

// GetCourse получает курс по уникальному идентификатору
//
// Проверяет принадлежность курса указанной категории.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//   - id: уникальный идентификатор курса
//
// Возвращает:
//   - *models.CourseResponse: ответ с курсом
//   - error: ошибка выполнения (если есть)
func (s *CourseService) GetCourse(ctx context.Context, categoryID, id string) (*models.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCourse")
	span.SetAttributes(attribute.String("course.id", id))
	defer span.End()

	data, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get course: %v", err))
	}

	if data == nil {
		return nil, exceptions.NotFoundError("Course", id)
	}

	if toString(data["category_id"]) != categoryID {
		return nil, exceptions.NotFoundError("Course", id)
	}

	course := &models.CourseResponse{
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
		},
	}

	return course, nil
}

// CreateCourse создает новый курс в указанной категории
//
// Устанавливает значения по умолчанию для уровня и видимости.
//
// Параметры:
//   - ctx: контекст выполнения
//   - input: данные для создания курса
//
// Возвращает:
//   - *models.CourseResponse: ответ с созданным курсом
//   - error: ошибка выполнения (если есть)
func (s *CourseService) CreateCourse(ctx context.Context, input models.CourseCreate) (*models.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.CreateCourse")
	span.SetAttributes(
		attribute.String("course.category_id", input.CategoryID),
		attribute.String("course.level", input.Level),
		attribute.String("course.visibility", input.Visibility),
		attribute.String("course.title", input.Title),
	)
	defer span.End()

	// Проверяем существование категории
	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, input.CategoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, exceptions.NotFoundError("Category", input.CategoryID)
	}

	// Устанавливаем значения по умолчанию
	if strings.TrimSpace(input.Level) == "" {
		input.Level = "medium"
	}
	if strings.TrimSpace(input.Visibility) == "" {
		input.Visibility = "draft"
	}

	// Создаем курс
	data, err := s.courseRepo.Create(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create course: %v", err))
	}

	course := &models.CourseResponse{
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
		},
	}

	return course, nil
}

// UpdateCourse обновляет существующий курс
//
// Проверяет принадлежность курса указанной категории.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//   - id: уникальный идентификатор курса
//   - input: данные для обновления курса
//
// Возвращает:
//   - *models.CourseResponse: ответ с обновленным курсом
//   - error: ошибка выполнения (если есть)
func (s *CourseService) UpdateCourse(ctx context.Context, categoryID, id string, input models.CourseUpdate) (*models.CourseResponse, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.UpdateCourse")
	span.SetAttributes(
		attribute.String("course.id", id),
		attribute.String("course.category_id", input.CategoryID),
		attribute.String("course.level", input.Level),
		attribute.String("course.visibility", input.Visibility),
		attribute.String("course.title", input.Title),
	)
	defer span.End()

	// Проверяем существование курса
	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil || toString(existing["category_id"]) != categoryID {
		return nil, exceptions.NotFoundError("Course", id)
	}

	// Категорию не меняем для соответствия пути
	input.CategoryID = categoryID

	// Обновляем курс
	data, err := s.courseRepo.Update(ctx, id, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to update course: %v", err))
	}

	course := &models.CourseResponse{
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
		},
	}

	return course, nil
}

// DeleteCourse удаляет курс по уникальному идентификатору
//
// Проверяет принадлежность курса указанной категории.
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//   - id: уникальный идентификатор курса
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (s *CourseService) DeleteCourse(ctx context.Context, categoryID, id string) error {
	ctx, span := courseTracer.Start(ctx, "CourseService.DeleteCourse")
	span.SetAttributes(attribute.String("course.id", id))
	defer span.End()

	// Проверяем существование курса
	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil || toString(existing["category_id"]) != categoryID {
		return exceptions.NotFoundError("Course", id)
	}

	// TODO: Проверяем наличие связанных уроков

	// Удаляем курс
	deleted, err := s.courseRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to delete course: %v", err))
	}

	if !deleted {
		return exceptions.InternalError("Failed to delete course")
	}

	return nil
}

// GetCategoryCourses получает все курсы указанной категории
//
// Параметры:
//   - ctx: контекст выполнения
//   - categoryID: уникальный идентификатор категории
//
// Возвращает:
//   - []models.Course: список курсов категории
//   - error: ошибка выполнения (если есть)
func (s *CourseService) GetCategoryCourses(ctx context.Context, categoryID string) ([]models.Course, error) {
	ctx, span := courseTracer.Start(ctx, "CourseService.GetCategoryCourses")
	span.SetAttributes(attribute.String("category.id", categoryID))
	defer span.End()

	// Проверяем существование категории
	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, exceptions.NotFoundError("Category", categoryID)
	}

	// Получаем курсы
	data, err := s.courseRepo.GetByCategory(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
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
