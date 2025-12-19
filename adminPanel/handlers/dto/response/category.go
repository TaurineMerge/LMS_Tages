package response

import "adminPanel/models"

// CategoryResponse - ответ API с одной категорией
//
// Используется для возврата данных об одной категории.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: объект категории с полной информацией
type CategoryResponse struct {
	Status string          `json:"status"`
	Data   models.Category `json:"data"`
}

// PaginatedCategoriesResponse - пагинированный ответ со списком категорий
//
// Используется для возврата списка категорий с информацией о пагинации.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: данные ответа, содержащие список категорий и информацию о пагинации
//   - Data.Items: массив категорий
//   - Data.Pagination: информация о пагинации
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []models.Category `json:"items"`
		Pagination models.Pagination `json:"pagination"`
	} `json:"data"`
}
