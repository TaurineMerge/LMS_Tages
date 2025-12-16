package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"adminPanel/exceptions"
	"adminPanel/models"
	"adminPanel/repositories"
)

// CategoryService - сервис для работы с категориями
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService создает сервис категорий
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// GetCategories - получение категорий с пагинацией
func (s *CategoryService) GetCategories(ctx context.Context, page, limit int) ([]models.Category, int, error) {
	offset := (page - 1) * limit

	data, err := s.categoryRepo.GetAllWithPagination(ctx, limit, offset)
	if err != nil {
		return nil, 0, exceptions.InternalError(fmt.Sprintf("Failed to get categories: %v", err))
	}

	total, err := s.categoryRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, exceptions.InternalError(fmt.Sprintf("Failed to count categories: %v", err))
	}

	categories := make([]models.Category, 0, len(data))
	for _, item := range data {
		category := models.Category{
			ID:        parseString(item["id"]),
			Title:     parseString(item["title"]),
			CreatedAt: parseTime(item["created_at"]),
			UpdatedAt: parseTime(item["updated_at"]),
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// GetCategory - получение категории по ID
func (s *CategoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	data, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to get category: %v", err))
	}

	if data == nil {
		return nil, exceptions.NotFoundError("Category", id)
	}

	category := &models.Category{
		ID:        parseString(data["id"]),
		Title:     parseString(data["title"]),
		CreatedAt: parseTime(data["created_at"]),
		UpdatedAt: parseTime(data["updated_at"]),
	}

	return category, nil
}

// CreateCategory - создание категории
func (s *CategoryService) CreateCategory(ctx context.Context, input models.CategoryCreate) (*models.Category, error) {
	// Проверяем, существует ли категория с таким названием
	existing, err := s.categoryRepo.GetByTitle(ctx, input.Title)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check existing category: %v", err))
	}

	if existing != nil {
		return nil, exceptions.ConflictError(fmt.Sprintf("Category with title '%s' already exists", input.Title))
	}

	// Создаем категорию
	data, err := s.categoryRepo.Create(ctx, input.Title)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, exceptions.ConflictError("Category with this title already exists")
		}
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to create category: %v", err))
	}

	category := &models.Category{
		ID:        parseString(data["id"]),
		Title:     parseString(data["title"]),
		CreatedAt: parseTime(data["created_at"]),
		UpdatedAt: parseTime(data["updated_at"]),
	}

	return category, nil
}

// UpdateCategory - обновление категории
func (s *CategoryService) UpdateCategory(ctx context.Context, id string, input models.CategoryUpdate) (*models.Category, error) {
	// Проверяем существование категории
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return nil, exceptions.NotFoundError("Category", id)
	}

	// Если title не предоставлен, используем существующий
	newTitle := input.Title
	if newTitle == "" {
		newTitle = parseString(existing["title"])
	}

	// Проверяем, не занято ли новое название другой категорией
	if newTitle != parseString(existing["title"]) {
		categoryWithTitle, err := s.categoryRepo.GetByTitle(ctx, newTitle)
		if err != nil {
			return nil, exceptions.InternalError(fmt.Sprintf("Failed to check title: %v", err))
		}
		if categoryWithTitle != nil && parseString(categoryWithTitle["id"]) != id {
			return nil, exceptions.ConflictError(fmt.Sprintf("Category with title '%s' already exists", newTitle))
		}
	}

	// Обновляем категорию
	data, err := s.categoryRepo.Update(ctx, id, newTitle)
	if err != nil {
		return nil, exceptions.InternalError(fmt.Sprintf("Failed to update category: %v", err))
	}

	category := &models.Category{
		ID:        parseString(data["id"]),
		Title:     parseString(data["title"]),
		CreatedAt: parseTime(data["created_at"]),
		UpdatedAt: parseTime(data["updated_at"]),
	}

	return category, nil
}

// DeleteCategory - удаление категории
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	// Проверяем существование категории
	existing, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return exceptions.InternalError(fmt.Sprintf("Failed to check category: %v", err))
	}

	if existing == nil {
		return exceptions.NotFoundError("Category", id)
	}

	// Проверяем, нет ли связанных курсов
	courseCount, err := s.categoryRepo.CountCoursesForCategory(ctx, id)
	if err != nil {
		return exceptions.InternalError(fmt.Sprintf("Failed to check associated courses: %v", err))
	}

	if courseCount > 0 {
		return exceptions.ConflictError("Category has related courses")
	}

	// Удаляем категорию
	deleted, err := s.categoryRepo.Delete(ctx, id)
	if err != nil {
		return exceptions.InternalError(fmt.Sprintf("Failed to delete category: %v", err))
	}

	if !deleted {
		return exceptions.InternalError("Failed to delete category")
	}

	return nil
}

// Вспомогательная функция для парсинга времени
func parseTime(value interface{}) time.Time {
	if value == nil {
		return time.Time{}
	}

	switch v := value.(type) {
	case time.Time:
		return v
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return t
		}
		return time.Time{}
	default:
		return time.Time{}
	}
}

// Вспомогательная функция для парсинга строк
func parseString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}
