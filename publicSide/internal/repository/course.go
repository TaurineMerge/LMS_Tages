package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/google/uuid"
)

type CourseRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Course, int, error)
	GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error)
	GetByID(ctx context.Context, id string) (domain.Course, error)
}

type courseMemoryRepository struct {
	courses []domain.Course
}

func NewCourseMemoryRepository() CourseRepository {
	courses := []domain.Course{
		{ID: "660e8400-e29b-41d4-a716-446655440001", Title: "Введение в Python", Description: "Базовый курс по программированию на Python для начинающих", Level: "easy", CategoryID: "550e8400-e29b-41d4-a716-446655440000", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "660e8400-e29b-41d4-a716-446655440002", Title: "Веб-разработка на Go", Description: "Создание веб-сервисов с использованием языка Go", Level: "medium", CategoryID: "550e8400-e29b-41d4-a716-446655440000", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "660e8400-e29b-41d4-a716-446655440003", Title: "UI/UX для начинающих", Description: "Основы дизайна пользовательских интерфейсов", Level: "easy", CategoryID: "550e8400-e29b-41d4-a716-446655440001", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	return &courseMemoryRepository{courses: courses}
}

func (r *courseMemoryRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Course, int, error) {
	total := len(r.courses)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []domain.Course{}, total, nil
	}
	if end > total {
		end = total
	}

	return r.courses[start:end], total, nil
}

func (r *courseMemoryRepository) GetAllByCategoryID(ctx context.Context, categoryID string, page, limit int) ([]domain.Course, int, error) {
	_, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid category id: %w", err)
	}

	var filtered []domain.Course
	for _, course := range r.courses {
		if course.CategoryID == categoryID {
			filtered = append(filtered, course)
		}
	}

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []domain.Course{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *courseMemoryRepository) GetByID(ctx context.Context, id string) (domain.Course, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return domain.Course{}, fmt.Errorf("invalid uuid: %w", err)
	}
	for _, course := range r.courses {
		if course.ID == id {
			return course, nil
		}
	}
	return domain.Course{}, fmt.Errorf("course with id %s not found", id)
}
