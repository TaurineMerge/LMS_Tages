package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// StripPrefixMiddleware удаляет префикс /admin из пути
func StripPrefixMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		const prefix = "/admin"

		if strings.HasPrefix(path, prefix+"/") {
			newPath := strings.TrimPrefix(path, prefix)
			c.Path(newPath)
		}

		return c.Next()
	}
}