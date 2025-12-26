// Package middleware предоставляет промежуточные обработчики для Fiber.
package middleware

import (
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// APIErrorHandler является обработчиком ошибок для маршрутов API.
// Он перехватывает ошибки, преобразует их в стандартизированный JSON-формат
// и отправляет клиенту с соответствующим HTTP-статусом.
func APIErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	// Пытаемся преобразовать ошибку в наш кастомный тип AppError.
	if errors.As(err, &appErr) {
		slog.Info("Handler error", "error", err)
		return c.Status(appErr.HTTPStatus).JSON(response.ErrorResponse{
			Status: response.StatusError,
			Error: response.ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}

	// Если это не AppError, считаем ее непредвиденной внутренней ошибкой.
	slog.Error("Unhandled API error", "error", err)
	appErr = apperrors.NewInternal().(*apperrors.AppError)
	return c.Status(appErr.HTTPStatus).JSON(response.ErrorResponse{
		Status: response.StatusError,
		Error: response.ErrorDetail{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
