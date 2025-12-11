package models

// HealthResponse - ответ на health check запрос
//
// Используется для проверки работоспособности сервиса.
// Может включать статус базы данных и версию приложения.
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
}

// ErrorDetails - детали ошибки
//
// Соответствует объекту error в Swagger спецификации.
// Содержит код и сообщение об ошибке.
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse - ответ с ошибкой
//
// Соответствует error envelope в Swagger спецификации.
// Используется для возврата ошибок API.
type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorDetails `json:"error"`
}

// StatusOnly - простой ответ с статусом
//
// Используется для подтверждения успешного выполнения операции.
type StatusOnly struct {
	Status string `json:"status"`
}

// MessageResponse - ответ с текстовым сообщением
//
// Используется для возврата текстовых сообщений от API.
type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
