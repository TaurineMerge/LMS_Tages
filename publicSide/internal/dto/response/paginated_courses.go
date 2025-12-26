// Package response содержит структуры данных для формирования HTTP-ответов.
package response

// PaginatedCoursesData представляет собой структуру данных для ответа,
// содержащего список курсов с информацией о пагинации.
type PaginatedCoursesData struct {
	Items      []CourseDTO `json:"items"`      // Срез (список) курсов на текущей странице.
	Pagination Pagination  `json:"pagination"` // Информация о пагинации.
}
