package main

type Course struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

type Lesson struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	CourseID string                 `json:"course_id"`
	Content  map[string]interface{} `json:"content,omitempty"`
}

type CreateCourseRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility,omitempty"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

type CreateLessonRequest struct {
	Title    string                 `json:"title"`
	CourseID string                 `json:"course_id"`
	Content  map[string]interface{} `json:"content,omitempty"`
}

type UpdateLessonRequest struct {
	Title   string                 `json:"title,omitempty"`
	Content map[string]interface{} `json:"content,omitempty"`
}