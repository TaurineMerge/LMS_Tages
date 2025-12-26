// Пакет request содержит структуры для запросов API.
package request

// CategoryCreate представляет запрос на создание новой категории.
type CategoryCreate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryUpdate представляет запрос на обновление категории.
type CategoryUpdate struct {
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
}
