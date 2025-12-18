package response

// PaginatedCategoriesData is a paginated list of categories.
type PaginatedCategoriesData struct {
	Items      []CategoryDTO `json:"items"`
	Pagination Pagination    `json:"pagination"`
}
