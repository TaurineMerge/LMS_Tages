// Package apiconst provides centralized constants for the API layer.
// This includes keys for URL parameters and dynamically constructed path segments.
package apiconst

// ParamKeys holds the string keys for URL parameters.
// These should be used when parsing parameters from the request context.
const (
	PathVariableCategoryID = "category_id"
	PathVariableCourseID   = "course_id"
	PathVariableLessonID   = "lesson_id"
)