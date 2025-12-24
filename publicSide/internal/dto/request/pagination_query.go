// Package request contains Data Transfer Objects (DTOs) for incoming client requests.
package request

// PaginationQuery represents the pagination query parameters.
type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}
