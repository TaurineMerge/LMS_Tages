package models

import "time"

type Category struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryCreate - DTO для создания категории
type CategoryCreate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryUpdate - DTO для обновления категории
type CategoryUpdate struct {
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
}

// CategoryResponse - ответ с категорией
type CategoryResponse struct {
	Status string   `json:"status"`
	Data   Category `json:"data"`
}

// PaginatedCategoriesResponse - пагинированный список категорий
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Category `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}

// Pagination - структура пагинации
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// StatusOnly - только статус
type StatusOnly struct {
	Status string `json:"status"`
}
