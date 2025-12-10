package models

// HealthResponse - ответ health check
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
}

// ErrorDetails соответствует swagger error object.
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse соответствует swagger error envelope.
type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorDetails `json:"error"`
}

// StatusOnly - простой успех со статусом
type StatusOnly struct {
	Status string `json:"status"`
}

// MessageResponse - вспомогательный ответ с текстом
type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
