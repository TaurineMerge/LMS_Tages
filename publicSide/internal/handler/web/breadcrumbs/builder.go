package breadcrumbs

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// home возвращает базовую крошку для главной страницы.
func home() viewmodel.Breadcrumb {
	return viewmodel.Breadcrumb{Text: "Главная", URL: routing.MakePathHome()}
}

// categories возвращает путь до списка категорий.
func categories() []viewmodel.Breadcrumb {
	return []viewmodel.Breadcrumb{
		home(),
		{Text: "Категории", URL: routing.MakePathCategories()},
	}
}

// ForCategoriesPage создает крошки для страницы со списком всех категорий.
func ForCategoriesPage() []viewmodel.Breadcrumb {
	crumbs := categories()
	if len(crumbs) > 0 {
		crumbs[len(crumbs)-1].URL = "" // Последний элемент не кликабельный
	}
	return crumbs
}

// ForCoursesPage создает крошки для страницы курсов конкретной категории.
func ForCoursesPage(category response.CategoryDTO) []viewmodel.Breadcrumb {
	crumbs := categories()
	crumbs = append(crumbs, viewmodel.Breadcrumb{Text: category.Title, URL: ""})
	return crumbs
}

// ForCoursePage создает крошки для страницы конкретного курса.
func ForCoursePage(category response.CategoryDTO, course response.CourseDTO) []viewmodel.Breadcrumb {
	crumbs := ForCoursesPage(category)
	if len(crumbs) > 1 {
		// Делаем предыдущий элемент (категорию) кликабельным
		crumbs[len(crumbs)-1].URL = routing.MakePathCourses(category.ID)
	}
	// Добавляем текущий, некликабельный элемент
	crumbs = append(crumbs, viewmodel.Breadcrumb{Text: course.Title, URL: ""})
	return crumbs
}

// ForLessonPage создает крошки для страницы урока.
func ForLessonPage(category response.CategoryDTO, course response.CourseDTO, lesson response.LessonDTODetailed) []viewmodel.Breadcrumb {
	crumbs := ForCoursePage(category, course)
	if len(crumbs) > 1 {
		// Делаем предыдущий элемент (курс) кликабельным
		crumbs[len(crumbs)-1].URL = routing.MakePathCourse(category.ID, course.ID)
	}
	// Добавляем текущий, некликабельный элемент
	crumbs = append(crumbs, viewmodel.Breadcrumb{Text: lesson.Title, URL: ""})
	return crumbs
}
