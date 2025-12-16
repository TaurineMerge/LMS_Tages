package models

// Lesson - модель урока курса
//
// Представляет урок, который является частью учебного курса.
// Урок содержит базовую информацию о себе.
//
// Поля:
//   - BaseModel: встроенная структура с общими полями (ID, CreatedAt, UpdatedAt)
//   - Title: название урока
//   - CategoryID: уникальный идентификатор категории
//   - CourseID: уникальный идентификатор курса, к которому принадлежит урок
type Lesson struct {
	BaseModel
	Title      string `json:"title"`
	CategoryID string `json:"category_id"`
	CourseID   string `json:"course_id"`
}

// LessonDetailed - детальная модель урока с контентом
//
// Расширенная модель урока, включающая в себя содержимое урока.
//
// Поля:
//   - Lesson: встроенная базовая модель урока
//   - Content: содержимое урока в формате JSON (map[string]interface{})
type LessonDetailed struct {
	Lesson
	Content map[string]interface{} `json:"content"`
}

// LessonResponse - ответ API с одним уроком
//
// Используется для возврата данных об одном уроке с его содержимым.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: объект детального урока с содержимым
type LessonResponse struct {
	Status string         `json:"status"`
	Data   LessonDetailed `json:"data"`
}

// LessonCreate - DTO для создания нового урока
//
// Используется в запросах на создание урока.
// Содержит валидацию полей.
//
// Поля:
//   - Title: название урока (обязательное, от 1 до 255 символов)
//   - CategoryID: идентификатор категории (опционально, UUID v4)
//   - Content: содержимое урока в формате JSON
type LessonCreate struct {
	Title      string                 `json:"title" validate:"required,min=1,max=255"`
	CategoryID string                 `json:"category_id" validate:"omitempty,uuid4"`
	Content    map[string]interface{} `json:"content"`
}

// LessonUpdate - DTO для обновления урока
//
// Используется в запросах на обновление урока.
// Все поля опциональны (omitempty).
//
// Поля:
//   - Title: новое название урока (опционально, от 1 до 255 символов)
//   - CategoryID: новый идентификатор категории (опционально, UUID v4)
//   - Content: новое содержимое урока в формате JSON (опционально)
type LessonUpdate struct {
	Title      string                 `json:"title" validate:"omitempty,min=1,max=255"`
	CategoryID string                 `json:"category_id" validate:"omitempty,uuid4"`
	Content    map[string]interface{} `json:"content"`
}

// LessonListResponse - ответ со списком уроков
//
// Используется для возврата списка уроков с информацией о пагинации.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: данные ответа, содержащие список уроков и информацию о пагинации
//   - Data.Items: массив уроков
//   - Data.Pagination: информация о пагинации
type LessonListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Lesson   `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}
