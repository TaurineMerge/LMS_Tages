// Package response contains Data Transfer Objects (DTOs) for outgoing server responses.
package response

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// SuccessResponse represents a success response.
type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
