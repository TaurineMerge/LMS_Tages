// Package web contains web handlers and related middleware
package web

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// AuthMiddleware contains dependencies for auth middleware.
type AuthMiddleware struct {
	provider *oidc.Provider
	clientID string
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(provider *oidc.Provider, clientID string) *AuthMiddleware {
	return &AuthMiddleware{
		provider: provider,
		clientID: clientID,
	}
}

// WithUser is a "light" middleware that checks for a session and populates
// c.Locals with user info if available, but does not redirect or block.
// This is useful for templates that need to show different states for logged-in
// and anonymous users (e.g., the header).
func (m *AuthMiddleware) WithUser(c *fiber.Ctx) error {
	// Set a default empty user
	c.Locals(domain.UserContextKey, domain.UserClaims{})

	rawIDToken := c.Cookies(domain.SessionTokenCookie)
	if rawIDToken == "" {
		return c.Next() // No token, just continue
	}

	ctx := context.Background()
	verifier := m.provider.Verifier(&oidc.Config{ClientID: m.clientID})

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		// Token is invalid or expired, just continue as an anonymous user
		return c.Next()
	}

	var claims domain.UserClaims
	if err := idToken.Claims(&claims); err != nil {
		// Claims are malformed, continue as anonymous
		return c.Next()
	}

	// Token is valid, populate locals
	c.Locals(domain.UserContextKey, claims)

	return c.Next()
}
