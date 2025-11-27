package models

type LoginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Continue string `json:"continue" binding`
}

type LoginResponse struct {
	User     UserResponse `json:"user"`
	Continue string       `json:"continue"`
}

type RegisterPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8,strong_password"`
}
