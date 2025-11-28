package models

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Continue string `json:"continue" validate:"required,url"`
}

type LoginResponse struct {
	User     UserResponse `json:"user"`
	Continue string       `json:"continue"`
}

type RegisterPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,strong_password"`
}
