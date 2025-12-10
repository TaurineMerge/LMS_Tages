// Package middleware provides HTTP middleware functions for the Fiber application.
// This includes interceptors for logging, authentication, error handling, etc.
package middleware

import (
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// GlobalErrorHandler is a centralized error handler for the application.
// It catches errors returned from handlers, logs them, and formats them
// into a consistent JSON error response.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// This is a custom application error
		return c.Status(appErr.HTTPStatus).JSON(response.ErrorResponse{
			Status: response.StatusError,
			Error: response.ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}

	// This is an unhandled, unexpected error
	slog.Error("Unhandled error", "error", err) // Log the unexpected error
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
		Status: response.StatusError,
		Error: response.ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected internal error occurred",
		},
	})
}
