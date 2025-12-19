package models

import "time"

// Category represents the structure for a category.
type Category struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Course represents the structure for a course.
type Course struct {
	ID          string
	Title       string
	Description string
	Level       string
	Visibility  string
	CategoryID  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Lesson represents the structure for a lesson.
type Lesson struct {
	ID        string
	Title     string
	CourseID  string
	Content   []ContentBlock
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ContentBlock represents a block of content within a lesson.
type ContentBlock struct {
	ContentType string // "text" or "image"
	Data        string
}
