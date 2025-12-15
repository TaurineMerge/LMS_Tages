package models

// Category - модель категории курсов
//
// Представляет категорию, к которой могут принадлежать курсы.
// Категория содержит базовую информацию и наследует поля от BaseModel.
type Category struct {
	BaseModel
	Title string `json:"title"`
}

// CategoryCreate - DTO для создания новой категории
//
// Используется в запросах на создание категории.
// Содержит валидацию полей.
type CategoryCreate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryUpdate - DTO для обновления категории
//
// Используется в запросах на обновление категории.
// Все поля опциональны (omitempty).
type CategoryUpdate struct {
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
}

// CategoryResponse - ответ API с одной категорией
//
// Используется для возврата данных об одной категории.
type CategoryResponse struct {
	Status string   `json:"status"`
	Data   Category `json:"data"`
}

// PaginatedCategoriesResponse - пагинированный ответ со списком категорий
//
// Используется для возврата списка категорий с информацией о пагинации.
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Category `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}
