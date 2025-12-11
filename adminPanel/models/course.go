package models

// Course - модель учебного курса
//
// Представляет учебный курс, который может содержать уроки.
// Курс принадлежит определенной категории и имеет уровень сложности.
type Course struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

// CourseCreate - DTO для создания нового курса
//
// Используется в запросах на создание курса.
// Содержит валидацию полей.
type CourseCreate struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
}

// CourseUpdate - DTO для обновления курса
//
// Используется в запросах на обновление курса.
// Все поля опциональны (omitempty).
type CourseUpdate struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"omitempty,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
}

// CourseResponse - ответ API с одним курсом
//
// Используется для возврата данных об одном курсе.
type CourseResponse struct {
	Status string `json:"status"`
	Data   Course `json:"data"`
}

// PaginatedCoursesResponse - пагинированный ответ со списком курсов
//
// Используется для возврата списка курсов с информацией о пагинации.
type PaginatedCoursesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Course   `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}

// CourseFilter - фильтр для поиска курсов
//
// Используется для фильтрации курсов по различным критериям:
// уровень сложности, видимость, категория, пагинация.
type CourseFilter struct {
	Level      string `query:"level"`
	Visibility string `query:"visibility"`
	CategoryID string `query:"category_id" validate:"omitempty,uuid4"`
	Page       int    `query:"page" validate:"min=1"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
}
