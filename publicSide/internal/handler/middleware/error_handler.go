package middleware

import (
	"errors"
	"log"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// GlobalErrorHandler is a centralized error handler for the application.
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// This is a custom application error
		return c.Status(appErr.HTTPStatus).JSON(dto.ErrorResponse{
			Status: dto.StatusError,
			Error: dto.ErrorDetail{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}

	// This is an unhandled, unexpected error
	log.Printf("Unhandled error: %v", err) // Log the unexpected error
	return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
		Status: dto.StatusError,
		Error: dto.ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected internal error occurred",
		},
	})
}
