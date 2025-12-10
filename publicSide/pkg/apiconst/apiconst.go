// Package apiconst provides centralized constants for the API layer.
// This includes keys for URL parameters and dynamically constructed path segments.
package apiconst

// ParamKeys holds the string keys for URL parameters.
// These should be used when parsing parameters from the request context.
const (
	ParamCategoryID = "category_id"
	ParamCourseID   = "course_id"
	ParamLessonID   = "lesson_id"
)

// PathSegments holds the dynamic parts of URL paths, including the colon prefix.
// These should be used when defining router groups and endpoints.
const (
	PathCategory = "/:" + ParamCategoryID
	PathCourse   = "/:" + ParamCourseID
	PathLesson   = "/:" + ParamLessonID
)
