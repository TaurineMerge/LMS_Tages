package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TrustProxyMiddleware возвращает промежуточное ПО для обработки заголовков прокси.
// Устанавливает заголовки X-Forwarded-For и X-Real-IP для корректного определения IP клиента.
func TrustProxyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
			c.Set("X-Forwarded-For", forwardedFor)
		}

		if realIP := c.Get("X-Real-IP"); realIP != "" {
			c.Set("X-Real-IP", realIP)
		}

		return c.Next()
	}
}
