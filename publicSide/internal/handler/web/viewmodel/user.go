package viewmodel

type UserClaims struct {
	ID       string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"preferred_username"`
}