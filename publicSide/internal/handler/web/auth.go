// Package web contains web handlers for rendering HTML pages.
package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// AuthHandler handles OIDC authentication flows.
type AuthHandler struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(provider *oidc.Provider, oauth2Config *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		provider:     provider,
		oauth2Config: oauth2Config,
	}
}

// Login initiates the OIDC login flow.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Generate a random state for CSRF protection
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Failed to generate random state for OIDC", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to initiate login")
	}
	state := hex.EncodeToString(b)

	// Save the state in a short-lived cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	// Redirect to the Keycloak login page
	authURL := h.oauth2Config.AuthCodeURL(state)
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// Callback handles the OIDC callback from Keycloak.
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	// 1. Check for the state cookie
	stateCookie := c.Cookies("oidc_state")
	if stateCookie == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing state cookie")
	}

	// 2. Compare the state from the cookie with the one from the query params
	if c.Query("state") != stateCookie {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid OIDC state")
	}

	// 3. Exchange the authorization code for tokens
	ctx := context.Background()
	tokens, err := h.oauth2Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		slog.Error("Failed to exchange code for tokens", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to log in")
	}

	// 4. Extract and verify the ID Token
	rawIDToken, ok := tokens.Extra("id_token").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "ID token missing")
	}

	verifier := h.provider.Verifier(&oidc.Config{ClientID: h.oauth2Config.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		slog.Error("Failed to verify ID token", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid session token")
	}

	// 5. (Optional) Extract claims if needed to store in a local session
	var claims struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Username string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("Failed to extract claims from ID token", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not process user information")
	}
	slog.Info("User logged in successfully", "user", claims.Username, "email", claims.Email)

	// 6. Set the ID token and Refresh token in secure, long-lived session cookies
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    rawIDToken,
		Expires:  idToken.Expiry,
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	// Store refresh token
	if tokens.RefreshToken != "" {
		c.Cookie(&fiber.Cookie{
			Name:     domain.RefreshTokenCookie,
			Value:    tokens.RefreshToken,
			Expires:  time.Now().Add(5 * 24 * time.Hour),
			HTTPOnly: true,
			Secure:   c.Protocol() == "https",
			SameSite: "Lax",
		})
	}
	
	// 7. Clean up the state cookie
	c.Cookie(&fiber.Cookie{
		Name:     "oidc_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	// 8. Redirect to the home page
	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}

// Logout clears the session cookie and redirects to the home page.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})
	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}
