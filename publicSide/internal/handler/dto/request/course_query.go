// Package request contains Data Transfer Objects (DTOs) for incoming client requests.
package request

// CourseQuery represents query parameters for listing courses with filters and sorting.
type CourseQuery struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Level  string `query:"level"`   // Filter by difficulty level: all, easy, medium, hard
	SortBy string `query:"sort_by"` // Sort order: updated_desc, updated_asc
}
