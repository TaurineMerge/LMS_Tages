package models

// Убираем неиспользуемый импорт time
// Используйте time только если он нужен в структурах

type Lesson struct {
	BaseModel
	Title      string `json:"title"`
	CategoryID string `json:"category_id"`
	CourseID   string `json:"course_id"`
}

// LessonDetailed - детальный урок с контентом
type LessonDetailed struct {
	Lesson
	Content map[string]interface{} `json:"content"`
}

// LessonResponse - ответ с уроком
type LessonResponse struct {
	Status string         `json:"status"`
	Data   LessonDetailed `json:"data"`
}

// LessonCreate - DTO для создания урока
type LessonCreate struct {
	Title      string                 `json:"title" validate:"required,min=1,max=255"`
	CategoryID string                 `json:"category_id" validate:"omitempty,uuid4"`
	Content    map[string]interface{} `json:"content"`
}

// LessonUpdate - DTO для обновления урока
type LessonUpdate struct {
	Title      string                 `json:"title" validate:"omitempty,min=1,max=255"`
	CategoryID string                 `json:"category_id" validate:"omitempty,uuid4"`
	Content    map[string]interface{} `json:"content"`
}

// LessonListResponse - список уроков
type LessonListResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Lesson   `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}
