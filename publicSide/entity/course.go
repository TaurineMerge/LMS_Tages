package entity

import "time"

// Course represents a course entity
// @Description Course information
type Course struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Slug        string    `json:"slug" gorm:"size:255;uniqueIndex;not null"`
	IsPublished bool      `json:"is_published" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Lessons     []Lesson  `json:"lessons,omitempty" gorm:"foreignKey:CourseID"`
}

// CourseResponse represents a course response for public API
// @Description Course response with lessons count
type CourseResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Slug         string    `json:"slug"`
	LessonsCount int       `json:"lessons_count"`
	CreatedAt    time.Time `json:"created_at"`
}
