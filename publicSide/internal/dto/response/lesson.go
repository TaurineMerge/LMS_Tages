// Package response содержит структуры данных для формирования HTTP-ответов.
package response

import "time"

// LessonDTO - это объект передачи данных (DTO) для урока (краткая версия).
// Используется для отправки информации об уроке без его содержимого, например, в списках.
type LessonDTO struct {
	ID        string    `json:"id"`         // Уникальный идентификатор урока.
	Title     string    `json:"title"`      // Название урока.
	CourseID  string    `json:"course_id"`  // ID курса, к которому относится урок.
	CreatedAt time.Time `json:"created_at"` // Время создания.
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления.
}
