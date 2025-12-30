// Package request содержит структуры данных для разбора входящих HTTP-запросов.
package request

// ListQuery представляет собой стандартные параметры запроса для списков,
// включая пагинацию и сортировку.
type ListQuery struct {
	Page  int    `query:"page"`  // Номер страницы.
	Limit int    `query:"limit"` // Количество элементов на странице.
	Sort  string `query:"sort"`  // Поле и направление для сортировки (например, "created_at" или "-title").
}
