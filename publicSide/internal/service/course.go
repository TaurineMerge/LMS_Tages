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

type CourseService interface {
	GetAll(ctx context.Context, page, limit int) ([]dto.CourseDTO, domain.Pagination, error)
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]dto.CourseDTO, domain.Pagination, error)
	GetByID(ctx context.Context, id string) (dto.CourseDTO, error)
}

type courseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func toCourseDTO(course domain.Course) dto.CourseDTO {
	return dto.CourseDTO{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Level:       course.Level,
		CategoryID:  course.CategoryID,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}

func (s *courseService) GetAll(ctx context.Context, page, limit int) ([]dto.CourseDTO, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courses, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, apperrors.NewInternal()
	}

	courseDTOs := make([]dto.CourseDTO, len(courses))
	for i, course := range courses {
		courseDTOs[i] = toCourseDTO(course)
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return courseDTOs, pagination, nil
}

func (s *courseService) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]dto.CourseDTO, domain.Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	courses, total, err := s.repo.GetAllByCategoryID(ctx, categoryID, page, limit)
	if err != nil {
		return nil, domain.Pagination{}, apperrors.NewInternal()
	}

	courseDTOs := make([]dto.CourseDTO, len(courses))
	for i, course := range courses {
		courseDTOs[i] = toCourseDTO(course)
	}

	pagination := domain.Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
		Pages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return courseDTOs, pagination, nil
}

func (s *courseService) GetByID(ctx context.Context, id string) (dto.CourseDTO, error) {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return dto.CourseDTO{}, apperrors.NewNotFound("Course")
		}
		return dto.CourseDTO{}, apperrors.NewInternal()
	}
	return toCourseDTO(course), nil
}
