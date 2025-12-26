// Package testing предоставляет клиент для взаимодействия с внешним сервисом тестирования.
package testing

import "errors"

var (
	// ErrTestNotFound возникает, когда сервис тестирования возвращает статус "not_found".
	ErrTestNotFound = errors.New("test not found")

	// ErrServiceUnavailable возникает при сетевых ошибках или таймаутах при обращении к сервису.
	ErrServiceUnavailable = errors.New("testing service is unavailable")

	// ErrInvalidResponse возникает, если ответ от сервиса не соответствует JSON-схеме или не может быть разобран.
	ErrInvalidResponse = errors.New("invalid response from testing service")
)
