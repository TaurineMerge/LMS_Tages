// Package domain определяет основные бизнес-сущности и модели данных,
// которые используются во всем приложении.
package domain

import "time"

// Lesson представляет собой урок в рамках курса.
type Lesson struct {
	ID        string    `json:"id"`        // Уникальный идентификатор
	Title     string    `json:"title"`     // Название урока
	CourseID  string    `json:"course_id"` // ID курса, к которому относится урок
	Content   string    `json:"content"`   // Содержимое урока (HTML/Markdown)
	CreatedAt time.Time `json:"created_at"`// Время создания
	UpdatedAt time.Time `json:"updated_at"`// Время последнего обновления
}
