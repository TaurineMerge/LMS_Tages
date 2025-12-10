// Package exceptions contains reusable application error types.
package exceptions

// AppError represents an error with a message, status code, and error code.
// @property {string} Message - The `Message` property in the `AppError` struct represents the error
// message or description associated with the error that occurred. It provides information about what
// went wrong in a human-readable format.
// @property {int} StatusCode - The `StatusCode` property in the `AppError` struct represents the HTTP
// status code associated with the error. It is used to indicate the specific type of error that
// occurred during the execution of the application.
// @property {string} Code - The `Code` property in the `AppError` struct is used to store a specific
// error code that can be used to identify different types of errors. This code can be helpful for
// debugging and handling errors in a more structured way.
type AppError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError создает новую ошибку
func NewAppError(message string, statusCode int, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}
}

// NotFoundError - ресурс не найден (аналог NotFoundError)
func NotFoundError(resource, identifier string) *AppError {
	message := resource + " not found"
	if identifier != "" {
		message = resource + " with id '" + identifier + "' not found"
	}
	return NewAppError(message, 404, "NOT_FOUND")
}

// ConflictError - конфликт (аналог ConflictError)
func ConflictError(message string) *AppError {
	return NewAppError(message, 409, "ALREADY_EXISTS")
}

// ValidationError - ошибка валидации (аналог ValidationError)
func ValidationError(message string) *AppError {
	return NewAppError(message, 422, "VALIDATION_ERROR")
}

// UnauthorizedError - ошибка авторизации
func UnauthorizedError(message string) *AppError {
	return NewAppError(message, 401, "UNAUTHORIZED")
}

// InternalError - внутренняя ошибка сервера
func InternalError(message string) *AppError {
	return NewAppError(message, 500, "SERVER_ERROR")
}
