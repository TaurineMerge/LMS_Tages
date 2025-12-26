// Package web содержит обработчики для рендеринга веб-страниц.
package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

// AuthHandler обрабатывает HTTP-запросы, связанные с аутентификацией через OIDC.
type AuthHandler struct {
	provider     *oidc.Provider
	oauth2Config *oauth2.Config
}

// NewAuthHandler создает новый экземпляр AuthHandler.
func NewAuthHandler(provider *oidc.Provider, oauth2Config *oauth2.Config) *AuthHandler {
	return &AuthHandler{
		provider:     provider,
		oauth2Config: oauth2Config,
	}
}

// Login инициирует процесс аутентификации OIDC.
// Он генерирует `state` для защиты от CSRF, сохраняет его в cookie
// и перенаправляет пользователя на страницу входа провайдера.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Failed to generate random state for OIDC", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to initiate login")
	}
	state := hex.EncodeToString(b)

	c.Cookie(&fiber.Cookie{
		Name:     "oidc_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	authURL := h.oauth2Config.AuthCodeURL(state)
	return c.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// Callback обрабатывает обратный вызов от OIDC провайдера после аутентификации.
// Он проверяет `state`, обменивает `code` на токены, верифицирует `id_token`
// и сохраняет его в сессионной cookie.
func (h *AuthHandler) Callback(c *fiber.Ctx) error {
	stateCookie := c.Cookies("oidc_state")
	if stateCookie == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing state cookie")
	}

	if c.Query("state") != stateCookie {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid OIDC state")
	}

	ctx := context.Background()
	tokens, err := h.oauth2Config.Exchange(ctx, c.Query("code"))
	if err != nil {
		slog.Error("Failed to exchange code for tokens", "error", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to log in")
	}

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

	c.Cookie(&fiber.Cookie{
		Name:     domain.SessionTokenCookie,
		Value:    rawIDToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})

	// Удаляем state cookie после успешного использования.
	c.Cookie(&fiber.Cookie{
		Name:     "oidc_state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}

// Logout выполняет выход пользователя из системы.
// Он удаляет сессионную cookie и перенаправляет на главную страницу.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     domain.SessionTokenCookie,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   c.Protocol() == "https",
		SameSite: "Lax",
	})
	return c.Redirect("/", fiber.StatusTemporaryRedirect)
}
