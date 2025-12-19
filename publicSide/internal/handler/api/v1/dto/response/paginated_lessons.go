package response

// PaginatedLessonsData is a paginated list of lessons.
type PaginatedLessonsData struct {
	Items      []LessonDTO `json:"items"`
	Pagination Pagination  `json:"pagination"`
}
