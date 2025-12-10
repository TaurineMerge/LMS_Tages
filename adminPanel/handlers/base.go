package handlers

import (
	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// isValidUUID проверяет валидность UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// getBaseURL возвращает базовый URL из контекста
func getBaseURL(c *fiber.Ctx) string {
	protocol := "http"
	if c.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}

	host := c.Get("Host")
	if forwardedHost := c.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return protocol + "://" + host
}
