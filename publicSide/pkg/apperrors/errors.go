// Package apperrors provides custom error types and factories for the application.
// This allows for consistent, structured error responses in the API layer.
package apperrors

import "fmt"

// AppError is our custom error type for the application.
// It includes an HTTP status code, an application-specific error code,
// and a user-facing message.
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

// NewNotFound creates a new 404 Not Found error.
func NewNotFound(resource string) error {
	return &AppError{
		HTTPStatus: 404,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

// NewInvalidUUID creates a new 400 Bad Request error for invalid UUIDs.
func NewInvalidUUID(resource string) error {
	return &AppError{
		HTTPStatus: 400,
		Code:       "INVALID_UUID",
		Message:    fmt.Sprintf("Invalid UUID format for %s", resource),
	}
}

// NewInvalidRequest creates a new 400 Bad Request error for general invalid requests.
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

// NewInternal creates a new 500 Internal Server Error.
// This should be used for unexpected errors that cannot be handled.
func NewInternal() error {
	return &AppError{
		HTTPStatus: 500,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected internal error occurred",
	}
}

// NewDatabaseError creates a new 500 Internal Server Error for database operations.
// The actual error is logged but a generic message is returned to the user.
func NewDatabaseError(message string, err error) error {
	// In a production environment, you would log the actual error here
	// For now, we'll just return a generic error to the user
	return &AppError{
		HTTPStatus: 500,
		Code:       "DATABASE_ERROR",
		Message:    message,
	}
}

// NewNotFoundError creates a new 404 Not Found error with a custom message.
func NewNotFoundError(message string) error {
	return &AppError{
		HTTPStatus: 404,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}
