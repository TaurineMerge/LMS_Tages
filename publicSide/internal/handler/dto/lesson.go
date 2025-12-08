package dto

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// --- Lesson DTOs ---

// LessonDTO is for lists, without content.
type LessonDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CourseID  string    `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LessonDTODetailed is for the detailed view, including content.
type LessonDTODetailed struct {
	ID        string                `json:"id"`
	Title     string                `json:"title"`
	CourseID  string                `json:"course_id"`
	Content   []domain.ContentBlock `json:"content"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// --- Paginated Response for Lessons ---

type PaginatedLessonsData struct {
	Items      []LessonDTO       `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}
