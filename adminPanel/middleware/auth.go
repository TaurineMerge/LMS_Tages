package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthConfig - конфигурация аутентификации
type AuthConfig struct {
	IssuerURL string
	Audience  string
	JWKSURL   string
}

var (
	authConfig *AuthConfig
	jwks       *keyfunc.JWKS
)

// InitAuth инициализирует аутентификацию
func InitAuth() error {
	issuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	if issuer == "" {
		log.Println("ℹ️  Authentication is not configured")
		return nil
	}

	audience := os.Getenv("KEYCLOAK_AUDIENCE")
	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL")
	if jwksURL == "" {
		jwksURL = strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
	}

	authConfig = &AuthConfig{
		IssuerURL: issuer,
		Audience:  audience,
		JWKSURL:   jwksURL,
	}

	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("❌ JWKS refresh error: %v", err)
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	}

	var err error
	jwks, err = keyfunc.Get(authConfig.JWKSURL, options)
	if err != nil {
		return err
	}

	log.Printf("✅ Authentication initialized (issuer=%s)", authConfig.IssuerURL)
	return nil
}

// AuthMiddleware - middleware для проверки JWT
// AuthMiddleware - middleware для проверки JWT
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Публичные эндпоинты
		path := c.Path()

		// Разрешаем доступ без авторизации к этим маршрутам
		if strings.HasPrefix(path, "/admin/swagger") ||
			strings.HasPrefix(path, "/health") ||
			path == "/favicon.ico" ||
			path == "/admin/swagger/doc.json" {
			return c.Next()
		}

		// Если аутентификация не настроена — блокируем запросы, чтобы не открыть API в проде
		if authConfig == nil || jwks == nil {
			log.Println("❌ Authentication not configured: KEYCLOAK_* env vars are missing")
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Authentication not configured",
				"code":  "AUTH_MISCONFIGURED",
			})
		}

		// Проверка заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid Authorization header",
				"code":  "UNAUTHORIZED",
			})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Empty bearer token",
				"code":  "UNAUTHORIZED",
			})
		}

		// Парсинг и валидация токена
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)

		if err != nil || !token.Valid {
			log.Printf("❌ Invalid JWT: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
				"code":  "UNAUTHORIZED",
			})
		}

		// Проверка issuer (можно закомментировать для разработки)
		iss, ok := claims["iss"].(string)
		if ok && authConfig.IssuerURL != "" && iss != authConfig.IssuerURL {
			log.Printf("⚠️  Token issuer mismatch. Expected: %s, Got: %s", authConfig.IssuerURL, iss)
		}

		// Проверка audience (опционально)
		if authConfig.Audience != "" && !verifyAudience(claims, authConfig.Audience) {
			log.Printf("⚠️  Token audience mismatch. Expected: %s", authConfig.Audience)
		}

		// Логируем информацию о пользователе
		if preferredUsername, ok := claims["preferred_username"].(string); ok {
			log.Printf("✅ Authenticated user: %s", preferredUsername)
		}

		// Сохраняем claims в контекст
		c.Locals("userClaims", claims)
		return c.Next()
	}
}

// verifyAudience проверяет audience в claims
func verifyAudience(claims jwt.MapClaims, expected string) bool {
	if expected == "" {
		return true
	}

	// aud может быть строкой
	if audStr, ok := claims["aud"].(string); ok {
		return audStr == expected
	}

	// или массивом
	if audSlice, ok := claims["aud"].([]interface{}); ok {
		for _, v := range audSlice {
			if s, ok := v.(string); ok && s == expected {
				return true
			}
		}
	}

	return false
}
