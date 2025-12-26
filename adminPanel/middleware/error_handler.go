package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AppError представляет пользовательскую ошибку приложения с кодом и статусом HTTP.
type AppError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
}

// Error реализует интерфейс error для AppError.
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError создает новый экземпляр AppError с заданными параметрами.
func NewAppError(message string, statusCode int, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}
}

// NotFoundError создает ошибку 404 для не найденного ресурса.
func NotFoundError(resource, identifier string) *AppError {
	message := resource + " not found"
	if identifier != "" {
		message = resource + " with id '" + identifier + "' not found"
	}
	return NewAppError(message, 404, "NOT_FOUND")
}

// ConflictError создает ошибку 409 для конфликта ресурсов.
func ConflictError(message string) *AppError {
	return NewAppError(message, 409, "ALREADY_EXISTS")
}

// ValidationError создает ошибку 422 для ошибок валидации.
func ValidationError(message string) *AppError {
	return NewAppError(message, 422, "VALIDATION_ERROR")
}

// UnauthorizedError создает ошибку 401 для неавторизованного доступа.
func UnauthorizedError(message string) *AppError {
	return NewAppError(message, 401, "UNAUTHORIZED")
}

// InternalError создает ошибку 500 для внутренних ошибок сервера.
func InternalError(message string) *AppError {
	return NewAppError(message, 500, "SERVER_ERROR")
}

// ErrorDetails содержит детали ошибки для ответа API.
type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse представляет структуру ответа с ошибкой для API.
type ErrorResponse struct {
	Status string       `json:"status"`
	Error  ErrorDetails `json:"error"`
}

// ErrorHandlerMiddleware возвращает промежуточное ПО для обработки ошибок.
// Преобразует ошибки в соответствующие HTTP-ответы для API или HTML.
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			log.Printf("Error occurred: %v", err)

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

// getErrorCode возвращает строковый код ошибки по HTTP-статус коду.
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
