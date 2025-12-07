package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/MicahParks/keyfunc/v2"
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
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Публичные эндпоинты
		path := c.Path()
		if strings.HasPrefix(path, "/swagger") || 
		   strings.HasPrefix(path, "/api/v1/health") || 
		   path == "/swagger.json" {
			return c.Next()
		}

		// Если аутентификация не настроена
		if authConfig == nil || jwks == nil {
			return c.Next()
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

		// Проверка issuer
		if iss, ok := claims["iss"].(string); ok && authConfig.IssuerURL != "" && iss != authConfig.IssuerURL {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token issuer",
				"code":  "UNAUTHORIZED",
			})
		}

		// Проверка audience
		if authConfig.Audience != "" && !verifyAudience(claims, authConfig.Audience) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token audience",
				"code":  "UNAUTHORIZED",
			})
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