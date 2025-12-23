package models

// Lesson - модель урока для списков (без контента)
type Lesson struct {
	BaseModel
	Title    string `json:"title"`
	CourseID string `json:"course_id"`
	Content  string `json:"content"`
}

// LessonDetailed - детальная модель урока с контентом
type LessonDetailed struct {
	Lesson
}
