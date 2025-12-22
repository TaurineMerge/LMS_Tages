package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

type UserViewModel struct {
	ID       string
	Email    string
	Name     string
	Username string
}

func NewUserViewModel(claims domain.UserClaims) *UserViewModel {
	return &UserViewModel{
		ID:       claims.ID,
		Email:    claims.Email,
		Name:     claims.Name,
		Username: claims.Username,
	}
}
