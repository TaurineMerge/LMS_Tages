package dto

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// --- Error Response ---

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status string      `json:"status"`
	Error  ErrorDetail `json:"error"`
}

// --- Success Response ---

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// --- Lesson DTOs ---

// LessonDTO is for lists, without content.
type LessonDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CourseID  string    `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Paginated Responses ---

type PaginatedCategoriesData struct {
	Items      []domain.Category `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

type PaginatedCoursesData struct {
	Items      []domain.Course   `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}

type PaginatedLessonsData struct {
	Items      []LessonDTO       `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}
