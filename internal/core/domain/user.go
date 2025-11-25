package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordMismatch  = errors.New("password does not match")
	ErrEmailNotVerified  = errors.New("email not verified")
)

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	PasswordHash  string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUser(name, email, passwordHash string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &User{
		ID:            id,
		Name:          name,
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (u *User) VerifyEmail() {
	u.EmailVerified = true
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) UpdatePasswordHash(newHash string) {
	u.PasswordHash = newHash
	u.UpdatedAt = time.Now().UTC()
}
