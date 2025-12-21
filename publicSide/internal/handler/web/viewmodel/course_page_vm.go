package viewmodel

import "github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"

// CoursePageViewModel содержит все данные для страницы детального отображения курса.
type CoursePageViewModel struct {
	PageHeader    PageHeaderViewModel
	CategoryTitle string
	CategoryID    string
	Course        response.CourseDTO
	LevelRu       string
	FirstLessonID string
}
