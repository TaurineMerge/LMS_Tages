package service

import (
	"context"
	"math"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
)

type CategoryService interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, domain.Pagination, error)
	GetByID(ctx context.Context, id string) (domain.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context, page, limit int) ([]domain.Category, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	categories, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return categories, pagination, nil
}

func (s *categoryService) GetByID(ctx context.Context, id string) (domain.Category, error) {
	return s.repo.GetByID(ctx, id)
}
