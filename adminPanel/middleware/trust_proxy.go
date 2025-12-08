package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TrustProxyMiddleware доверяет заголовкам от прокси
func TrustProxyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Forwarded-For", c.Get("X-Forwarded-For"))
		c.Set("X-Real-IP", c.Get("X-Real-IP"))

		return c.Next()
	}
}
