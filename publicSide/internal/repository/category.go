package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/google/uuid"
)

type CategoryRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	GetByID(ctx context.Context, id string) (domain.Category, error)
}

type categoryMemoryRepository struct {
	categories []domain.Category
}

func NewCategoryMemoryRepository() CategoryRepository {
	categories := []domain.Category{
		{ID: "550e8400-e29b-41d4-a716-446655440000", Title: "Программирование", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "550e8400-e29b-41d4-a716-446655440001", Title: "Дизайн", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Title: "Маркетинг", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	return &categoryMemoryRepository{categories: categories}
}

func (r *categoryMemoryRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	total := len(r.categories)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []domain.Category{}, total, nil
	}
	if end > total {
		end = total
	}

	return r.categories[start:end], total, nil
}

func (r *categoryMemoryRepository) GetByID(ctx context.Context, id string) (domain.Category, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return domain.Category{}, fmt.Errorf("invalid uuid: %w", err)
	}

	for _, category := range r.categories {
		if category.ID == id {
			return category, nil
		}
	}
	return domain.Category{}, fmt.Errorf("category with id %s not found", id)
}
