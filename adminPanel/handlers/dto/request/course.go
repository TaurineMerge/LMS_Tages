package request

// CourseCreate - DTO для создания нового курса
//
// Используется в запросах на создание курса.
// Содержит валидацию полей.
//
// Поля:
//   - Title: название курса (обязательное, от 1 до 255 символов)
//   - Description: описание курса (опционально)
//   - Level: уровень сложности (опционально, "hard", "medium", "easy")
//   - CategoryID: идентификатор категории (обязательное, UUID v4)
//   - Visibility: видимость курса (опционально, "draft", "public", "private")
//   - ImageKey: ключ изображения в S3 (опционально)
type CourseCreate struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
	ImageKey    string `json:"image_key"`
}

// CourseUpdate - DTO для обновления курса
//
// Используется в запросах на обновление курса.
// Все поля опциональны (omitempty).
//
// Поля:
//   - Title: новое название курса (опционально, от 1 до 255 символов)
//   - Description: новое описание курса (опционально)
//   - Level: новый уровень сложности (опционально, "hard", "medium", "easy")
//   - CategoryID: новый идентификатор категории (опционально, UUID v4)
//   - Visibility: новая видимость курса (опционально, "draft", "public", "private")
//   - ImageKey: новый ключ изображения в S3 (опционально)
type CourseUpdate struct {
	Title       string `json:"title" validate:"omitempty,min=1,max=255"`
	Description string `json:"description"`
	Level       string `json:"level" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"omitempty,uuid4"`
	Visibility  string `json:"visibility" validate:"omitempty,oneof=draft public private"`
	ImageKey    string `json:"image_key"`
}

// CourseFilter - фильтр для поиска курсов
//
// Используется для фильтрации курсов по различным критериям:
// уровень сложности, видимость, категория, пагинация.
//
// Поля:
//   - Level: уровень сложности курса ("hard", "medium", "easy")
//   - Visibility: видимость курса ("draft", "public", "private")
//   - CategoryID: уникальный идентификатор категории для фильтрации
//   - Page: номер страницы для пагинации (минимум 1)
//   - Limit: количество элементов на странице (от 1 до 100)
type CourseFilter struct {
	Level      string `query:"level"`
	Visibility string `query:"visibility"`
	CategoryID string `query:"category_id" validate:"omitempty,uuid4"`
	Page       int    `query:"page" validate:"min=1"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
}
