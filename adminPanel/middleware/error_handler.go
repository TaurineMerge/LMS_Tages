package middleware

import (
	"log"
	"strings"

	"adminPanel/exceptions"
	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
)

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
		// Выполняем следующий middleware/handler
		err := c.Next()

		// Если есть ошибка - обрабатываем её
		if err != nil {
			log.Printf("Error occurred: %v", err)

			switch e := err.(type) {
			case *exceptions.AppError:
				return c.Status(e.StatusCode).JSON(models.ErrorResponse{
					Status: "error",
					Error: models.ErrorDetails{
						Code:    e.Code,
						Message: e.Message,
					},
				})

			case *fiber.Error:
				return c.Status(e.Code).JSON(models.ErrorResponse{
					Status: "error",
					Error: models.ErrorDetails{
						Code:    getErrorCode(e.Code),
						Message: e.Message,
					},
				})

			default:
				errMsg := strings.ToLower(err.Error())

				switch {
				case strings.Contains(errMsg, "no rows in result set"):
					return c.Status(404).JSON(models.ErrorResponse{
						Status: "error",
						Error: models.ErrorDetails{
							Code:    "NOT_FOUND",
							Message: "Resource not found",
						},
					})

				case strings.Contains(errMsg, "duplicate key"):
					return c.Status(409).JSON(models.ErrorResponse{
						Status: "error",
						Error: models.ErrorDetails{
							Code:    "ALREADY_EXISTS",
							Message: "Resource already exists",
						},
					})

				case strings.Contains(errMsg, "violates foreign key constraint"):
					return c.Status(400).JSON(models.ErrorResponse{
						Status: "error",
						Error: models.ErrorDetails{
							Code:    "INVALID_REFERENCE",
							Message: "Invalid reference",
						},
					})

				default:
					return c.Status(500).JSON(models.ErrorResponse{
						Status: "error",
						Error: models.ErrorDetails{
							Code:    "SERVER_ERROR",
							Message: "Internal server error",
						},
					})
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
