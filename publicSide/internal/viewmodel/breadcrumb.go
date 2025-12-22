package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// Breadcrumb представляет один элемент "хлебных крошек".
type Breadcrumb struct {
	Text string
	URL  string // URL может быть пустым для последнего, неактивного элемента
}

// home возвращает базовую крошку для главной страницы.
func home() Breadcrumb {
	return Breadcrumb{Text: "Главная", URL: routing.MakePathHome()}
}

// categories возвращает путь до списка категорий.
func categories() []Breadcrumb {
	return []Breadcrumb{
		home(),
		{Text: "Категории", URL: routing.MakePathCategories()},
	}
}

// BreadcrumbsForCategoriesPage создает крошки для страницы со списком всех категорий.
func BreadcrumbsForCategoriesPage() []Breadcrumb {
	crumbs := categories()
	if len(crumbs) > 0 {
		crumbs[len(crumbs)-1].URL = "" // Последний элемент не кликабельный
	}
	return crumbs
}

// BreadcrumbsForCoursesPage создает крошки для страницы курсов конкретной категории.
func BreadcrumbsForCoursesPage(category response.CategoryDTO) []Breadcrumb {
	crumbs := categories()
	crumbs = append(crumbs, Breadcrumb{Text: category.Title, URL: ""})
	return crumbs
}

// BreadcrumbsForCoursePage создает крошки для страницы конкретного курса.
func BreadcrumbsForCoursePage(category response.CategoryDTO, course response.CourseDTO) []Breadcrumb {
	crumbs := BreadcrumbsForCoursesPage(category)
	if len(crumbs) > 1 {
		// Делаем предыдущий элемент (категорию) кликабельным
		crumbs[len(crumbs)-1].URL = routing.MakePathCourses(category.ID)
	}
	// Добавляем текущий, некликабельный элемент
	crumbs = append(crumbs, Breadcrumb{Text: course.Title, URL: ""})
	return crumbs
}

// BreadcrumbsForLessonPage создает крошки для страницы урока.
func BreadcrumbsForLessonPage(category response.CategoryDTO, course response.CourseDTO, lesson response.LessonDTODetailed) []Breadcrumb {
	crumbs := BreadcrumbsForCoursePage(category, course)
	if len(crumbs) > 1 {
		// Делаем предыдущий элемент (курс) кликабельным
		crumbs[len(crumbs)-1].URL = routing.MakePathCourse(category.ID, course.ID)
	}
	// Добавляем текущий, некликабельный элемент
	crumbs = append(crumbs, Breadcrumb{Text: lesson.Title, URL: ""})
	return crumbs
}
