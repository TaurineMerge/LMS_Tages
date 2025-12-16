// Package services реализует бизнес-логику для admin panel
//
// Пакет содержит сервисы для работы с сущностями:
//   - CategoryService: работа с категориями
//   - CourseService: работа с курсами
//   - LessonService: работа с уроками
//
// Каждый сервис отвечает за:
//   - Валидацию данных
//   - Бизнес-логику операций
//   - Интеграцию с репозиториями
//   - Трассировку операций
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

// CategoryService - сервис для работы с категориями курсов
//
// Сервис предоставляет методы для управления категориями:
//   - Получение списка категорий
//   - Получение категории по ID
//   - Создание новой категории
//   - Обновление категории
//   - Удаление категории
//
// Особенности:
//   - Проверка уникальности названий
//   - Контроль связанных курсов при удалении
//   - Интеграция с OpenTelemetry для трассировки
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

var categoryTracer = otel.Tracer("admin-panel/category-service")

// NewCategoryService создает новый сервис для работы с категориями
//
// Параметры:
//   - categoryRepo: репозиторий для работы с категориями
//
// Возвращает:
//   - *CategoryService: указатель на новый сервис
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetCategories получает список всех категорий
//
// Метод возвращает все доступные категории, отсортированные по названию.
//
// Параметры:
//   - ctx: контекст выполнения
//
// Возвращает:
//   - []models.Category: список категорий
//   - error: ошибка выполнения (если есть)
func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.GetCategories")
	defer span.End()

	data, err := s.categoryRepo.GetAll(ctx, 100, 0, "title", "ASC")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get categories: %v", err))
	}

	categories := make([]models.Category, 0, len(data))
	for _, item := range data {
		category := models.Category{
			BaseModel: models.BaseModel{
				ID:        toString(item["id"]),
				CreatedAt: parseTime(item["created_at"]),
				UpdatedAt: parseTime(item["updated_at"]),
			},
			Title: toString(item["title"]),
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategory получает категорию по уникальному идентификатору
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор категории
//
// Возвращает:
//   - *models.Category: указатель на категорию
//   - error: ошибка выполнения (если есть)
func (s *CategoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.GetCategory")
	span.SetAttributes(attribute.String("category.id", id))
	defer span.End()

	data, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get category: %v", err))
	}

	if data == nil {
		return nil, exceptions.NotFoundError("Category", id)
	}

	category := &models.Category{
		BaseModel: models.BaseModel{
			ID:        toString(data["id"]),
			CreatedAt: parseTime(data["created_at"]),
			UpdatedAt: parseTime(data["updated_at"]),
		},
		Title: toString(data["title"]),
	}

	return category, nil
}

// CreateCategory создает новую категорию
//
// Перед созданием проверяет уникальность названия категории.
//
// Параметры:
//   - ctx: контекст выполнения
//   - input: данные для создания категории
//
// Возвращает:
//   - *models.Category: созданная категория
//   - error: ошибка выполнения (если есть)
func (s *CategoryService) CreateCategory(ctx context.Context, input models.CategoryCreate) (*models.Category, error) {
	existing, err := s.categoryRepo.GetByTitle(ctx, input.Title)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check existing category: %v", err))
	}

	if existing != nil {
		return nil, exceptions.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
	}

	data, err := s.categoryRepo.Create(ctx, input.Title)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, exceptions.ConflictError("Category with this title already exists")
		}
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create category: %v", err))
	}

	category := &models.Category{
		BaseModel: models.BaseModel{
			ID:        toString(data["id"]),
			CreatedAt: parseTime(data["created_at"]),
			UpdatedAt: parseTime(data["updated_at"]),
		},
		Title: toString(data["title"]),
	}

	return category, nil
}

// UpdateCategory обновляет существующую категорию
//
// Проверяет существование категории и уникальность нового названия.
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор категории
//   - input: данные для обновления категории
//
// Возвращает:
//   - *models.Category: обновленная категория
//   - error: ошибка выполнения (если есть)
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, input models.CategoryUpdate) (*models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.UpdateCategory")
	span.SetAttributes(
		attribute.String("category.id", id),
		attribute.String("category.title", input.Title),
	)
	defer span.End()

	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Category", id)
	}

	if input.Title != fmt.Sprintf("%v", existing["title"]) {
		categoryWithTitle, err := s.categoryRepo.GetByTitle(ctx, input.Title)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, exceptions.InternalError(fmt.Sprintf("Failed to check title: %v", err))
		}
		if categoryWithTitle != nil {
			return nil, exceptions.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
		}
	}

	data, err := s.categoryRepo.Update(ctx, id, input.Title)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to update category: %v", err))
	}

	category := &models.Category{
		BaseModel: models.BaseModel{
			ID:        toString(data["id"]),
			CreatedAt: parseTime(data["created_at"]),
			UpdatedAt: parseTime(data["updated_at"]),
		},
		Title: toString(data["title"]),
	}

	return category, nil
}

// DeleteCategory удаляет категорию по уникальному идентификатору
//
// Перед удалением проверяет наличие связанных курсов.
// Категория с курсами не может быть удалена.
//
// Параметры:
//   - ctx: контекст выполнения
//   - id: уникальный идентификатор категории
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.DeleteCategory")
	span.SetAttributes(attribute.String("category.id", id))
	defer span.End()

	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return exceptions.NotFoundError("Category", id)
	}

	courseCount, err := s.categoryRepo.CountCoursesForCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check associated courses: %v", err))
	}

	if courseCount > 0 {
		return exceptions.ConflictError("Cannot delete category with associated courses")
	}

	deleted, err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to delete category: %v", err))
	}

	if !deleted {
		return exceptions.InternalError("Failed to delete category")
	}

	return nil
}
