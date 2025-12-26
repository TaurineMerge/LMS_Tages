// Package response содержит структуры данных для формирования HTTP-ответов.
package response

import "time"

// CourseDTO - это объект передачи данных (DTO) для курса.
// Используется для отправки информации о курсе клиенту.
type CourseDTO struct {
	ID          string    `json:"id"`           // Уникальный идентификатор курса.
	Title       string    `json:"title"`        // Название курса.
	Description string    `json:"description"`  // Описание курса.
	Level       string    `json:"level"`        // Уровень сложности.
	CategoryID  string    `json:"category_id"`  // ID категории, к которой относится курс.
	ImageURL    string    `json:"image_url"`    // URL изображения курса.
	CreatedAt   time.Time `json:"created_at"`   // Время создания.
	UpdatedAt   time.Time `json:"updated_at"`   // Время последнего обновления.
}
