// Пакет exceptions содержит типы ошибок приложения
//
// Пакет предоставляет:
//   - AppError: базовую структуру для всех ошибок
//   - Конструкторы для стандартных ошибок
//   - Единый формат ошибок для API
//
// Стандартные ошибки:
//   - NotFoundError: ресурс не найден
//   - ConflictError: конфликт данных
//   - ValidationError: ошибка валидации
//   - UnauthorizedError: ошибка авторизации
//   - InternalError: внутренняя ошибка сервера
package exceptions

// AppError - базовая структура для всех ошибок приложения
//
// Содержит:
//   - Message: текстовое описание ошибки
//   - StatusCode: HTTP статус-код
//   - Code: строковый код ошибки для программной обработки
type AppError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError создает новую ошибку приложения
//
// Параметры:
//   - message: текстовое описание ошибки
//   - statusCode: HTTP статус-код
//   - code: строковый код ошибки
//
// Возвращает:
//   - *AppError: указатель на новую ошибку
func NewAppError(message string, statusCode int, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}
}

// NotFoundError создает ошибку "Ресурс не найден"
//
// Параметры:
//   - resource: тип ресурса (например, "Course", "Category")
//   - identifier: идентификатор ресурса (опционально)
//
// Возвращает:
//   - *AppError: ошибка с кодом 404
func NotFoundError(resource, identifier string) *AppError {
	message := resource + " not found"
	if identifier != "" {
		message = resource + " with id '" + identifier + "' not found"
	}
	return NewAppError(message, 404, "NOT_FOUND")
}

// ConflictError создает ошибку "Конфликт данных"
//
// Используется при попытке создать дубликат или нарушении
// уникальности данных.
//
// Параметры:
//   - message: описание конфликта
//
// Возвращает:
//   - *AppError: ошибка с кодом 409
func ConflictError(message string) *AppError {
	return NewAppError(message, 409, "ALREADY_EXISTS")
}

// ValidationError создает ошибку "Ошибка валидации"
//
// Используется при нарушении правил валидации данных.
//
// Параметры:
//   - message: описание ошибки валидации
//
// Возвращает:
//   - *AppError: ошибка с кодом 422
func ValidationError(message string) *AppError {
	return NewAppError(message, 422, "VALIDATION_ERROR")
}

// UnauthorizedError создает ошибку "Неавторизованный доступ"
//
// Используется при отсутствии или невалидности аутентификации.
//
// Параметры:
//   - message: описание ошибки авторизации
//
// Возвращает:
//   - *AppError: ошибка с кодом 401
func UnauthorizedError(message string) *AppError {
	return NewAppError(message, 401, "UNAUTHORIZED")
}

// InternalError создает ошибку "Внутренняя ошибка сервера"
//
// Используется для неожиданных ошибок, возникающих на сервере.
//
// Параметры:
//   - message: описание внутренней ошибки
//
// Возвращает:
//   - *AppError: ошибка с кодом 500
func InternalError(message string) *AppError {
	return NewAppError(message, 500, "SERVER_ERROR")
}
