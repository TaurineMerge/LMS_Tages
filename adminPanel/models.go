package main

import "time"

// ============ МОДЕЛИ ============

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type Category struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryCreate struct {
	Title string `json:"title"`
}

type CategoryUpdate struct {
	Title string `json:"title"`
}

type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	CategoryID  string    `json:"category_id"`
	Visibility  string    `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

type CourseUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

type PaginatedCourses struct {
	Data  []Course `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Pages int      `json:"pages"`
}

type Lesson struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CourseID  string                 `json:"course_id"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type LessonCreate struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

type LessonUpdate struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}
