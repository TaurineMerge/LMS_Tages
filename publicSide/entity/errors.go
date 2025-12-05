// Добавьте в entity/course.go или создайте entity/errors.go:

package entity

import "errors"

// Common errors
var (
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidSlug        = errors.New("invalid slug")
	ErrCourseNotFound     = errors.New("course not found")
	ErrCourseNotPublished = errors.New("course not published")
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrLessonNotPublished = errors.New("lesson not published")
	ErrNoNextLesson       = errors.New("no next lesson available")
	ErrNoPrevLesson       = errors.New("no previous lesson available")
)
