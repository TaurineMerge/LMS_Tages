package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
)

// CourseView представляет курс в списке для страницы courses.
type CourseView struct {
	ID          string
	Title       string
	Description string
	Level       string
	LevelRu     string
}

// CoursesPageViewModel содержит все данные для страницы со списком курсов в категории.
type CoursesPageViewModel struct {
	PageHeader    PageHeaderViewModel
	CategoryTitle string
	CategoryID    string
	Courses       []CourseView
	Pagination    response.Pagination
	CurrentPage   int
	Level         string // for filter state
	SortBy        string // for filter state
}
