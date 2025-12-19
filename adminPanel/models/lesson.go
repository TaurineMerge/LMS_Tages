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
