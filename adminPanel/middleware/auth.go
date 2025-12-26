// Пакет middleware содержит промежуточное ПО для аутентификации,
// обработки ошибок, валидации JSON и доверия прокси.
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

// AuthConfig содержит конфигурацию для аутентификации через Keycloak.
type AuthConfig struct {
	IssuerURL string
	Audience  string
	JWKSURL   string
}

var (
	authConfig *AuthConfig
	jwks       *keyfunc.JWKS
)

// InitAuth инициализирует аутентификацию, загружая JWKS из Keycloak.
// Если переменная окружения KEYCLOAK_ISSUER_URL не установлена, аутентификация пропускается.
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
	for i := 0; i < 6; i++ {
		jwks, err = keyfunc.Get(authConfig.JWKSURL, options)
		if err == nil {
			break
		}
		log.Printf("❌ JWKS init failed (attempt %d/6): %v", i+1, err)
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		return err
	}

	log.Printf("✅ Authentication initialized (issuer=%s)", authConfig.IssuerURL)
	return nil
}

// AuthMiddleware возвращает промежуточное ПО для аутентификации JWT-токенов.
// Пропускает определенные пути без проверки и проверяет токены для остальных.
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		if strings.HasPrefix(path, "/admin/swagger") ||
			strings.HasPrefix(path, "/health") ||
			path == "/favicon.ico" ||
			path == "/admin/swagger/doc.json" {
			return c.Next()
		}

		if authConfig == nil || jwks == nil {
			log.Println("⚠️  Authentication not configured, skipping auth check")
			return c.Next()
		}

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

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)

		if err != nil || !token.Valid {
			log.Printf("❌ Invalid JWT: %v", err)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
				"code":  "UNAUTHORIZED",
			})
		}

		iss, ok := claims["iss"].(string)
		if ok && authConfig.IssuerURL != "" && iss != authConfig.IssuerURL {
			log.Printf("⚠️  Token issuer mismatch. Expected: %s, Got: %s", authConfig.IssuerURL, iss)
		}

		if authConfig.Audience != "" && !verifyAudience(claims, authConfig.Audience) {
			log.Printf("⚠️  Token audience mismatch. Expected: %s", authConfig.Audience)
		}

		if preferredUsername, ok := claims["preferred_username"].(string); ok {
			log.Printf("✅ Authenticated user: %s", preferredUsername)
		}

		c.Locals("userClaims", claims)
		return c.Next()
	}
}

// verifyAudience проверяет, соответствует ли аудитория токена ожидаемой.
// Поддерживает как строковую, так и массивную форму аудитории.
func verifyAudience(claims jwt.MapClaims, expected string) bool {
	if expected == "" {
		return true
	}

	if audStr, ok := claims["aud"].(string); ok {
		return audStr == expected
	}

	if audSlice, ok := claims["aud"].([]interface{}); ok {
		for _, v := range audSlice {
			if s, ok := v.(string); ok && s == expected {
				return true
			}
		}
	}

	return false
}
