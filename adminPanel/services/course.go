package services

import (
	"context"
	"fmt"

	"adminPanel/exceptions"
	"adminPanel/models"
	"adminPanel/repositories"
)

// CourseService - сервис для работы с курсами
type CourseService struct {
	courseRepo    *repositories.CourseRepository
	categoryRepo  *repositories.CategoryRepository
}

// NewCourseService создает сервис курсов
func NewCourseService(
	courseRepo *repositories.CourseRepository,
	categoryRepo *repositories.CategoryRepository,
) *CourseService {
	return &CourseService{
		courseRepo:   courseRepo,
		categoryRepo: categoryRepo,
	}
}

// GetCourses - получение курсов с фильтрацией
func (s *CourseService) GetCourses(ctx context.Context, filter models.CourseFilter) (*models.PaginatedCourseResponse, error) {
	// Устанавливаем значения по умолчанию
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	data, total, err := s.courseRepo.GetFiltered(ctx, filter)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
	}

	courses := make([]models.CourseResponse, 0, len(data))
	for _, item := range data {
		course := models.Course{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", item["id"]),
				CreatedAt: parseTime(item["created_at"]),
				UpdatedAt: parseTime(item["updated_at"]),
			},
			Title:       fmt.Sprintf("%v", item["title"]),
			Description: fmt.Sprintf("%v", item["description"]),
			Level:       fmt.Sprintf("%v", item["level"]),
			CategoryID:  fmt.Sprintf("%v", item["category_id"]),
			Visibility:  fmt.Sprintf("%v", item["visibility"]),
		}
		courses = append(courses, models.CourseResponse{Course: course})
	}

	pages := (total + filter.Limit - 1) / filter.Limit
	if pages == 0 {
		pages = 1
	}

	return &models.PaginatedCourseResponse{
		Data:  courses,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
		Pages: pages,
	}, nil
}

// GetCourse - получение курса по ID
func (s *CourseService) GetCourse(ctx context.Context, id string) (*models.CourseResponse, error) {
	data, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get course: %v", err))
	}

	if data == nil {
		return nil, exceptions.NotFoundError("Course", id)
	}

	course := &models.CourseResponse{
		Course: models.Course{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       fmt.Sprintf("%v", data["title"]),
			Description: fmt.Sprintf("%v", data["description"]),
			Level:       fmt.Sprintf("%v", data["level"]),
			CategoryID:  fmt.Sprintf("%v", data["category_id"]),
			Visibility:  fmt.Sprintf("%v", data["visibility"]),
		},
	}

	return course, nil
}

// CreateCourse - создание курса
func (s *CourseService) CreateCourse(ctx context.Context, input models.CourseCreate) (*models.CourseResponse, error) {
	// Проверяем существование категории
	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, input.CategoryID)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, exceptions.NotFoundError("Category", input.CategoryID)
	}

	// Создаем курс
	data, err := s.courseRepo.Create(ctx, input)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create course: %v", err))
	}

	course := &models.CourseResponse{
		Course: models.Course{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       fmt.Sprintf("%v", data["title"]),
			Description: fmt.Sprintf("%v", data["description"]),
			Level:       fmt.Sprintf("%v", data["level"]),
			CategoryID:  fmt.Sprintf("%v", data["category_id"]),
			Visibility:  fmt.Sprintf("%v", data["visibility"]),
		},
	}

	return course, nil
}

// UpdateCourse - обновление курса
func (s *CourseService) UpdateCourse(ctx context.Context, id string, input models.CourseUpdate) (*models.CourseResponse, error) {
	// Проверяем существование курса
	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Course", id)
	}

	// Если меняется категория, проверяем её существование
	if input.CategoryID != "" && input.CategoryID != fmt.Sprintf("%v", existing["category_id"]) {
		categoryExists, err := s.courseRepo.ExistsByCategory(ctx, input.CategoryID)
		if err != nil {
			return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
		}

		if !categoryExists {
			return nil, exceptions.NotFoundError("Category", input.CategoryID)
		}
	}

	// Обновляем курс
	data, err := s.courseRepo.Update(ctx, id, input)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to update course: %v", err))
	}

	course := &models.CourseResponse{
		Course: models.Course{
			BaseModel: models.BaseModel{
				ID:        fmt.Sprintf("%v", data["id"]),
				CreatedAt: parseTime(data["created_at"]),
				UpdatedAt: parseTime(data["updated_at"]),
			},
			Title:       fmt.Sprintf("%v", data["title"]),
			Description: fmt.Sprintf("%v", data["description"]),
			Level:       fmt.Sprintf("%v", data["level"]),
			CategoryID:  fmt.Sprintf("%v", data["category_id"]),
			Visibility:  fmt.Sprintf("%v", data["visibility"]),
		},
	}

	return course, nil
}

// DeleteCourse - удаление курса
func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	// Проверяем существование курса
	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return exceptions.InternalError(fmt.Sprintf("Failed to check course: %v", err))
	}

	if existing == nil {
		return exceptions.NotFoundError("Course", id)
	}

	// TODO: Проверяем наличие связанных уроков

	// Удаляем курс
	deleted, err := s.courseRepo.Delete(ctx, id)
	if err != nil {
		return exceptions.InternalError(fmt.Sprintf("Failed to delete course: %v", err))
	}

	if !deleted {
		return exceptions.InternalError("Failed to delete course")
	}

	return nil
}

// GetCategoryCourses - получение курсов категории
func (s *CourseService) GetCategoryCourses(ctx context.Context, categoryID string) ([]models.CourseResponse, error) {
	// Проверяем существование категории
	categoryExists, err := s.courseRepo.ExistsByCategory(ctx, categoryID)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if !categoryExists {
		return nil, exceptions.NotFoundError("Category", categoryID)
	}

	// Получаем курсы
	data, err := s.courseRepo.GetByCategory(ctx, categoryID)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get courses: %v", err))
	}

	courses := make([]models.CourseResponse, 0, len(data))
	for _, item := range data {
		course := models.CourseResponse{
			Course: models.Course{
				BaseModel: models.BaseModel{
					ID:        fmt.Sprintf("%v", item["id"]),
					CreatedAt: parseTime(item["created_at"]),
					UpdatedAt: parseTime(item["updated_at"]),
				},
				Title:       fmt.Sprintf("%v", item["title"]),
				Description: fmt.Sprintf("%v", item["description"]),
				Level:       fmt.Sprintf("%v", item["level"]),
				CategoryID:  fmt.Sprintf("%v", item["category_id"]),
				Visibility:  fmt.Sprintf("%v", item["visibility"]),
			},
		}
		courses = append(courses, course)
	}

	return courses, nil
}