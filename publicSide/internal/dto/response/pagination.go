// Package response содержит структуры данных для формирования HTTP-ответов.
package response

// Pagination содержит информацию, необходимую для постраничной навигации.
type Pagination struct {
	Page  int `json:"page"`  // Текущий номер страницы.
	Limit int `json:"limit"` // Количество элементов на странице.
	Total int `json:"total"` // Общее количество элементов.
	Pages int `json:"pages"` // Общее количество страниц.
}
