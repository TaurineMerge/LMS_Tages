package repository

import (
	"context"

	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"gorm.io/gorm"
)

type CourseRepository interface {
	FindAllPublished(ctx context.Context) ([]entity.Course, error)
	FindByID(ctx context.Context, id int64) (*entity.Course, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Course, error)
	FindWithLessons(ctx context.Context, id int64) (*entity.Course, error)
	FindPublishedBySlug(ctx context.Context, slug string) (*entity.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) FindAllPublished(ctx context.Context) ([]entity.Course, error) {
	var courses []entity.Course
	err := r.db.WithContext(ctx).
		Where("is_published = ?", true).
		Order("created_at DESC").
		Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByID(ctx context.Context, id int64) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).First(&course, id).Error
	return &course, err
}

func (r *courseRepository) FindBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&course).Error
	return &course, err
}

func (r *courseRepository) FindWithLessons(ctx context.Context, id int64) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).
		Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = ?", true).Order("order_index ASC")
		}).
		Where("is_published = ?", true).
		First(&course, id).Error
	return &course, err
}

func (r *courseRepository) FindPublishedBySlug(ctx context.Context, slug string) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).
		Where("slug = ? AND is_published = ?", slug, true).
		First(&course).Error
	return &course, err
}
