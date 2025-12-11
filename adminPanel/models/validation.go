// Package models содержит модели данных для API.
// Этот пакет предоставляет структуры для запросов, ответов и валидации.
package models

// ValidationErrorResponse представляет ответ с ошибками валидации.
// Используется для возврата клиенту информации о том, какие поля не прошли валидацию.
//
// Поля:
//   - Status: статус ответа (обычно "error")
//   - Error: детали общей ошибки валидации
//   - Errors: карта с ошибками валидации для конкретных полей
//     Ключ - имя поля, Значение - сообщение об ошибке
//
// Пример использования:
//
//	{
//	  "status": "error",
//	  "error": {
//	    "code": "VALIDATION_ERROR",
//	    "message": "Validation failed"
//	  },
//	  "errors": {
//	    "title": "Title is required",
//	    "email": "Invalid email format"
//	  }
//	}
type ValidationErrorResponse struct {
	Status string            `json:"status"`
	Error  ErrorDetails      `json:"error"`
	Errors map[string]string `json:"errors,omitempty"`
}
