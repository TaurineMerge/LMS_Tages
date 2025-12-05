package service

import (
	"context"

	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"github.com/TaurineMerge/LMS_Tages/publicSide/repository"
)

// CourseService defines the interface for course business logic
type CourseService interface {
	GetAllPublished(ctx context.Context) ([]entity.Course, error)
	GetByID(ctx context.Context, id int64) (*entity.Course, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Course, error)
	GetWithLessons(ctx context.Context, id int64) (*entity.Course, error)
	GetPublishedBySlug(ctx context.Context, slug string) (*entity.Course, error)
	SearchCourses(ctx context.Context, query string, limit int) ([]entity.Course, error)
}

// courseService implements CourseService
type courseService struct {
	courseRepo repository.CourseRepository
}

// NewCourseService creates a new course service
func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{
		courseRepo: courseRepo,
	}
}

// GetAllPublished returns all published courses
func (s *courseService) GetAllPublished(ctx context.Context) ([]entity.Course, error) {
	courses, err := s.courseRepo.FindAllPublished(ctx)
	if err != nil {
		return nil, err
	}

	// Здесь можно добавить дополнительную бизнес-логику:
	// - Фильтрация
	// - Сортировка
	// - Ограничение количества
	// - Добавление статистики

	return courses, nil
}

// GetByID returns course by ID
func (s *courseService) GetByID(ctx context.Context, id int64) (*entity.Course, error) {
	if id <= 0 {
		return nil, entity.ErrInvalidID
	}

	course, err := s.courseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли курс
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return course, nil
}

// GetBySlug returns course by slug
func (s *courseService) GetBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	if slug == "" {
		return nil, entity.ErrInvalidSlug
	}

	course, err := s.courseRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли курс
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return course, nil
}

// GetWithLessons returns course with published lessons
func (s *courseService) GetWithLessons(ctx context.Context, id int64) (*entity.Course, error) {
	if id <= 0 {
		return nil, entity.ErrInvalidID
	}

	course, err := s.courseRepo.FindWithLessons(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли курс
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	// Фильтруем уроки (уже сделано в репозитории)
	// Можно добавить дополнительную логику:
	// - Проверка порядка уроков
	// - Добавление метаданных
	// - Кэширование

	return course, nil
}

// GetPublishedBySlug returns published course by slug with lessons
func (s *courseService) GetPublishedBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	if slug == "" {
		return nil, entity.ErrInvalidSlug
	}

	course, err := s.courseRepo.FindPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Уже проверено в репозитории, что курс опубликован
	return course, nil
}

// SearchCourses searches courses by title or description
func (s *courseService) SearchCourses(ctx context.Context, query string, limit int) ([]entity.Course, error) {
	if query == "" {
		return s.GetAllPublished(ctx)
	}

	// В реальном проекте здесь будет вызов репозитория с поиском
	// Пока возвращаем все курсы и фильтруем на уровне сервиса
	courses, err := s.courseRepo.FindAllPublished(ctx)
	if err != nil {
		return nil, err
	}

	// Простой поиск по названию (в реальном проекте используйте полнотекстовый поиск)
	var results []entity.Course
	for _, course := range courses {
		if contains(course.Title, query) || contains(course.Description, query) {
			results = append(results, course)
		}

		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, nil
}

// Helper function for simple search
func contains(s, substr string) bool {
	// В реальном проекте используйте strings.Contains с учетом регистра
	return len(s) >= len(substr) &&
		(len(substr) == 0 ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())
}
