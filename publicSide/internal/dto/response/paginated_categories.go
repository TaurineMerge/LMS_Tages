// Package response содержит структуры данных для формирования HTTP-ответов.
package response

// PaginatedCategoriesData представляет собой структуру данных для ответа,
// содержащего список категорий с информацией о пагинации.
type PaginatedCategoriesData struct {
	Items      []CategoryDTO `json:"items"`      // Срез (список) категорий на текущей странице.
	Pagination Pagination    `json:"pagination"` // Информация о пагинации.
}
