package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
)

// LessonPageViewModel содержит все данные, необходимые для отображения страницы урока.
// Это строго типизированная модель представления, предназначенная специально для шаблона 'pages/lesson.hbs'.
type LessonPageViewModel struct {
	// Вся основная информация об уроке
	Lesson     response.LessonDTODetailed
	Course     response.CourseDTO
	Category   response.CategoryDTO
	PrevLesson response.LessonDTO
	NextLesson response.LessonDTO
	AllLessons []response.LessonDTO

	// Вложенное поле, которое содержит данные для хедера
	PageHeader PageHeaderViewModel
}
