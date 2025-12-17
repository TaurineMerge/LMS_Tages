// Package service contains the business logic of the application.
// It orchestrates data from repositories and prepares it for the handler layer.
package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CategoryService defines the interface for category-related business logic.
type CategoryService interface {
	// GetAll retrieves a paginated list of all categories.
	GetAll(ctx context.Context, page, limit int) ([]response.CategoryDTO, response.Pagination, error)
	// GetByID retrieves a single category by its ID.
	GetByID(ctx context.Context, categoryID string) (response.CategoryDTO, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new instance of a category service.
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func toCategoryDTO(category domain.Category) response.CategoryDTO {
	return response.CategoryDTO{
		ID:        category.ID,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (s *categoryService) GetAll(ctx context.Context, page, limit int) ([]response.CategoryDTO, response.Pagination, error) {
	ctx, span := otel.Tracer("categoryService").Start(ctx, "GetAll")
	defer span.End()

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	categories, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, response.Pagination{}, err
	}

	categoryDTOs := make([]response.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = toCategoryDTO(category)
	}

	pagination := response.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return categoryDTOs, pagination, nil
}

func (s *categoryService) GetByID(ctx context.Context, categoryID string) (response.CategoryDTO, error) {
	ctx, span := otel.Tracer("categoryService").Start(ctx, "GetByID")
	span.SetAttributes(attribute.String("category.id", categoryID))
	defer span.End()

	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.CategoryDTO{}, apperrors.NewNotFound("Category")
		}
		return response.CategoryDTO{}, err
	}
	return toCategoryDTO(category), nil
}
