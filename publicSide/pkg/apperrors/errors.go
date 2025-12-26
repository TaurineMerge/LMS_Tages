// Package apperrors предоставляет стандартизированные типы ошибок, используемые во всем приложении.
// Это позволяет последовательно обрабатывать ошибки и преобразовывать их в соответствующие HTTP-ответы.
package apperrors

import "fmt"

// AppError представляет собой стандартную ошибку приложения с дополнительной информацией
// для преобразования в HTTP-ответ.
type AppError struct {
	HTTPStatus int    // HTTP-статус, который должен быть возвращен клиенту.
	Code       string // Уникальный код ошибки для программной обработки.
	Message    string // Человекочитаемое сообщение об ошибке.
}

// Error реализует стандартный интерфейс error.
func (e *AppError) Error() string {
	return e.Message
}

// ServiceUnavailableError представляет ошибку, возникающую, когда внешний сервис недоступен.
type ServiceUnavailableError struct {
	ServiceName string // Имя недоступного сервиса.
}

// Error реализует стандартный интерфейс error.
func (e *ServiceUnavailableError) Error() string {
	return fmt.Sprintf("service %s is unavailable", e.ServiceName)
}

// NewNotFound создает новую ошибку AppError для случаев, когда ресурс не найден (HTTP 404).
func NewNotFound(resource string) error {
	return &AppError{
		HTTPStatus: 404,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s not found", resource),
	}
}

// NewInvalidUUID создает новую ошибку AppError для неверного формата UUID (HTTP 400).
func NewInvalidUUID(resource string) error {
	return &AppError{
		HTTPStatus: 400,
		Code:       "INVALID_UUID",
		Message:    fmt.Sprintf("Invalid UUID format for %s", resource),
	}
}

// NewInvalidRequest создает новую ошибку AppError для неверных параметров запроса (HTTP 400).
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

// NewInternal создает новую ошибку AppError для непредвиденных внутренних ошибок сервера (HTTP 500).
func NewInternal() error {
	return &AppError{
		HTTPStatus: 500,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected internal error occurred",
	}
}

// NewServiceUnavailable создает новую ошибку ServiceUnavailableError.
func NewServiceUnavailable(serviceName string) error {
	return &ServiceUnavailableError{ServiceName: serviceName}
}
