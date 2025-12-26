// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// LessonViewModel представляет данные для отображения одного урока в списке (например, в боковой панели).
type LessonViewModel struct {
	Title string
	Ref   string // URL-адрес урока.
}

// NewLessonViewModel создает новую модель представления для элемента списка уроков.
func NewLessonViewModel(lessonDTO response.LessonDTO, categoryID, courseID string) *LessonViewModel {
	vm := LessonViewModel{
		Title: lessonDTO.Title,
	}
	if lessonDTO.ID != "" {
		vm.Ref = routing.MakePathLesson(categoryID, courseID, lessonDTO.ID)
	}
	return &vm
}

// LessonDetailedViewModel расширяет LessonViewModel, добавляя контент урока для детального отображения.
type LessonDetailedViewModel struct {
	LessonViewModel
	Content string
}

// NewLessonDetailedViewModel создает новую модель представления для детальной информации об уроке.
func NewLessonDetailedViewModel(lessonDTO response.LessonDTODetailed, categoryId string) *LessonDetailedViewModel {
	return &LessonDetailedViewModel{
		LessonViewModel: LessonViewModel{
			Title: lessonDTO.Title,
			Ref:   routing.MakePathLesson(categoryId, lessonDTO.CourseID, lessonDTO.ID),
		},
		Content: lessonDTO.Content,
	}
}

// LessonPageViewModel представляет данные для страницы одного урока.
type LessonPageViewModel struct {
	PageHeader *PageHeaderViewModel
	Lesson     *LessonDetailedViewModel // Текущий урок.
	NextLesson *LessonViewModel         // Следующий урок для навигации.
	PrevLesson *LessonViewModel         // Предыдущий урок для навигации.
	Lessons    []LessonViewModel        // Полный список уроков курса для боковой панели.
}

// NewLessonPageViewModel создает новую модель представления для страницы урока.
func NewLessonPageViewModel(
	lessonDTODetailed response.LessonDTODetailed,
	courseDTO response.CourseDTO,
	categoryDTO response.CategoryDTO,
	nextLessonDTO response.LessonDTO,
	prevLessonDTO response.LessonDTO,
	lessonsDTOs []response.LessonDTO,
) *LessonPageViewModel {
	lessons := make([]LessonViewModel, len(lessonsDTOs))
	for i, ldto := range lessonsDTOs {
		lessons[i] = *NewLessonViewModel(ldto, categoryDTO.ID, courseDTO.ID)
	}

	return &LessonPageViewModel{
		PageHeader: NewPageHeaderViewModel("Урок: "+lessonDTODetailed.Title, BreadcrumbsForLessonPage(categoryDTO, courseDTO, lessonDTODetailed)),
		Lesson:     NewLessonDetailedViewModel(lessonDTODetailed, categoryDTO.ID),
		NextLesson: NewLessonViewModel(nextLessonDTO, categoryDTO.ID, courseDTO.ID),
		PrevLesson: NewLessonViewModel(prevLessonDTO, categoryDTO.ID, courseDTO.ID),
		Lessons:    lessons,
	}
}
