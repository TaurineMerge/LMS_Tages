// Package web содержит обработчики для рендеринга веб-страниц.
package web

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// AuthMiddleware предоставляет middleware для аутентификации.
type AuthMiddleware struct {
	provider *oidc.Provider
	clientID string
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(provider *oidc.Provider, clientID string) *AuthMiddleware {
	return &AuthMiddleware{
		provider: provider,
		clientID: clientID,
	}
}

// WithUser является middleware, которое проверяет сессионную cookie,
// верифицирует JWT (ID Token) и помещает информацию о пользователе (claims)
// в `c.Locals` для дальнейшего использования в обработчиках и шаблонах.
// Если токен отсутствует или невалиден, он просто передает управление дальше,
// оставляя в `c.Locals` пустую структуру UserClaims (гостевой пользователь).
func (m *AuthMiddleware) WithUser(c *fiber.Ctx) error {
	// Инициализируем пустыми данными на случай, если пользователь гость.
	c.Locals(domain.UserContextKey, domain.UserClaims{})

	rawIDToken := c.Cookies(domain.SessionTokenCookie)
	if rawIDToken == "" {
		return c.Next()
	}

	ctx := context.Background()
	verifier := m.provider.Verifier(&oidc.Config{ClientID: m.clientID})

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// Если токен невалиден (например, истек), считаем пользователя гостем.
		return c.Next()
	}

	var claims domain.UserClaims
	if err := idToken.Claims(&claims); err != nil {
		// Если не удалось извлечь claims, считаем пользователя гостем.
		return c.Next()
	}

	// Сохраняем claims в контексте для доступа в последующих обработчиках.
	c.Locals(domain.UserContextKey, claims)

	return c.Next()
}
