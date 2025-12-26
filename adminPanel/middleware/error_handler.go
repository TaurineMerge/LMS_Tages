package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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

// ErrorDetails - детализированная информация об ошибке
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

// ErrorHandlerMiddleware - middleware для централизованной обработки ошибок
//
// Middleware перехватывает все ошибки, возникающие в ходе выполнения
// запроса, и возвращает их в стандартизированном формате.
//
// Обрабатываемые типы ошибок:
//   - AppError: кастомные ошибки приложения
//   - fiber.Error: ошибки фреймворка Fiber
//   - Database errors: ошибки базы данных (404, 409, 400)
//   - Unknown errors: неизвестные ошибки (500)
//
// Возвращает:
//   - fiber.Handler: middleware для использования в Fiber
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			log.Printf("Error occurred: %v", err)

			// Определяем тип запроса: API или веб
			isAPIRequest := strings.HasPrefix(c.Path(), "/api/")

			switch e := err.(type) {
			case *AppError:
				if isAPIRequest {
					return c.Status(e.StatusCode).JSON(ErrorResponse{
						Status: "error",
						Error: ErrorDetails{
							Code:    e.Code,
							Message: e.Message,
						},
					})
				} else {
					return c.Status(e.StatusCode).Render("pages/error", fiber.Map{
						"title":      "Ошибка",
						"HTTPStatus": e.StatusCode,
						"Message":    e.Message,
					}, "layouts/main")
				}

			case *fiber.Error:
				if isAPIRequest {
					return c.Status(e.Code).JSON(ErrorResponse{
						Status: "error",
						Error: ErrorDetails{
							Code:    getErrorCode(e.Code),
							Message: e.Message,
						},
					})
				} else {
					return c.Status(e.Code).Render("pages/error", fiber.Map{
						"title":      "Ошибка",
						"HTTPStatus": e.Code,
						"Message":    e.Message,
					}, "layouts/main")
				}

			default:
				errMsg := strings.ToLower(err.Error())

				switch {
				case strings.Contains(errMsg, "no rows in result set"):
					if isAPIRequest {
						return c.Status(404).JSON(ErrorResponse{
							Status: "error",
							Error: ErrorDetails{
								Code:    "NOT_FOUND",
								Message: "Resource not found",
							},
						})
					} else {
						return c.Status(404).Render("pages/error", fiber.Map{
							"title":      "Ошибка",
							"HTTPStatus": 404,
							"Message":    "Resource not found",
						}, "layouts/main")
					}

				case strings.Contains(errMsg, "duplicate key"):
					if isAPIRequest {
						return c.Status(409).JSON(ErrorResponse{
							Status: "error",
							Error: ErrorDetails{
								Code:    "ALREADY_EXISTS",
								Message: "Resource already exists",
							},
						})
					} else {
						return c.Status(409).Render("pages/error", fiber.Map{
							"title":      "Ошибка",
							"HTTPStatus": 409,
							"Message":    "Resource already exists",
						}, "layouts/main")
					}

				case strings.Contains(errMsg, "violates foreign key constraint"):
					if isAPIRequest {
						return c.Status(400).JSON(ErrorResponse{
							Status: "error",
							Error: ErrorDetails{
								Code:    "INVALID_REFERENCE",
								Message: "Invalid reference",
							},
						})
					} else {
						return c.Status(400).Render("pages/error", fiber.Map{
							"title":      "Ошибка",
							"HTTPStatus": 400,
							"Message":    "Invalid reference",
						}, "layouts/main")
					}

				default:
					if isAPIRequest {
						return c.Status(500).JSON(ErrorResponse{
							Status: "error",
							Error: ErrorDetails{
								Code:    "SERVER_ERROR",
								Message: "Internal server error",
							},
						})
					} else {
						return c.Status(500).Render("pages/error", fiber.Map{
							"title":      "Ошибка",
							"HTTPStatus": 500,
							"Message":    "Internal server error",
						}, "layouts/main")
					}
				}
			}
		}

		return nil
	}
}

// getErrorCode преобразует HTTP статус-код в строковый код ошибки
//
// Функция возвращает соответствующий строковый код ошибки
// для стандартных HTTP статусов.
//
// Параметры:
//   - statusCode: HTTP статус-код
//
// Возвращает:
//   - string: строковый код ошибки
func getErrorCode(statusCode int) string {
	switch statusCode {
	case 400:
		return "BAD_REQUEST"
	case 401:
		return "UNAUTHORIZED"
	case 403:
		return "FORBIDDEN"
	case 404:
		return "NOT_FOUND"
	case 409:
		return "ALREADY_EXISTS"
	case 422:
		return "VALIDATION_ERROR"
	case 500:
		return "SERVER_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}
