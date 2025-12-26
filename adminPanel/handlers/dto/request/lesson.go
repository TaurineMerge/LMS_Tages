// Пакет request содержит структуры для запросов API.
package request

// LessonCreate представляет запрос на создание нового урока.
// Содержит заголовок и содержимое урока.
type LessonCreate struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"omitempty"`
}

// LessonUpdate представляет запрос на обновление существующего урока.
// Все поля опциональны для частичного обновления.
type LessonUpdate struct {
	Title   string `json:"title" validate:"omitempty,min=1,max=255"`
	Content string `json:"content" validate:"omitempty"`
}
