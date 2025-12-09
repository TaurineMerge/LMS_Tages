package models

// HealthResponse - ответ health check
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
}

// ErrorResponse - стандартный ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// MessageResponse - простой ответ с сообщением
type MessageResponse struct {
	Message string `json:"message"`
}