package models

// ValidationErrorResponse - ответ с ошибками валидации
type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Code   string            `json:"code"`
	Errors map[string]string `json:"errors"`
}