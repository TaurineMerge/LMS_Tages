package models

// Course represents a learning course.
type Course struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

// CourseCreate - DTO для создания курса
type CourseCreate struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
}

// CourseUpdate - DTO для обновления курса
type CourseUpdate struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"omitempty,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
}

// CourseResponse - ответ с курсом
type CourseResponse struct {
	Status string `json:"status"`
	Data   Course `json:"data"`
}

// PaginatedCoursesResponse - пагинированный ответ с курсами
type PaginatedCoursesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Course   `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}

// CourseFilter - фильтр для курсов
type CourseFilter struct {
	Level      string `query:"level"`
	Visibility string `query:"visibility"`
	CategoryID string `query:"category_id" validate:"omitempty,uuid4"`
	Page       int    `query:"page" validate:"min=1"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
}
