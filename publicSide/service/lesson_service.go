package service

import (
	"context"

	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"github.com/TaurineMerge/LMS_Tages/publicSide/repository"
)

// LessonService defines the interface for lesson business logic
type LessonService interface {
	GetAllPublished(ctx context.Context, courseID int64) ([]entity.Lesson, error)
	GetByID(ctx context.Context, id int64) (*entity.Lesson, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Lesson, error)
	GetWithCourse(ctx context.Context, id int64) (*entity.Lesson, error)
	GetPublishedBySlug(ctx context.Context, slug string) (*entity.Lesson, error)
	GetCourseLessons(ctx context.Context, courseID int64) ([]entity.Lesson, error)
	GetNextLesson(ctx context.Context, currentLessonID int64) (*entity.Lesson, error)
	GetPrevLesson(ctx context.Context, currentLessonID int64) (*entity.Lesson, error)
}

// lessonService implements LessonService
type lessonService struct {
	lessonRepo repository.LessonRepository
	courseRepo repository.CourseRepository
}

// NewLessonService creates a new lesson service
func NewLessonService(lessonRepo repository.LessonRepository, courseRepo repository.CourseRepository) LessonService {
	return &lessonService{
		lessonRepo: lessonRepo,
		courseRepo: courseRepo,
	}
}

// GetAllPublished returns all published lessons
func (s *lessonService) GetAllPublished(ctx context.Context, courseID int64) ([]entity.Lesson, error) {
	// Если указан courseID, проверяем что курс существует и опубликован
	if courseID > 0 {
		course, err := s.courseRepo.FindByID(ctx, courseID)
		if err != nil {
			return nil, err
		}
		if !course.IsPublished {
			return nil, entity.ErrCourseNotPublished
		}
	}

	lessons, err := s.lessonRepo.FindAllPublished(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// Здесь можно добавить бизнес-логику:
	// - Добавление прогресса пользователя
	// - Фильтрация по доступности
	// - Добавление метаданных

	return lessons, nil
}

// GetByID returns lesson by ID
func (s *lessonService) GetByID(ctx context.Context, id int64) (*entity.Lesson, error) {
	if id <= 0 {
		return nil, entity.ErrInvalidID
	}

	lesson, err := s.lessonRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли урок
	if !lesson.IsPublished {
		return nil, entity.ErrLessonNotPublished
	}

	// Проверяем, опубликован ли родительский курс
	course, err := s.courseRepo.FindByID(ctx, lesson.CourseID)
	if err != nil {
		return nil, err
	}
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return lesson, nil
}

// GetBySlug returns lesson by slug
func (s *lessonService) GetBySlug(ctx context.Context, slug string) (*entity.Lesson, error) {
	if slug == "" {
		return nil, entity.ErrInvalidSlug
	}

	lesson, err := s.lessonRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли урок
	if !lesson.IsPublished {
		return nil, entity.ErrLessonNotPublished
	}

	// Проверяем, опубликован ли родительский курс
	course, err := s.courseRepo.FindByID(ctx, lesson.CourseID)
	if err != nil {
		return nil, err
	}
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return lesson, nil
}

// GetWithCourse returns lesson with course info
func (s *lessonService) GetWithCourse(ctx context.Context, id int64) (*entity.Lesson, error) {
	if id <= 0 {
		return nil, entity.ErrInvalidID
	}

	lesson, err := s.lessonRepo.FindWithCourse(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли урок
	if !lesson.IsPublished {
		return nil, entity.ErrLessonNotPublished
	}

	// Проверяем, опубликован ли родительский курс
	if lesson.Course != nil && !lesson.Course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return lesson, nil
}

// GetPublishedBySlug returns published lesson by slug with course
func (s *lessonService) GetPublishedBySlug(ctx context.Context, slug string) (*entity.Lesson, error) {
	if slug == "" {
		return nil, entity.ErrInvalidSlug
	}

	lesson, err := s.lessonRepo.FindPublishedBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Уже проверено в репозитории, что урок опубликован
	// Дополнительно проверяем курс
	if lesson.Course != nil && !lesson.Course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	return lesson, nil
}

// GetCourseLessons returns all published lessons for a course
func (s *lessonService) GetCourseLessons(ctx context.Context, courseID int64) ([]entity.Lesson, error) {
	if courseID <= 0 {
		return nil, entity.ErrInvalidID
	}

	// Проверяем что курс существует и опубликован
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if !course.IsPublished {
		return nil, entity.ErrCourseNotPublished
	}

	lessons, err := s.lessonRepo.FindAllPublished(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return lessons, nil
}

// GetNextLesson returns the next lesson in the course
func (s *lessonService) GetNextLesson(ctx context.Context, currentLessonID int64) (*entity.Lesson, error) {
	if currentLessonID <= 0 {
		return nil, entity.ErrInvalidID
	}

	// Получаем текущий урок
	currentLesson, err := s.lessonRepo.FindByID(ctx, currentLessonID)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли текущий урок
	if !currentLesson.IsPublished {
		return nil, entity.ErrLessonNotPublished
	}

	// Получаем все уроки курса
	lessons, err := s.lessonRepo.FindAllPublished(ctx, currentLesson.CourseID)
	if err != nil {
		return nil, err
	}

	// Ищем следующий урок по order_index
	for i, lesson := range lessons {
		if lesson.ID == currentLessonID && i+1 < len(lessons) {
			nextLesson := lessons[i+1]
			return &nextLesson, nil
		}
	}

	// Следующего урока нет
	return nil, entity.ErrNoNextLesson
}

// GetPrevLesson returns the previous lesson in the course
func (s *lessonService) GetPrevLesson(ctx context.Context, currentLessonID int64) (*entity.Lesson, error) {
	if currentLessonID <= 0 {
		return nil, entity.ErrInvalidID
	}

	// Получаем текущий урок
	currentLesson, err := s.lessonRepo.FindByID(ctx, currentLessonID)
	if err != nil {
		return nil, err
	}

	// Проверяем, опубликован ли текущий урок
	if !currentLesson.IsPublished {
		return nil, entity.ErrLessonNotPublished
	}

	// Получаем все уроки курса
	lessons, err := s.lessonRepo.FindAllPublished(ctx, currentLesson.CourseID)
	if err != nil {
		return nil, err
	}

	// Ищем предыдущий урок по order_index
	for i, lesson := range lessons {
		if lesson.ID == currentLessonID && i-1 >= 0 {
			prevLesson := lessons[i-1]
			return &prevLesson, nil
		}
	}

	// Предыдущего урока нет
	return nil, entity.ErrNoPrevLesson
}
