// Package domain определяет основные бизнес-сущности и модели данных,
// которые используются во всем приложении.
package domain

import "time"

// Course представляет собой учебный курс.
type Course struct {
	ID          string    `json:"id"`          // Уникальный идентификатор
	Title       string    `json:"title"`       // Название курса
	Description string    `json:"description"` // Описание курса
	Level       string    `json:"level"`       // Уровень сложности (easy, medium, hard)
	Visibility  string    `json:"visibility"`  // Видимость (draft, public)
	CategoryID  string    `json:"category_id"` // ID категории, к которой относится курс
	ImageKey    string    `json:"image_key"`   // Ключ изображения в S3/MinIO
	CreatedAt   time.Time `json:"created_at"`  // Время создания
	UpdatedAt   time.Time `json:"updated_at"`  // Время последнего обновления
}
