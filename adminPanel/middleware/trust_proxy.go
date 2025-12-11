// Package middleware содержит middleware-функции для обработки HTTP-запросов.
// Этот пакет предоставляет функции для работы с прокси, валидации, аутентификации и обработки ошибок.
package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TrustProxyMiddleware доверяет заголовкам от прокси-сервера.
// Этот middleware устанавливает заголовки X-Forwarded-For и X-Real-IP,
// которые могут быть установлены прокси-сервером (например, nginx).
//
// Особенности:
//   - Сохраняет оригинальные IP-адреса клиентов при работе через прокси
//   - Устанавливает заголовки, которые могут быть использованы другими middleware
//   - Должен быть установлен до других middleware, использующих IP-адреса
//
// Возвращает:
//   - fiber.Handler: middleware-функцию для использования в Fiber приложении
func TrustProxyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Устанавливаем заголовок X-Forwarded-For, если он пришел от прокси
		if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {
			c.Set("X-Forwarded-For", forwardedFor)
		}

		// Устанавливаем заголовок X-Real-IP, если он пришел от прокси
		if realIP := c.Get("X-Real-IP"); realIP != "" {
			c.Set("X-Real-IP", realIP)
		}

		// Передаем управление следующему middleware
		return c.Next()
	}
}
