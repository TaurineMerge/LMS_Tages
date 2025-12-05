package repository

import (
	"context"

	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"gorm.io/gorm"
)

type LessonRepository interface {
	FindAllPublished(ctx context.Context, courseID int64) ([]entity.Lesson, error)
	FindByID(ctx context.Context, id int64) (*entity.Lesson, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Lesson, error)
	FindPublishedBySlug(ctx context.Context, slug string) (*entity.Lesson, error)
	FindWithCourse(ctx context.Context, id int64) (*entity.Lesson, error)
}

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) FindAllPublished(ctx context.Context, courseID int64) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	query := r.db.WithContext(ctx).
		Where("is_published = ?", true).
		Order("order_index ASC")

	if courseID > 0 {
		query = query.Where("course_id = ?", courseID)
	}

	err := query.Find(&lessons).Error
	return lessons, err
}

func (r *lessonRepository) FindByID(ctx context.Context, id int64) (*entity.Lesson, error) {
	var lesson entity.Lesson
	err := r.db.WithContext(ctx).First(&lesson, id).Error
	return &lesson, err
}

func (r *lessonRepository) FindBySlug(ctx context.Context, slug string) (*entity.Lesson, error) {
	var lesson entity.Lesson
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&lesson).Error
	return &lesson, err
}

func (r *lessonRepository) FindPublishedBySlug(ctx context.Context, slug string) (*entity.Lesson, error) {
	var lesson entity.Lesson
	err := r.db.WithContext(ctx).
		Where("slug = ? AND is_published = ?", slug, true).
		First(&lesson).Error
	return &lesson, err
}

func (r *lessonRepository) FindWithCourse(ctx context.Context, id int64) (*entity.Lesson, error) {
	var lesson entity.Lesson
	err := r.db.WithContext(ctx).
		Preload("Course", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = ?", true)
		}).
		Where("is_published = ?", true).
		First(&lesson, id).Error
	return &lesson, err
}
