package request

// LessonCreate - DTO для создания нового урока
type LessonCreate struct {
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Content string `json:"content" validate:"omitempty"`
}

// LessonUpdate - DTO для обновления урока
type LessonUpdate struct {
	Title   string `json:"title" validate:"omitempty,min=1,max=255"`
	Content string `json:"content" validate:"omitempty"`
}
