package response

// PaginatedCoursesData is a paginated list of courses.
type PaginatedCoursesData struct {
	CategoryID   string      `json:"category_id"`
	CategoryName string      `json:"category_name"`
	Items        []CourseDTO `json:"items"`
	Pagination   Pagination  `json:"pagination"`
}
