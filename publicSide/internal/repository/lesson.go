package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/google/uuid"
)

type LessonRepository interface {
	GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]domain.Lesson, int, error)
	GetByID(ctx context.Context, id string) (domain.Lesson, error)
}

type lessonMemoryRepository struct {
	lessons []domain.Lesson
}

func NewLessonMemoryRepository() LessonRepository {
	lessons := []domain.Lesson{
		{ID: "770e8400-e29b-41d4-a716-446655440001", Title: "Урок 1: Введение", CourseID: "660e8400-e29b-41d4-a716-446655440001", Content: []domain.ContentBlock{{ContentType: "text", Data: "Это первый урок."}}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "770e8400-e29b-41d4-a716-446655440002", Title: "Урок 2: Переменные", CourseID: "660e8400-e29b-41d4-a716-446655440001", Content: []domain.ContentBlock{{ContentType: "text", Data: "Это второй урок."}}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "770e8400-e29b-41d4-a716-446655440003", Title: "Урок 1: HTTP", CourseID: "660e8400-e29b-41d4-a716-446655440002", Content: []domain.ContentBlock{{ContentType: "text", Data: "Это первый урок по Go."}}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	return &lessonMemoryRepository{lessons: lessons}
}

func (r *lessonMemoryRepository) GetAllByCourseID(ctx context.Context, courseID string, page, limit int) ([]domain.Lesson, int, error) {
	_, err := uuid.Parse(courseID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid course id: %w", err)
	}

	var filtered []domain.Lesson
	for _, lesson := range r.lessons {
		if lesson.CourseID == courseID {
			filtered = append(filtered, lesson)
		}
	}

	total := len(filtered)
	start := (page - 1) * limit
	end := start + limit

	if start > total {
		return []domain.Lesson{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (r *lessonMemoryRepository) GetByID(ctx context.Context, id string) (domain.Lesson, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("invalid uuid: %w", err)
	}
	for _, lesson := range r.lessons {
		if lesson.ID == id {
			return lesson, nil
		}
	}
	return domain.Lesson{}, fmt.Errorf("lesson with id %s not found", id)
}
