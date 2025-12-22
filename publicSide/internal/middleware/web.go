// Package middleware provides HTTP middleware functions for the Fiber application.
package middleware

import (
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// WebErrorHandler handles errors for web routes, rendering an HTML error page.
func WebErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		slog.Info("Handler error", "error", err)
		appErr = err.(*apperrors.AppError)
	} else {
		slog.Error("Unhandled web error", "error", err)
		appErr = apperrors.NewInternal().(*apperrors.AppError)
	}

	return c.Status(appErr.HTTPStatus).Render("pages/error", fiber.Map{
		"Title":   "Error",
		"HTTPStatus":    appErr.HTTPStatus,
		"Message": appErr.Message,
	}, "layouts/main")
}

// NoCache is a middleware that sets headers to prevent browser caching.
// Useful for development to ensure fresh assets are always served.
func NoCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.Next()
	}
}
