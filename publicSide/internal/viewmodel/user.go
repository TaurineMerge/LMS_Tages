// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// UserViewModel представляет данные о пользователе для отображения в шаблонах.
type UserViewModel struct {
	ID       string
	Email    string
	Name     string
	Username string
}

// NewUserViewModel создает новую модель представления для пользователя на основе claims из JWT.
func NewUserViewModel(claims domain.UserClaims) *UserViewModel {
	return &UserViewModel{
		ID:       claims.ID,
		Email:    claims.Email,
		Name:     claims.Name,
		Username: claims.Username,
	}
}
