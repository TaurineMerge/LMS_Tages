// Package response содержит структуры данных для формирования HTTP-ответов.
package response

import "time"

// CategoryDTO - это объект передачи данных (DTO) для категории.
// Используется для отправки информации о категории клиенту.
type CategoryDTO struct {
	ID        string    `json:"id"`         // Уникальный идентификатор категории.
	Title     string    `json:"title"`      // Название категории.
	CreatedAt time.Time `json:"created_at"` // Время создания.
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления.
}
