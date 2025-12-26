// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// CategoryViewModel представляет данные для отображения одной карточки категории.
type CategoryViewModel struct {
	Title        string
	CoursesRef   string
	TotalCourses int
	Courses      []CourseViewModel
	CoursesLimit int
}

// NewCategoryViewModel создает новую модель представления для одной категории.
func NewCategoryViewModel(categoryDTO response.CategoryDTO, coursesDTO []response.CourseDTO, coursesPagination response.Pagination, coursesLimit int) CategoryViewModel {
	courses := make([]CourseViewModel, 0, len(coursesDTO))
	for _, c := range coursesDTO {
		courses = append(courses, *NewCourseViewModel(&c, 0))
	}

	return CategoryViewModel{
		Title:        categoryDTO.Title,
		CoursesRef:   routing.MakePathCourses(categoryDTO.ID),
		TotalCourses: coursesPagination.Total,
		Courses:      courses,
		CoursesLimit: coursesLimit,
	}
}

// CategoriesPageViewModel представляет данные для страницы со списком всех категорий.
type CategoriesPageViewModel struct {
	PageHeader *PageHeaderViewModel
	Categories []CategoryViewModel
	Pagination *PaginationViewModel
}

// NewCategoriesPageViewMode создает новую модель представления для страницы категорий.
func NewCategoriesPageViewMode(
	categories []CategoryViewModel,
	pagination response.Pagination,
) *CategoriesPageViewModel {
	return &CategoriesPageViewModel{
		PageHeader: NewPageHeaderViewModel("Категории курсов", BreadcrumbsForCategoriesPage()),
		Categories: categories,
		Pagination: NewPaginationViewModel(pagination, routing.MakePathCategories()),
	}
}
