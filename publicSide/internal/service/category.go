// Package service предоставляет бизнес-логику приложения, работая как промежуточный
// слой между обработчиками (handlers) и репозиториями (repositories).
package service

import (
	"context"
	"math"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CategoryService определяет интерфейс для бизнес-логики, связанной с категориями.
type CategoryService interface {
	// GetAll получает все категории с пагинацией.
	GetAll(ctx context.Context, page, limit int) ([]response.CategoryDTO, response.Pagination, error)
	// GetAllNotEmpty получает все категории, в которых есть хотя бы один публичный курс.
	GetAllNotEmpty(ctx context.Context, page, limit int) ([]response.CategoryDTO, response.Pagination, error)
	// GetByID получает категорию по ее ID.
	GetByID(ctx context.Context, categoryID string) (response.CategoryDTO, error)
}

// categoryService является реализацией CategoryService.
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService создает новый экземпляр categoryService.
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

// toCategoryDTO преобразует доменную модель Category в DTO CategoryDTO.
func toCategoryDTO(category domain.Category) response.CategoryDTO {
	return response.CategoryDTO{
		ID:        category.ID,
		Title:     category.Title,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

// GetAll обрабатывает запрос на получение всех категорий, валидирует параметры
// пагинации, вызывает репозиторий и преобразует результат в DTO.
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

// GetAllNotEmpty обрабатывает запрос на получение непустых категорий, валидирует
// параметры пагинации, вызывает репозиторий и преобразует результат в DTO.
func (s *categoryService) GetAllNotEmpty(ctx context.Context, page, limit int) ([]response.CategoryDTO, response.Pagination, error) {
	ctx, span := otel.Tracer("categoryService").Start(ctx, "GetAllNotEmpty")
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

	categories, total, err := s.repo.GetAllNotEmpty(ctx, page, limit)
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

// GetByID находит категорию по ID. Если категория не найдена,
// возвращает стандартизированную ошибку `apperrors.NewNotFound`.
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
