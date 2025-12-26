// Package domain определяет основные бизнес-сущности и модели данных,
// которые используются во всем приложении.
package domain

const (
	// UserContextKey - ключ для хранения информации о пользователе в контексте Fiber.
	UserContextKey = "user"
	// SessionTokenCookie - имя cookie, в котором хранится сессионный токен (ID Token).
	SessionTokenCookie = "access_token"
)

// UserClaims представляет информацию о пользователе, извлеченную из ID Token'а.
type UserClaims struct {
	ID       string `json:"sub"`                // Уникальный идентификатор пользователя (Subject)
	Email    string `json:"email"`              // Email пользователя
	Name     string `json:"name"`               // Полное имя пользователя
	Username string `json:"preferred_username"` // Предпочитаемое имя пользователя (логин)
}
