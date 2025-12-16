package models

// Lesson - модель урока для списков (без контента)
type Lesson struct {
	BaseModel
	Title    string `json:"title"`
	CourseID string `json:"course_id"`
}

// LessonDetailed - детальная модель урока с контентом
type LessonDetailed struct {
	Lesson
	Content ContentSlice `json:"content"`
}

// LessonCreate - DTO для создания нового урока
type LessonCreate struct {
	Title   string       `json:"title" validate:"required,min=1,max=255"`
	Content ContentSlice `json:"content"`
}

// LessonUpdate - DTO для обновления урока
type LessonUpdate struct {
	Title   string       `json:"title" validate:"omitempty,min=1,max=255"`
	Content ContentSlice `json:"content"`
}

// LessonResponse - ответ API с одним уроком
type LessonResponse struct {
	Status string         `json:"status"`
	Data   LessonDetailed `json:"data"`
}

// LessonListResponse - ответ со списком уроков
type LessonListResponse struct {
	Status string                        `json:"status"`
	Data   ResponsePaginationLessonsList `json:"data"`
}
