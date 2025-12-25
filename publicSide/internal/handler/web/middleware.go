// Package web contains web handlers and related middleware
package web

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2" // Added import

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// AuthMiddleware contains dependencies for auth middleware.
type AuthMiddleware struct {
	provider     *oidc.Provider
	clientID     string
	oauth2Config *oauth2.Config // Added oauth2Config
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(provider *oidc.Provider, clientID string, oauth2Config *oauth2.Config) *AuthMiddleware {
	return &AuthMiddleware{
		provider:     provider,
		clientID:     clientID,
		oauth2Config: oauth2Config, // Initialize new field
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
		// If token is expired, try to refresh it
		if !errors.Is(err, oidc.ErrTokenExpired) {
			slog.Debug("Failed to verify token, but not because it's expired", "error", err)
			clearAuthCookies(c) // Clear invalid tokens
			return c.Next()
		}

		slog.Info("ID token expired, attempting refresh")
		refreshToken := c.Cookies(domain.RefreshTokenCookie)
		if refreshToken == "" {
			slog.Debug("No refresh token found, treating user as anonymous")
			clearAuthCookies(c) // No refresh token, so clear session token
			return c.Next()
		}

		// Perform token refresh
		ts := m.oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
		newTokens, err := ts.Token()
		if err != nil {
			slog.Debug("Failed to refresh token", "error", err)
			clearAuthCookies(c) // Refresh failed, clear all tokens
			return c.Next()
		}

		// Extract new ID token and re-verify
		newRawIDToken, ok := newTokens.Extra("id_token").(string)
		if !ok {
			slog.Error("id_token missing from refresh response")
			clearAuthCookies(c)
			return c.Next()
		}

		idToken, err = verifier.Verify(ctx, newRawIDToken)
		if err != nil {
			slog.Error("Failed to verify newly refreshed ID token", "error", err)
			clearAuthCookies(c)
			return c.Next()
		}

		// Update cookies with new tokens
		setAuthCookies(c, newRawIDToken, newTokens.RefreshToken, idToken.Expiry)
		slog.Info("Successfully refreshed tokens")
	}

	var claims domain.UserClaims
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to extract claims from token", "error", err)
		clearAuthCookies(c)
		return c.Next()
	}

	// Token is valid, populate locals
	c.Locals(domain.UserContextKey, claims)

	return c.Next()
}

// setAuthCookies is a helper to set both session and refresh cookies.
func setAuthCookies(c *fiber.Ctx, idToken, refreshToken string, expiry time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     domain.SessionTokenCookie,
		Value:    idToken,
		Expires:  expiry,
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	if refreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     domain.RefreshTokenCookie,
			Value:    refreshToken,
			Expires:  time.Now().Add(30 * 24 * time.Hour), // Set a long life for the refresh token
			HTTPOnly: true,
			Secure:   c.Protocol() == "https",
			SameSite: "Lax",
		})
	}
}

// clearAuthCookies is a helper to clear both session and refresh cookies.
func clearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     domain.SessionTokenCookie,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     domain.RefreshTokenCookie,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
}
