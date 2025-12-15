// Package middleware provides HTTP middleware functions for the Fiber application.
// This includes interceptors for logging, authentication, error handling, etc.
package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RequestResponseLogger logs request and response details.
// It adds logs as events to the current OpenTelemetry span and also logs to slog at DEBUG level.
// WARNING: This can log sensitive information like request/response bodies and is not
// recommended for a production environment without proper data masking.
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the current span from the context provided by otelfiber middleware.
		span := trace.SpanFromContext(c.UserContext())

		// Log request details to console for local debugging.
		slog.Debug("Incoming request",
			"method", c.Method(),
			"path", c.Path(),
			"body", string(c.Body()),
		)

		// Add request body as a span event.
		span.AddEvent("Incoming request", trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("http.request.body", string(c.Body())),
		))

		// Proceed to the next middleware/handler.
		err := c.Next()

		// Log response details to console.
		slog.Debug("Outgoing response",
			"status", c.Response().StatusCode(),
			"body", string(c.Response().Body()),
		)

		// Add response body as a span event.
		span.AddEvent("Outgoing response", trace.WithAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
			attribute.String("http.response.body", string(c.Response().Body())),
		))

		return err
	}
}
