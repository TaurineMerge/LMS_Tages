package domain

import "time"

// ContentBlock represents a block of content in a lesson (text or image).
type ContentBlock struct {
	ContentType string `json:"content_type"` // "text" or "image"
	Data        string `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
	Alt         string `json:"alt,omitempty"`
}

// Lesson represents a lesson entity.
type Lesson struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	CourseID  string         `json:"course_id"`
	Content   []ContentBlock `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
