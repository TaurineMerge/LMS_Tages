package dto

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// CourseDTO represents a course data transfer object.
type CourseDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- Paginated Response for Courses ---

type PaginatedCoursesData struct {
	Items      []CourseDTO       `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}
