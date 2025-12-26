// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// CourseViewModel представляет данные для отображения одной карточки курса в списке.
type CourseViewModel struct {
	Title         string
	Ref           string
	Level         string
	LevelRu       string
	Description   string
	LessonsAmount int
	UpdatedAt     time.Time
	CreatedAt     time.Time
	ImageURL      string
}

// NewCourseViewModel создает новую модель представления для карточки курса.
func NewCourseViewModel(courseDTO *response.CourseDTO, lessonAmount int) *CourseViewModel {
	return &CourseViewModel{
		Title:         courseDTO.Title,
		Ref:           routing.MakePathCourse(courseDTO.CategoryID, courseDTO.ID),
		Level:         courseDTO.Level,
		LevelRu:       "ПУСТО!!!", // Это поле заполняется позже в обработчике
		Description:   courseDTO.Description,
		LessonsAmount: lessonAmount,
		UpdatedAt:     courseDTO.UpdatedAt,
		CreatedAt:     courseDTO.CreatedAt,
		ImageURL:      courseDTO.ImageURL,
	}
}

// CourseDetailViewModel расширяет CourseViewModel, добавляя список уроков для детальной страницы курса.
type CourseDetailViewModel struct {
	CourseViewModel
	Lessons []LessonViewModel
}

// NewCourseDetailViewModel создает новую модель представления для детальной информации о курсе.
func NewCourseDetailViewModel(courseDTO *response.CourseDTO, lessonsDTO []response.LessonDTO) *CourseDetailViewModel {
	lessons := make([]LessonViewModel, 0, len(lessonsDTO))
	for _, lDTO := range lessonsDTO {
		lessons = append(lessons, *NewLessonViewModel(lDTO, courseDTO.CategoryID, courseDTO.ID))
	}

	return &CourseDetailViewModel{
		CourseViewModel: *NewCourseViewModel(courseDTO, len(lessonsDTO)),
		Lessons:         lessons,
	}
}

// CoursesPageViewModel представляет данные для страницы со списком всех курсов в категории.
type CoursesPageViewModel struct {
	PageHeader *PageHeaderViewModel
	Courses    []CourseViewModel
	Pagination *PaginationViewModel
	Level      string // Текущий выбранный фильтр уровня
	SortBy     string // Текущий выбранный метод сортировки
}

// NewCoursesPageViewModel создает новую модель представления для страницы списка курсов.
func NewCoursesPageViewModel(categoryDTO response.CategoryDTO, coursesDTO []response.CourseDTO, coursesPagination response.Pagination, lessonsAmount []int, level string, sortBy string) *CoursesPageViewModel {
	courses := make([]CourseViewModel, 0, len(coursesDTO))
	for i, c := range coursesDTO {
		courses = append(courses, *NewCourseViewModel(&c, lessonsAmount[i]))
	}

	return &CoursesPageViewModel{
		PageHeader: NewPageHeaderViewModel("Курсы в категории: "+categoryDTO.Title, BreadcrumbsForCoursesPage(categoryDTO)),
		Courses:    courses,
		Pagination: NewPaginationViewModel(coursesPagination, routing.MakePathCourses(categoryDTO.ID)),
		Level:      level,
		SortBy:     sortBy,
	}
}

// CoursePageViewModel представляет данные для детальной страницы одного курса.
type CoursePageViewModel struct {
	PageHeader               *PageHeaderViewModel
	Course                   *CourseDetailViewModel
	Test                     *TestViewModel
	TestIsNotFound           bool // Флаг, что тест для курса не найден.
	TestServiceIsUnavailable bool // Флаг, что сервис тестов недоступен.
}

// NewCoursePageViewModel создает новую модель представления для страницы курса.
func NewCoursePageViewModel(
	categoryDTO response.CategoryDTO,
	courseDTO response.CourseDTO,
	lessonsDTO []response.LessonDTO,
	testVM *TestViewModel,
	testIsNotFound bool,
	testServiceIsUnavailable bool,
) *CoursePageViewModel {
	return &CoursePageViewModel{
		PageHeader:               NewPageHeaderViewModel("Курс: "+courseDTO.Title, BreadcrumbsForCoursePage(categoryDTO, courseDTO)),
		Course:                   NewCourseDetailViewModel(&courseDTO, lessonsDTO),
		Test:                     testVM,
		TestIsNotFound:           testIsNotFound,
		TestServiceIsUnavailable: testServiceIsUnavailable,
	}
}
