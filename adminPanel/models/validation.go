package models

// ValidationErrorResponse - ответ с ошибками валидации
type ValidationErrorResponse struct {
	Status string            `json:"status"`
	Error  ErrorDetails      `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}
