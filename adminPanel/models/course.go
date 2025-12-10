package models

// Убираем неиспользуемый импорт time
// Используйте time только если он нужен в структурах

type Course struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

// CourseResponse - ответ с курсом
type CourseResponse struct {
	Course
}

// CourseCreate - DTO для создания курса
type CourseCreate struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"required,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Visibility  string `json:"visibility" validate:"required,oneof=draft public private"`
}

// CourseUpdate - DTO для обновления курса
type CourseUpdate struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"omitempty,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
}

// PaginatedCourseResponse - пагинированный ответ с курсами
type PaginatedCourseResponse struct {
	Data  []CourseResponse `json:"data"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
	Pages int              `json:"pages"`
}

// CourseFilter - фильтр для курсов
type CourseFilter struct {
	Level      string `query:"level"`
	Visibility string `query:"visibility"`
	CategoryID string `query:"category_id" validate:"omitempty,uuid4"`
	Page       int    `query:"page" validate:"min=1"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
}