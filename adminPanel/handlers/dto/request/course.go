// Пакет request содержит структуры для запросов API.
package request

// CourseCreate представляет запрос на создание нового курса.
// Содержит все необходимые поля для создания курса с валидацией.
type CourseCreate struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
	ImageKey    string `json:"image_key"`
}

// CourseUpdate представляет запрос на обновление существующего курса.
// Все поля опциональны для частичного обновления.
type CourseUpdate struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"omitempty,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
	ImageKey    string `json:"image_key"`
}

// CourseFilter представляет фильтр для поиска курсов.
// Используется для пагинации и фильтрации по различным критериям.
type CourseFilter struct {
	Level      string `query:"level"`
	Visibility string `query:"visibility"`
	CategoryID string `query:"category_id" validate:"omitempty,uuid4"`
	Page       int    `query:"page" validate:"min=1"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
}
