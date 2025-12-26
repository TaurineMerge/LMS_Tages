// Package middleware предоставляет промежуточные обработчики для Fiber.
package middleware

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RequestResponseLogger логирует информацию о входящем запросе и исходящем ответе.
// Он также добавляет события в текущий спан OpenTelemetry для трассировки.
func RequestResponseLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		span := trace.SpanFromContext(c.UserContext())

		slog.Debug("Incoming request",
			"method", c.Method(),
			"path", c.Path(),
			"body", string(c.Body()),
		)

		span.AddEvent("Incoming request", trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("http.request.body", string(c.Body())),
		))

		err := c.Next()

		slog.Debug("Outgoing response",
			"status", c.Response().StatusCode(),
			"body", string(c.Response().Body()),
		)

		span.AddEvent("Outgoing response", trace.WithAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
			attribute.String("http.response.body", string(c.Response().Body())),
		))

		return err
	}
}

// CommonErrorHandler является единым обработчиком ошибок, который делегирует
// обработку конкретному обработчику в зависимости от пути запроса.
// Если путь начинается с "/api", используется `APIErrorHandler`, иначе `WebErrorHandler`.
func CommonErrorHandler(c *fiber.Ctx, err error) error {
	if strings.HasPrefix(c.Path(), "/api") {
		return APIErrorHandler(c, err)
	}
	return WebErrorHandler(c, err)
}
