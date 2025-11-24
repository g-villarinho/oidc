package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
)

type AuthService struct {
	userRepository ports.UserRepository
	hasher         ports.Hasher
}

func NewAuthService(userRepo ports.UserRepository, hasher ports.Hasher) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		hasher:         hasher,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) error {
	passwordHash, err := s.hasher.Hash(ctx, password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user, err := domain.NewUser(name, email, passwordHash)
	if err != nil {
		return fmt.Errorf("create new user: %w", err)
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		if err == ports.ErrAlreadyExists {
			return fmt.Errorf("persist user: %w", domain.ErrUserAlreadyExists)
		}

		return fmt.Errorf("persist user: %w", err)
	}

	return nil
}
