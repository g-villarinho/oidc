package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
)

type AuthService struct {
	userService    *UserService
	tokenGenerator ports.TokenGenerator
}

func NewAuthService(userService *UserService, tokenGenerator ports.TokenGenerator) *AuthService {
	return &AuthService{
		userService:    userService,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) error {
	user, err := s.userService.CreateUser(ctx, name, email, password)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	fmt.Printf("User registered successfully: %+v\n", user)

	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*domain.LoginResponse, error) {
	user, err := s.userService.Authenticate(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(ctx, user.ID, nil, domain.AccessTokenLoginDuration)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &domain.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   domain.AccessTokenLoginDuration,
	}, nil
}
