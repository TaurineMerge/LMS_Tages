// Package models содержит базовые структуры данных для приложения adminPanel.
// Включает общие модели, такие как BaseModel, и структуры для пагинации и запросов.
package models

import "time"

// BaseModel представляет базовую модель с общими полями для всех сущностей.
// Содержит ID, время создания и обновления.
type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Pagination содержит информацию о пагинации для списков.
// Включает общее количество элементов, текущую страницу, лимит и общее количество страниц.
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// QueryList представляет параметры запроса для получения списков.
// Используется для парсинга query-параметров: page, limit, sort.
type QueryList struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

// ResponsePaginationLessonsList представляет ответ с пагинированным списком уроков.
// Содержит массив уроков и информацию о пагинации.
type ResponsePaginationLessonsList struct {
	Items      []Lesson   `json:"items"`
	Pagination Pagination `json:"pagination"`
}
