package main

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

// Конфигурация аутентификации через Keycloak.
// Ожидается, что будет создан отдельный realm и клиент для adminPanel.

type AuthConfig struct {
	IssuerURL string
	Audience  string
	JWKSURL   string
}

var (
	authConfig AuthConfig
	jwks       *keyfunc.JWKS
)

// verifyAudience проверяет, что aud в клеймах содержит ожидаемую audience.
// Поддерживает варианты, которые обычно выдаёт Keycloak: строка или массив.
func verifyAudience(claims jwt.MapClaims, expected string) bool {
	if expected == "" {
		return true
	}

	// aud может быть строкой
	if audStr, ok := claims["aud"].(string); ok {
		return audStr == expected
	}

	// или слайсом (интерфейсов / строк)
	if audSlice, ok := claims["aud"].([]interface{}); ok {
		for _, v := range audSlice {
			if s, ok := v.(string); ok && s == expected {
				return true
			}
		}
	}
	if audSlice, ok := claims["aud"].([]string); ok {
		for _, s := range audSlice {
			if s == expected {
				return true
			}
		}
	}

	return false
}

// initAuth инициализирует конфигурацию и JWKS-клиент для проверки JWT.
// Ожидаемые переменные окружения:
//   KEYCLOAK_ISSUER_URL  - URL iss (например, http://keycloak:8080/realms/admin-panel)
//   KEYCLOAK_AUDIENCE    - audience / client_id
//   KEYCLOAK_JWKS_URL    - URL до JWKS (обычно: <issuer>/protocol/openid-connect/certs)
func initAuth() error {
	issuer := os.Getenv("KEYCLOAK_ISSUER_URL")
	if issuer == "" {
		return nil // аутентификация не включена, работаем без неё
	}

	aud := os.Getenv("KEYCLOAK_AUDIENCE")
	jwksURL := os.Getenv("KEYCLOAK_JWKS_URL")
	if jwksURL == "" {
		// по умолчанию строим из issuer
		jwksURL = strings.TrimRight(issuer, "/") + "/protocol/openid-connect/certs"
	}

	authConfig = AuthConfig{
		IssuerURL: issuer,
		Audience:  aud,
		JWKSURL:   jwksURL,
	}

	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("❌ Ошибка обновления JWKS: %v", err)
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

	log.Printf("✅ Инициализирована аутентификация через Keycloak (issuer=%s)", authConfig.IssuerURL)
	return nil
}

// authMiddleware проверяет JWT в заголовке Authorization: Bearer <token>
// и валидирует его против Keycloak JWKS. Если переменная KEYCLOAK_ISSUER_URL
// не задана, middleware пропускает запросы без проверки (режим без аутентификации).
func authMiddleware(c *fiber.Ctx) error {
	// Публичные эндпоинты можно пропускать без проверки
	path := c.Path()
	if strings.HasPrefix(path, "/swagger") || path == "/swagger.json" || path == "/api/v1/health" {
		return c.Next()
	}

	if authConfig.IssuerURL == "" || jwks == nil {
		// Аутентификация не настроена, пропускаем (например, локальная разработка)
		return c.Next()
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(http.StatusUnauthorized).JSON(newUnauthorizedError("Missing or invalid Authorization header"))
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(newUnauthorizedError("Empty bearer token"))
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
	if err != nil || !token.Valid {
		log.Printf("❌ Невалидный JWT: %v", err)
		return c.Status(http.StatusUnauthorized).JSON(newUnauthorizedError("Invalid token"))
	}

	// Дополнительные проверки issuer / audience, если заданы
	if iss, ok := claims["iss"].(string); ok && authConfig.IssuerURL != "" && iss != authConfig.IssuerURL {
		return c.Status(http.StatusUnauthorized).JSON(newUnauthorizedError("Invalid token issuer"))
	}

	if authConfig.Audience != "" {
		if !verifyAudience(claims, authConfig.Audience) {
			return c.Status(http.StatusUnauthorized).JSON(newUnauthorizedError("Invalid token audience"))
		}
	}

	// Сохраняем клеймы в контекст для дальнейшей авторизации (ролей и т.п.)
	c.Locals("userClaims", claims)

	return c.Next()
}
