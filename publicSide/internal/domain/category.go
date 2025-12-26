// Package domain определяет основные бизнес-сущности и модели данных,
// которые используются во всем приложении.
package domain

import "time"

// Category представляет собой категорию курсов.
type Category struct {
	ID        string    `json:"id"`        // Уникальный идентификатор
	Title     string    `json:"title"`     // Название категории
	CreatedAt time.Time `json:"created_at"`// Время создания
	UpdatedAt time.Time `json:"updated_at"`// Время последнего обновления
}
