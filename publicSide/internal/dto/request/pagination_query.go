// Package request содержит структуры данных для разбора входящих HTTP-запросов.
package request

// PaginationQuery представляет собой стандартные параметры запроса для пагинации.
type PaginationQuery struct {
	Page  int `query:"page"`  // Номер страницы.
	Limit int `query:"limit"` // Количество элементов на странице.
}
