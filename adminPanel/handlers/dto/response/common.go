package response

// HealthResponse представляет ответ на health check запрос.
// Содержит статус сервиса, базы данных и версию.
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
}

// ErrorDetails содержит детали ошибки для ответа API.
// Включает код и сообщение об ошибке.
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse представляет ответ с ошибкой.
// Стандартизированный формат для ошибок API.
type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorDetails `json:"error"`
}

// StatusOnly представляет простой ответ с статусом.
// Используется для подтверждения успешных операций.
type StatusOnly struct {
	Status string `json:"status"`
}

// MessageResponse представляет ответ с текстовым сообщением.
// Содержит статус и сообщение для пользователя.
type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ValidationErrorResponse представляет ответ с ошибками валидации.
// Включает детали ошибки и карту полей с ошибками.
type ValidationErrorResponse struct {
	Status string            `json:"status"`
	Error  ErrorDetails      `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}
