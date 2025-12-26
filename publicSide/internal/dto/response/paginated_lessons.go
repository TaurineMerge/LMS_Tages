// Package response содержит структуры данных для формирования HTTP-ответов.
package response

// PaginatedLessonsData представляет собой структуру данных для ответа,
// содержащего список уроков с информацией о пагинации.
type PaginatedLessonsData struct {
	Items      []LessonDTO `json:"items"`      // Срез (список) уроков на текущей странице.
	Pagination Pagination  `json:"pagination"` // Информация о пагинации.
}
