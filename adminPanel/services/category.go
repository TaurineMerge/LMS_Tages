// Package services implements the business logic for the admin panel.
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"adminPanel/exceptions"
	"adminPanel/models"
	"adminPanel/repositories"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CategoryService - сервис для работы с категориями
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

var categoryTracer = otel.Tracer("admin-panel/category-service")

// NewCategoryService создает сервис категорий
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetCategories - получение всех категорий
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

// GetCategory - получение категории по ID
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

// CreateCategory - создание категории
func (s *CategoryService) CreateCategory(ctx context.Context, input models.CategoryCreate) (*models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.CreateCategory")
	span.SetAttributes(attribute.String("category.title", input.Title))
	defer span.End()

	// Проверяем, существует ли категория с таким названием
	existing, err := s.categoryRepo.GetByTitle(ctx, input.Title)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check existing category: %v", err))
	}

	if existing != nil {
		return nil, exceptions.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
	}

	// Создаем категорию
	data, err := s.categoryRepo.Create(ctx, input.Title)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, exceptions.ConflictError("Category with this title already exists")
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

// UpdateCategory - обновление категории
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, input models.CategoryUpdate) (*models.Category, error) {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.UpdateCategory")
	span.SetAttributes(
		attribute.String("category.id", id),
		attribute.String("category.title", input.Title),
	)
	defer span.End()

	// Проверяем существование категории
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Category", id)
	}

	// Проверяем, не занято ли новое название
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

	// Обновляем категорию
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

// DeleteCategory - удаление категории
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	ctx, span := categoryTracer.Start(ctx, "CategoryService.DeleteCategory")
	span.SetAttributes(attribute.String("category.id", id))
	defer span.End()

	// Проверяем существование категории
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return exceptions.NotFoundError("Category", id)
	}

	// Проверяем, нет ли связанных курсов
	courseCount, err := s.categoryRepo.CountCoursesForCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return exceptions.InternalError(fmt.Sprintf("Failed to check associated courses: %v", err))
	}

	if courseCount > 0 {
		return exceptions.ConflictError("Cannot delete category with associated courses")
	}

	// Удаляем категорию
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

// Вспомогательная функция для парсинга времени
func parseTime(value interface{}) time.Time {
	if str, ok := value.(string); ok {
		if t, err := time.Parse(time.RFC3339, str); err == nil {
			return t
		}
	}
	if t, ok := value.(time.Time); ok {
		return t
	}
	return time.Time{}
}
