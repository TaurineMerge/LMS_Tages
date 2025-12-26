// Package services предоставляет бизнес-логику для приложения adminPanel.
// Включает сервисы для работы с категориями, курсами, уроками и т.д.
package services

import (
	"context"
	"fmt"
	"strings"

	"adminPanel/handlers/dto/request"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/repositories"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CategoryService предоставляет бизнес-логику для работы с категориями.
// Содержит репозиторий для доступа к данным и методы для CRUD операций.
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// categoryTracer трассировщик для сервиса категорий.
// Используется для отслеживания операций с категориями.
var categoryTracer = otel.Tracer("admin-panel/category-service")

// NewCategoryService создает новый экземпляр CategoryService.
// Принимает репозиторий категорий.
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetCategories получает все категории, отсортированные по заголовку.
// Возвращает список моделей Category.
func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.GetCategories")
	defer span.End()

	data, err := s.categoryRepo.GetAll(ctx, 100, 0, "title", "ASC")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get categories: %v", err))
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

// GetCategory получает категорию по ID.
// Возвращает модель Category или ошибку, если не найдена.
func (s *CategoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.GetCategory")
	span.SetAttributes(attribute.String("category.id", id))
	defer span.End()

	data, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to get category: %v", err))
	}

	if data == nil {
		return nil, middleware.NotFoundError("Category", id)
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

// CreateCategory создает новую категорию на основе данных из request.CategoryCreate.
// Проверяет уникальность заголовка и возвращает созданную категорию.
func (s *CategoryService) CreateCategory(ctx context.Context, input request.CategoryCreate) (*models.Category, error) {
	existing, err := s.categoryRepo.GetByTitle(ctx, input.Title)
	if err != nil {
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check existing category: %v", err))
	}

	if existing != nil {
		return nil, middleware.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
	}

	data, err := s.categoryRepo.Create(ctx, input.Title)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, middleware.ConflictError("Category with this title already exists")
		}
		return nil, middleware.InternalError(fmt.Sprintf("Failed to create category: %v", err))
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

// UpdateCategory обновляет категорию по ID на основе данных из request.CategoryUpdate.
// Проверяет существование и уникальность заголовка, возвращает обновленную категорию.
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, input request.CategoryUpdate) (*models.Category, error) {
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
		return nil, middleware.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return nil, middleware.NotFoundError("Category", id)
	}

	if input.Title != fmt.Sprintf("%v", existing["title"]) {
		categoryWithTitle, err := s.categoryRepo.GetByTitle(ctx, input.Title)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, middleware.InternalError(fmt.Sprintf("Failed to check title: %v", err))
		}
		if categoryWithTitle != nil {
			return nil, middleware.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
		}
	}

	data, err := s.categoryRepo.Update(ctx, id, input.Title)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, middleware.InternalError(fmt.Sprintf("Failed to update category: %v", err))
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

// DeleteCategory удаляет категорию по ID.
// Проверяет существование и отсутствие связанных курсов перед удалением.
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.DeleteCategory")
	span.SetAttributes(attribute.String("category.id", id))
	defer span.End()

	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return middleware.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return middleware.NotFoundError("Category", id)
	}

	courseCount, err := s.categoryRepo.CountCoursesForCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return middleware.InternalError(fmt.Sprintf("Failed to check associated courses: %v", err))
	}

	if courseCount > 0 {
		return middleware.ConflictError("Cannot delete category with associated courses")
	}

	deleted, err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return middleware.InternalError(fmt.Sprintf("Failed to delete category: %v", err))
	}

	if !deleted {
		return middleware.InternalError("Failed to delete category")
	}

	return nil
}
