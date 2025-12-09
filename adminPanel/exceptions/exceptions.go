package exceptions

// AppError - базовый тип ошибки (аналог Python AppException)
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
	return NewAppError(message, 409, "CONFLICT")
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
	return NewAppError(message, 500, "INTERNAL_ERROR")
}