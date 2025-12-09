package response

import "time"

// LessonDTO is for lists, without content.
type LessonDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CourseID  string    `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
