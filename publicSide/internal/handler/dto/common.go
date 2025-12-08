package dto

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// PaginationQuery represents the pagination query parameters.
type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

// --- Error Response ---

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status string      `json:"status"`
	Error  ErrorDetail `json:"error"`
}

// --- Success Response ---

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
