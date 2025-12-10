package models

// Убираем неиспользуемый импорт time
// Используйте time только если он нужен в структурах

type Lesson struct {
	BaseModel
	Title    string                 `json:"title"`
	CourseID string                 `json:"course_id"`
	Content  map[string]interface{} `json:"content"`
}

// LessonResponse - ответ с уроком
type LessonResponse struct {
	Lesson
}

// LessonCreate - DTO для создания урока
type LessonCreate struct {
	Title   string                 `json:"title" validate:"required,min=1,max=255"`
	Content map[string]interface{} `json:"content"`
}

// LessonUpdate - DTO для обновления урока
type LessonUpdate struct {
	Title   string                 `json:"title" validate:"omitempty,min=1,max=255"`
	Content map[string]interface{} `json:"content"`
}

// LessonListResponse - список уроков
type LessonListResponse struct {
	Data  []LessonResponse `json:"data"`
	Total int              `json:"total"`
}