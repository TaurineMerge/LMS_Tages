package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"adminPanel/exceptions"
	"adminPanel/models"
)

// ErrorHandlerMiddleware - централизованная обработка ошибок
func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Выполняем следующий middleware/handler
		err := c.Next()
		
		// Если есть ошибка - обрабатываем её
		if err != nil {
			log.Printf("Error occurred: %v", err)
			
			// Определяем тип ошибки
			switch e := err.(type) {
			case *exceptions.AppError:
				// Наша кастомная ошибка
				return c.Status(e.StatusCode).JSON(models.ErrorResponse{
					Error: e.Message,
					Code:  e.Code,
				})
				
			case *fiber.Error:
				// Ошибка Fiber (404, 500 и т.д.)
				return c.Status(e.Code).JSON(models.ErrorResponse{
					Error: e.Message,
					Code:  getErrorCode(e.Code),
				})
				
			default:
				// Неизвестная ошибка
				// Проверяем на common database errors
				errMsg := strings.ToLower(err.Error())
				
				switch {
				case strings.Contains(errMsg, "no rows in result set"):
					return c.Status(404).JSON(models.ErrorResponse{
						Error: "Resource not found",
						Code:  "NOT_FOUND",
					})
					
				case strings.Contains(errMsg, "duplicate key"):
					return c.Status(409).JSON(models.ErrorResponse{
						Error: "Resource already exists",
						Code:  "CONFLICT",
					})
					
				case strings.Contains(errMsg, "violates foreign key constraint"):
					return c.Status(400).JSON(models.ErrorResponse{
						Error: "Invalid reference",
						Code:  "BAD_REQUEST",
					})
					
				default:
					return c.Status(500).JSON(models.ErrorResponse{
						Error: "Internal server error",
						Code:  "INTERNAL_ERROR",
					})
				}
			}
		}
		
		return nil
	}
}

// getErrorCode преобразует HTTP код в строковый код ошибки
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
		return "CONFLICT"
	case 422:
		return "VALIDATION_ERROR"
	case 500:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}