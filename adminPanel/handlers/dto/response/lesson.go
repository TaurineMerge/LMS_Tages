package response

import "adminPanel/models"

// LessonResponse представляет ответ API с одним уроком.
// Содержит статус и данные урока.
type LessonResponse struct {
	Status string        `json:"status"`
	Data   models.Lesson `json:"data"`
}

// LessonListResponse представляет ответ со списком уроков.
// Включает пагинацию и список уроков для курса.
type LessonListResponse struct {
	Status string                               `json:"status"`
	Data   models.ResponsePaginationLessonsList `json:"data"`
}
