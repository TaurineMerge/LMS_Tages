// Package models содержит модели данных для передачи и доменные модели
//
// Пакет предоставляет:
//   - BaseModel: базовую модель с общими полями
//   - Pagination: модель для пагинации
//   - Category: модель категории
//   - Course: модель курса
//   - Lesson: модель урока
//   - Response: модели для ответов API
package models

import "time"

// BaseModel - базовая модель, содержащая общие поля для всех сущностей
//
// Включает в себя:
//   - ID: уникальный идентификатор
//   - CreatedAt: дата и время создания
//   - UpdatedAt: дата и время последнего обновления
type BaseModel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Pagination - модель для пагинации результатов
//
// Соответствует объекту Pagination в Swagger спецификации.
// Используется для постраничного вывода данных.
type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

// QueryList - общая структура для обработки параметров запроса со списком
//
// Содержит:
//   - Page: номер текущей страницы
//   - Limit: количество элементов на странице
//   - Sort: строка сортировки (например, "created_at" или "-title")
type QueryList struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"`
}

type ResponsePaginationLessonsList struct {
	Items      []Lesson   `json:"items"`
	Pagination Pagination `json:"pagination"`
}
