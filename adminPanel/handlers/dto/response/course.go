package response

import "adminPanel/models"

// CourseResponse представляет ответ API с одним курсом.
// Содержит статус и данные курса.
type CourseResponse struct {
	Status string        `json:"status"`
	Data   models.Course `json:"data"`
}

// PaginatedCoursesResponse представляет пагинированный ответ со списком курсов.
// Включает список курсов и информацию о пагинации.
type PaginatedCoursesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []models.Course   `json:"items"`
		Pagination models.Pagination `json:"pagination"`
	} `json:"data"`
}
