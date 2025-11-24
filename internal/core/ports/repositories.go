package ports

import (
	"context"
	"errors"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrAlreadyExists = errors.New("already exists in the repository")
	ErrNotFound      = errors.New("not found in the repository")
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) error
	GetByClientID(ctx context.Context, clientID string) (*domain.Client, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
