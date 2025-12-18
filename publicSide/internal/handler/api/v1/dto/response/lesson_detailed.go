package response

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// LessonDTODetailed is for the detailed view, including content.
type LessonDTODetailed struct {
	ID        string                `json:"id"`
	Title     string                `json:"title"`
	CourseID  string                `json:"course_id"`
	Content   []domain.ContentBlock `json:"content"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
