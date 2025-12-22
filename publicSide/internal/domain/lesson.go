package domain

import "time"

type Lesson struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CourseID  string    `json:"course_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
