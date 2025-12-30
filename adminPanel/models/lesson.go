package models

// Lesson представляет урок в системе.
// Встраивает BaseModel и содержит поля для заголовка, ID курса и контента урока.
type Lesson struct {
	BaseModel
	Title    string `json:"title"`
	CourseID string `json:"course_id"`
	Content  string `json:"content"`
}
