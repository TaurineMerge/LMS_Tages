package routing

import "fmt"

// --- Path Variable Names (for Fiber) ---
// Used for Fiber to extract parameters from the URL
const (
	PathVariableCategoryID = "category_id"
	PathVariableCourseID   = "course_id"
	PathVariableLessonID   = "lesson_id"
)

// --- Route Definitions (for Fiber's app.Group and app.Get patterns) ---
const (
	RouteHome = "/"
	RouteLogin = "/login"
	RouteLogout = "/logout"
	RouteAuthCallback = "/auth/callback"
	RouteReg = "/reg"

	RouteProfile = "account/profile"

	RouteAPIV1 = "/api/v1"

	RouteCategories = "/categories"
	RouteCategory = "/categories/:" + PathVariableCategoryID
	RouteCourses = "/categories/:" + PathVariableCategoryID + "/courses"
	RouteCourse = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID
	RouteLessons = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID + "/lessons"
	RouteLesson = "/categories/:" + PathVariableCategoryID + "/courses/:" + PathVariableCourseID + "/lessons/:" + PathVariableLessonID
)

// --- Path Constructors (for generating URLs in templates, redirects, etc.) ---

func MakePathHome() string {
	return RouteHome
}

func MakePathCategories() string {
	return RouteCategories
}

func MakePathCourses(categoryID string) string {
	return fmt.Sprintf("/categories/%s/courses", categoryID)
}

func MakePathCourse(categoryID, courseID string) string {
	return fmt.Sprintf("%s/%s", MakePathCourses(categoryID), courseID)
}

func MakePathLesson(categoryID, courseID, lessonID string) string {
	return fmt.Sprintf("%s/lessons/%s", MakePathCourse(categoryID, courseID), lessonID)
}
