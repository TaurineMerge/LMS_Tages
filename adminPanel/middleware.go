package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// stripAdminPrefixMiddleware убирает префикс /admin,
// так как nginx проксирует admin-panel по пути /admin/...
// Внутри приложения мы продолжаем работать с /api/v1/... маршрутом.
func stripAdminPrefixMiddleware(c *fiber.Ctx) error {
	path := c.Path()
	const prefix = "/admin"

	if strings.HasPrefix(path, prefix+"/") {
		newPath := strings.TrimPrefix(path, prefix)
		c.Path(newPath)
	}

	return c.Next()
}

// Вспомогательные функции валидации

func isValidLevel(level string) bool {
	switch strings.ToLower(level) {
	case "hard", "medium", "easy":
		return true
	default:
		return false
	}
}

func isValidVisibility(visibility string) bool {
	switch strings.ToLower(visibility) {
	case "draft", "public", "private":
		return true
	default:
		return false
	}
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
