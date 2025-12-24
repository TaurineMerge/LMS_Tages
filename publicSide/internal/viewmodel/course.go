package viewmodel

import (
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

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

func NewCourseViewModel(courseDTO *response.CourseDTO, lessonAmount int) *CourseViewModel {
	return &CourseViewModel{
		Title:         courseDTO.Title,
		Ref:           routing.MakePathCourse(courseDTO.CategoryID, courseDTO.ID),
		Level:         courseDTO.Level,
		LevelRu:       "ПУСТО!!!",
		Description:   courseDTO.Description,
		LessonsAmount: lessonAmount,
		UpdatedAt:     courseDTO.UpdatedAt,
		CreatedAt:     courseDTO.CreatedAt,
		ImageURL:      courseDTO.ImageURL,
	}
}

type CourseDetailViewModel struct {
	CourseViewModel
	Lessons []LessonViewModel
}

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

type CoursesPageViewModel struct {
	PageHeader *PageHeaderViewModel
	Courses    []CourseViewModel
	Pagination *PaginationViewModel
	Level      string
	SortBy     string
}

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

type CoursePageViewModel struct {
	PageHeader *PageHeaderViewModel
	Course     *CourseDetailViewModel
}

func NewCoursePageViewModel(categoryDTO response.CategoryDTO, courseDTO response.CourseDTO, lessonsDTO []response.LessonDTO) *CoursePageViewModel {
	return &CoursePageViewModel{
		PageHeader: NewPageHeaderViewModel("Курс: "+courseDTO.Title, BreadcrumbsForCoursePage(categoryDTO, courseDTO)),
		Course:     NewCourseDetailViewModel(&courseDTO, lessonsDTO),
	}
}
