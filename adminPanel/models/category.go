package models

// Category - модель категории курсов
//
// Представляет категорию, к которой могут принадлежать курсы.
// Категория содержит базовую информацию и наследует поля от BaseModel.
//
// Поля:
//   - BaseModel: встроенная структура с общими полями (ID, CreatedAt, UpdatedAt)
//   - Title: название категории
type Category struct {
	BaseModel
	Title string `json:"title"`
}

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

// CategoryResponse - ответ API с одной категорией
//
// Используется для возврата данных об одной категории.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: объект категории с полной информацией
type CategoryResponse struct {
	Status string   `json:"status"`
	Data   Category `json:"data"`
}

// PaginatedCategoriesResponse - пагинированный ответ со списком категорий
//
// Используется для возврата списка категорий с информацией о пагинации.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Data: данные ответа, содержащие список категорий и информацию о пагинации
//   - Data.Items: массив категорий
//   - Data.Pagination: информация о пагинации
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Category `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}
