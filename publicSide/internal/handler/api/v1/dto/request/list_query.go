// Package request contains Data Transfer Objects (DTOs) for incoming client requests.
package request

// ListQuery represents the common query parameters for listing resources, including pagination and sorting.
type ListQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Sort  string `query:"sort"` // Field and order for sorting (e.g., "title", "-created_at")
}
