package models

// HealthResponse - ответ на health check запрос
//
// Используется для проверки работоспособности сервиса.
// Может включать статус базы данных и версию приложения.
//
// Поля:
//   - Status: статус работоспособности сервиса (обычно "ok" или "healthy")
//   - Database: статус подключения к базе данных (опционально, "connected" или "disconnected")
//   - Version: версия приложения (обычно semver, например "1.0.0")
type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
}

// ErrorDetails - детали ошибки
//
// Соответствует объекту error в Swagger спецификации.
// Содержит код и сообщение об ошибке для стандартизированной обработки на клиенте.
//
// Поля:
//   - Code: строковый код ошибки для программной обработки (например, "NOT_FOUND", "VALIDATION_ERROR")
//   - Message: понятное человеку описание ошибки
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse - ответ с ошибкой
//
// Соответствует error envelope в Swagger спецификации.
// Используется для возврата ошибок API в стандартизированном формате.
//
// Поля:
//   - Status: статус ответа (всегда "error")
//   - Error: детализированная информация об ошибке
type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorDetails `json:"error"`
}

// StatusOnly - простой ответ с статусом
//
// Используется для подтверждения успешного выполнения операции
// без необходимости возврата дополнительных данных.
//
// Поля:
//   - Status: статус выполнения операции (обычно "success")
type StatusOnly struct {
	Status string `json:"status"`
}

// MessageResponse - ответ с текстовым сообщением
//
// Используется для возврата текстовых сообщений от API,
// например, подтверждений удаления или информационных сообщений.
//
// Поля:
//   - Status: статус ответа (обычно "success")
//   - Message: текстовое сообщение для пользователя
type MessageResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
