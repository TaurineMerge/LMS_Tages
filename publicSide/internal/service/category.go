package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
)

type CategoryService interface {
	GetAll(ctx context.Context, page, limit int) ([]dto.CategoryDTO, domain.Pagination, error)
	GetByID(ctx context.Context, id string) (dto.CategoryDTO, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func toCategoryDTO(category domain.Category) dto.CategoryDTO {
	return dto.CategoryDTO{
		ID:        category.ID,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func (s *categoryService) GetAll(ctx context.Context, page, limit int) ([]dto.CategoryDTO, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	categories, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, apperrors.NewInternal()
	}

	categoryDTOs := make([]dto.CategoryDTO, len(categories))
	for i, category := range categories {
		categoryDTOs[i] = toCategoryDTO(category)
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return categoryDTOs, pagination, nil
}

func (s *categoryService) GetByID(ctx context.Context, id string) (dto.CategoryDTO, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return dto.CategoryDTO{}, apperrors.NewNotFound("Category")
		}
		return dto.CategoryDTO{}, apperrors.NewInternal()
	}
	return toCategoryDTO(category), nil
}
