package response

import "adminPanel/models"

// LessonResponse - ответ API с одним уроком
type LessonResponse struct {
	Status string                `json:"status"`
	Data   models.LessonDetailed `json:"data"`
}

// LessonListResponse - ответ со списком уроков
type LessonListResponse struct {
	Status string                               `json:"status"`
	Data   models.ResponsePaginationLessonsList `json:"data"`
}
