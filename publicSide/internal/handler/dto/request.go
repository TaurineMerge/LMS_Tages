package dto

// PaginationQuery represents the pagination query parameters.
type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
