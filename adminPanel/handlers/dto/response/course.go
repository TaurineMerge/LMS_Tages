package response

import "adminPanel/models"

// CourseResponse - ответ API с одним курсом
//
// Используется для возврата данных об одном курсе.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: объект курса с полной информацией
type CourseResponse struct {
	Status string        `json:"status"`
	Data   models.Course `json:"data"`
}

// PaginatedCoursesResponse - ответ со списком курсов
//
// Используется для возврата списка курсов с информацией о пагинации.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: данные ответа, содержащие список курсов и информацию о пагинации
//   - Data.Items: массив курсов
//   - Data.Pagination: информация о пагинации (текущая страница, общее количество и т.д.)
type PaginatedCoursesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []models.Course   `json:"items"`
		Pagination models.Pagination `json:"pagination"`
	} `json:"data"`
}
