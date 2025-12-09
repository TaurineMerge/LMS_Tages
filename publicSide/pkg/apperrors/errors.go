package apperrors

import "fmt"

// AppError is our custom error type for the application.
type AppError struct {
	HTTPStatus int    // HTTP Status code to return
	Code       string // Application-specific error code
	Message    string // User-facing error message
}

// Error makes AppError implement the standard error interface.
func (e *AppError) Error() string {
	return e.Message
}

// Factory functions for creating specific application errors.

func NewNotFound(resource string) error {
	return &AppError{
		HTTPStatus: 404,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

func NewInvalidUUID(resource string) error {
	return &AppError{
		HTTPStatus: 400,
		Code:       "INVALID_UUID",
		Message:    fmt.Sprintf("Invalid UUID format for %s", resource),
	}
}

func NewInvalidRequest(message string) error {
	if message == "" {
		message = "Invalid request parameters"
	}
	return &AppError{
		HTTPStatus: 400,
		Code:       "INVALID_PARAMETERS",
		Message:    message,
	}
}

func NewInternal() error {
	return &AppError{
		HTTPStatus: 500,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected internal error occurred",
	}
}
