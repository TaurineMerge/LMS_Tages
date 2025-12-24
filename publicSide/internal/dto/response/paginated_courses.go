package response

// PaginatedCoursesData is a paginated list of courses.
type PaginatedCoursesData struct {
	Items      []CourseDTO `json:"items"`
	Pagination Pagination  `json:"pagination"`
}
