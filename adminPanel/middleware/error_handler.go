package middleware

import (
	"log"
	"strings"

	"adminPanel/exceptions"
	dto "adminPanel/handlers/dto/response"

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
		err := c.Next()

		if err != nil {
			log.Printf("Error occurred: %v", err)

			// Определяем тип запроса: API или веб
			isAPIRequest := strings.HasPrefix(c.Path(), "/api/")

			switch e := err.(type) {
			case *exceptions.AppError:
				if isAPIRequest {
					return c.Status(e.StatusCode).JSON(dto.ErrorResponse{
						Status: "error",
						Error: dto.ErrorDetails{
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
					return c.Status(e.Code).JSON(dto.ErrorResponse{
						Status: "error",
						Error: dto.ErrorDetails{
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
						return c.Status(404).JSON(dto.ErrorResponse{
							Status: "error",
							Error: dto.ErrorDetails{
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
						return c.Status(409).JSON(dto.ErrorResponse{
							Status: "error",
							Error: dto.ErrorDetails{
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
						return c.Status(400).JSON(dto.ErrorResponse{
							Status: "error",
							Error: dto.ErrorDetails{
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
						return c.Status(500).JSON(dto.ErrorResponse{
							Status: "error",
							Error: dto.ErrorDetails{
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
