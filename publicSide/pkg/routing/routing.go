// Package routing определяет централизованные константы и функции для управления
// маршрутизацией в веб-приложении. Это помогает избежать "магических строк" и
// обеспечивает единый источник правды для всех URL-адресов.
package routing

import "fmt"

// --- Path Variable Names (для Fiber) ---
// Эти константы определяют имена переменных в URL-путях, которые Fiber использует
// для извлечения параметров. Например, в /categories/:category_id,
// `c.Params(PathVariableCategoryID)` вернет значение :category_id.
const (
	PathVariableCategoryID = "category_id" // Имя переменной для ID категории.
	PathVariableCourseID   = "course_id"   // Имя переменной для ID курса.
	PathVariableLessonID   = "lesson_id"   // Имя переменной для ID урока.
)

// --- Route Definitions (для шаблонов Fiber `app.Get` и `app.Group`) ---
// Эти константы представляют собой шаблоны путей, используемые при определении
// маршрутов в приложении Fiber. Они содержат плейсхолдеры (например, :category_id),
// которые Fiber сопоставляет с PathVariable-константами.
const (
	// Основные маршруты
	RouteHome         = "/"
	RouteLogin        = "/login"
	RouteLogout       = "/logout"
	RouteAuthCallback = "/auth/callback"
	RouteReg          = "/reg"

	// Внешние сервисы
	ExternalServiceRouteProfile = "http://localhost/account/profile"

	// API
	RouteAPIV1 = "/api/v1"

	// Ресурсы
	RouteCategories = "/categories"
	RouteCategory   = "/categories/:" + PathVariableCategoryID
	RouteCourses    = "/categories/:" + PathVariableCategoryID + "/courses"
	RouteCourse     = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID
	RouteLessons    = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID + "/lessons"
	RouteLesson     = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID + "/lessons/:" + PathVariableLessonID
)

// --- Path Constructors (для генерации URL в шаблонах, редиректах и т.д.) ---
// Эти функции создают конкретные, готовые к использованию URL-адреса.
// Они принимают реальные ID и подставляют их в шаблоны маршрутов, обеспечивая
// консистентность и избегая ошибок при ручном формировании URL.

// MakePathHome создает путь к домашней странице.
func MakePathHome() string {
	return RouteHome
}

// MakePathCategories создает путь к странице со списком всех категорий.
func MakePathCategories() string {
	return RouteCategories
}

// MakePathCourses создает путь к странице курсов для указанной категории.
func MakePathCourses(categoryID string) string {
	return fmt.Sprintf("/categories/%s/courses", categoryID)
}

// MakePathCourse создает путь к странице конкретного курса.
func MakePathCourse(categoryID, courseID string) string {
	return fmt.Sprintf("%s/%s", MakePathCourses(categoryID), courseID)
}

// MakePathLesson создает путь к странице конкретного урока.
func MakePathLesson(categoryID, courseID, lessonID string) string {
	return fmt.Sprintf("%s/lessons/%s", MakePathCourse(categoryID, courseID), lessonID)
}
