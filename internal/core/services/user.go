package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) (*domain.User, error)
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type UserServiceImpl struct {
	userRepository ports.UserRepository
	hasher         ports.Hasher
	logger         *slog.Logger
}

func NewUserService(userRepo ports.UserRepository, hasher ports.Hasher, logger *slog.Logger) UserService {
	return &UserServiceImpl{
		userRepository: userRepo,
		hasher:         hasher,
		logger:         logger.With("service", "user"),
	}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	existingUser, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, ports.ErrNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	passwordHash, err := s.hasher.Hash(ctx, password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := domain.NewUser(name, email, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("create new user: %w", err)
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		if err == ports.ErrUniqueKeyViolation {
			return nil, domain.ErrUserAlreadyExists
		}

		return nil, fmt.Errorf("persist user: %w", err)
	}

	return user, nil
}

func (s *UserServiceImpl) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	logger := s.logger.With("method", "Authenticate")

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	if err := s.hasher.Compare(ctx, password, user.PasswordHash); err != nil {
		logger.Warn("error to compare password", slog.Any("err", err))
		return nil, domain.ErrPasswordMismatch
	}

	if !user.EmailVerified {
		return nil, domain.ErrEmailNotVerified
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}
