// Package response содержит структуры данных для формирования HTTP-ответов.
package response

// ErrorDetail содержит стандартизированную информацию об ошибке.
type ErrorDetail struct {
	Code    string `json:"code"`    // Уникальный код ошибки для программной обработки.
	Message string `json:"message"` // Человекочитаемое сообщение об ошибке.
}

// ErrorResponse представляет собой стандартную структуру JSON-ответа для ошибок.
type ErrorResponse struct {
	Status string      `json:"status"` // Статус ответа, обычно "error".
	Error  ErrorDetail `json:"error"`  // Детали ошибки.
}
