package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type AuthService struct {
	userService       *UserService
	userRepository    ports.UserRepository
	sessionRepository ports.SessionRepository
	sessionConfig     config.Session
}

func NewAuthService(userService *UserService, userRepository ports.UserRepository, sessionRepository ports.SessionRepository, config *config.Config) *AuthService {
	return &AuthService{
		userService:       userService,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		sessionConfig:     config.Session,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) error {
	user, err := s.userService.CreateUser(ctx, name, email, password)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	// TODO: Send welcome email to user confirming registration

	fmt.Printf("User registered successfully: %+v\n", user)

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.Session, *domain.User, error) {
	user, err := s.userService.Authenticate(ctx, email, password)
	if err != nil {
		return nil, nil, fmt.Errorf("login user: %w", err)
	}

	session, err := domain.NewSession(user.ID, s.sessionConfig.Duration)
	if err != nil {
		return nil, nil, fmt.Errorf("create session: %w", err)
	}

	if err := s.sessionRepository.Create(ctx, session); err != nil {
		return nil, nil, fmt.Errorf("store session: %w", err)
	}

	return session, user, nil
}

func (s *AuthService) GetSessionUser(ctx context.Context, sessionID uuid.UUID) (*domain.User, error) {
	session, err := s.sessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	if session == nil {
		return nil, domain.ErrSessionNotFound
	}

	if session.IsExpired() {
		if err := s.sessionRepository.Delete(ctx, session.ID); err != nil {
			return nil, fmt.Errorf("delete expired session: %w", err)
		}

		return nil, domain.ErrSessionExpired
	}

	user, err := s.userRepository.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user by session: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}
