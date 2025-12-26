// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
// Они агрегируют и форматируют данные из сервисного слоя для удобного отображения.
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// Breadcrumb представляет собой один элемент в навигационной цепочке "хлебных крошек".
type Breadcrumb struct {
	Text string // Текст элемента.
	URL  string // URL-адрес элемента. Пустая строка для некликабельного элемента.
}

// home создает базовый элемент "хлебных крошек" для главной страницы.
func home() Breadcrumb {
	return Breadcrumb{Text: "Главная", URL: routing.MakePathHome()}
}

// categories создает базовую цепочку "хлебных крошек" до страницы категорий.
func categories() []Breadcrumb {
	return []Breadcrumb{
		home(),
		{Text: "Категории", URL: routing.MakePathCategories()},
	}
}

// BreadcrumbsForCategoriesPage генерирует "хлебные крошки" для страницы со списком категорий.
// Последний элемент делается некликабельным.
func BreadcrumbsForCategoriesPage() []Breadcrumb {
	crumbs := categories()
	if len(crumbs) > 0 {
		crumbs[len(crumbs)-1].URL = "" // Последний элемент некликабельный
	}
	return crumbs
}

// BreadcrumbsForCoursesPage генерирует "хлебные крошки" для страницы курсов в определенной категории.
func BreadcrumbsForCoursesPage(category response.CategoryDTO) []Breadcrumb {
	crumbs := categories()
	crumbs = append(crumbs, Breadcrumb{Text: category.Title, URL: ""})
	return crumbs
}

// BreadcrumbsForCoursePage генерирует "хлебные крошки" для страницы конкретного курса.
func BreadcrumbsForCoursePage(category response.CategoryDTO, course response.CourseDTO) []Breadcrumb {
	crumbs := BreadcrumbsForCoursesPage(category)
	if len(crumbs) > 1 {
		// Предыдущий элемент (название категории) делаем кликабельным.
		crumbs[len(crumbs)-1].URL = routing.MakePathCourses(category.ID)
	}
	crumbs = append(crumbs, Breadcrumb{Text: course.Title, URL: ""})
	return crumbs
}

// BreadcrumbsForLessonPage генерирует "хлебные крошки" для страницы конкретного урока.
func BreadcrumbsForLessonPage(category response.CategoryDTO, course response.CourseDTO, lesson response.LessonDTODetailed) []Breadcrumb {
	crumbs := BreadcrumbsForCoursePage(category, course)
	if len(crumbs) > 1 {
		// Предыдущий элемент (название курса) делаем кликабельным.
		crumbs[len(crumbs)-1].URL = routing.MakePathCourse(category.ID, course.ID)
	}
	crumbs = append(crumbs, Breadcrumb{Text: lesson.Title, URL: ""})
	return crumbs
}
