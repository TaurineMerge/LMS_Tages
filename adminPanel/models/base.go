package models

import (
	"time"  // Оставляем здесь, если BaseModel использует time
)

// BaseModel - базовая модель с общими полями
type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginatedResponse - пагинированный ответ
type PaginatedResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}