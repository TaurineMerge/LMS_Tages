package entity

import "time"

// Lesson represents a lesson entity
// @Description Lesson information
type Lesson struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	CourseID    int64     `json:"course_id" gorm:"index;not null"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Content     string    `json:"content" gorm:"type:text"`
	Slug        string    `json:"slug" gorm:"size:255;not null"`
	OrderIndex  int       `json:"order_index" gorm:"default:0"`
	IsPublished bool      `json:"is_published" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Course      *Course   `json:"course,omitempty" gorm:"foreignKey:CourseID"`
}

// LessonResponse represents a lesson response for public API
// @Description Lesson response with course info
type LessonResponse struct {
	ID         int64           `json:"id"`
	CourseID   int64           `json:"course_id"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	Slug       string          `json:"slug"`
	OrderIndex int             `json:"order_index"`
	CreatedAt  time.Time       `json:"created_at"`
	Course     *CourseResponse `json:"course,omitempty"`
}
