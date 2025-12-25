// Package domain contains domain models
package domain

const (
	UserContextKey      = "user"
	SessionTokenCookie  = "session_token"
	RefreshTokenCookie  = "refresh_token"
)

type UserClaims struct {
	ID       string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"preferred_username"`
}
