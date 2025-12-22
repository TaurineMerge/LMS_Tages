package request

// CategoryCreate - DTO для создания новой категории
//
// Используется в запросах на создание категории.
// Содержит валидацию полей.
//
// Поля:
//   - Title: название категории (обязательное, от 1 до 255 символов)
type CategoryCreate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryUpdate - DTO для обновления категории
//
// Используется в запросах на обновление категории.
// Все поля опциональны (omitempty).
//
// Поля:
//   - Title: новое название категории (опционально, от 1 до 255 символов)
type CategoryUpdate struct {
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
}
