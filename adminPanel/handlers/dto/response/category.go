// Пакет response содержит структуры для ответов API.
package response

import "adminPanel/models"

// CategoryResponse представляет ответ API с одной категорией.
// Содержит статус и данные категории.
type CategoryResponse struct {
	Status string          `json:"status"`
	Data   models.Category `json:"data"`
}

// PaginatedCategoriesResponse представляет пагинированный ответ со списком категорий.
// Включает список категорий и информацию о пагинации.
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []models.Category `json:"items"`
		Pagination models.Pagination `json:"pagination"`
	} `json:"data"`
}
