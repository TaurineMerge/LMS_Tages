package service

import (
	"context"
	"math"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
)

type CourseService interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Course, domain.Pagination, error)
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, domain.Pagination, error)
	GetByID(ctx context.Context, id string) (domain.Course, error)
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) GetAll(ctx context.Context, page, limit int) ([]domain.Course, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courses, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return courses, pagination, nil
}

func (s *courseService) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courses, total, err := s.repo.GetAllByCategoryID(ctx, categoryID, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return courses, pagination, nil
}

func (s *courseService) GetByID(ctx context.Context, id string) (domain.Course, error) {
	return s.repo.GetByID(ctx, id)
}
