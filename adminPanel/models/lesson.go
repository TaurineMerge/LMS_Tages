package models

// Lesson - модель урока для списков (без контента)
type Lesson struct {
	BaseModel
	Title       string `json:"title"`
	CourseID    string `json:"course_id"`
	HTMLContent string `json:"html_content"` // HTML контент урока для WYSIWYG редактора
}

// LessonDetailed - детальная модель урока с контентом
type LessonDetailed struct {
	Lesson
	Content ContentSlice `json:"content"` // Структурированный контент (legacy)
}
