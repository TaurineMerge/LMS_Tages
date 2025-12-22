// Package middleware provides HTTP middleware functions for the Fiber application.
// This includes interceptors for logging, authentication, error handling, etc.
package middleware

import (
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// APIErrorHandler handles errors for API routes, always returning a JSON response.
func APIErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
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

	// This is an unhandled, unexpected error
	slog.Error("Unhandled API error", "error", err) // Log the unexpected error
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
		Status: response.StatusError,
		Error: response.ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected internal error occurred",
		},
	})
}
