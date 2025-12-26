// Package response содержит структуры данных для формирования HTTP-ответов.
package response

import (
	"time"
)

// LessonDTODetailed - это объект передачи данных (DTO) для урока (детальная версия).
// Используется для отправки полной информации об уроке, включая его содержимое.
type LessonDTODetailed struct {
	ID        string    `json:"id"`         // Уникальный идентификатор урока.
	Title     string    `json:"title"`      // Название урока.
	CourseID  string    `json:"course_id"`  // ID курса, к которому относится урок.
	Content   string    `json:"content"`    // Содержимое урока (HTML/Markdown).
	CreatedAt time.Time `json:"created_at"` // Время создания.
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления.
}
