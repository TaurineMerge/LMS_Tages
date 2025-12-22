package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

type LessonViewModel struct {
	Title string
	Ref   string
}

func NewLessonViewModel(lessonDTO response.LessonDTO, categoryID, courseID string) *LessonViewModel {
	vm := LessonViewModel{
		Title: lessonDTO.Title,
	}
	if lessonDTO.ID != "" {
		vm.Ref = routing.MakePathLesson(categoryID, courseID, lessonDTO.ID)
	}
	return &vm
}

type LessonDetailedViewModel struct {
	LessonViewModel
	Content string
}

func NewLessonDetailedViewModel(lessonDTO response.LessonDTODetailed, categoryId string) *LessonDetailedViewModel {
	return &LessonDetailedViewModel{
		LessonViewModel: LessonViewModel{
			Title: lessonDTO.Title,
			Ref:   routing.MakePathLesson(categoryId, lessonDTO.CourseID, lessonDTO.ID),
		},
		Content: lessonDTO.Content,
	}
}

type LessonPageViewModel struct {
	PageHeader *PageHeaderViewModel
	Lesson     *LessonDetailedViewModel
	NextLesson *LessonViewModel
	PrevLesson *LessonViewModel
	Lessons    []LessonViewModel
}

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
