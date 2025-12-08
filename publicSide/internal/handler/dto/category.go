package dto

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// CategoryDTO represents a category data transfer object.
type CategoryDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Paginated Response for Categories ---

type PaginatedCategoriesData struct {
	Items      []CategoryDTO       `json:"items"`
	Pagination domain.Pagination `json:"pagination"`
}
