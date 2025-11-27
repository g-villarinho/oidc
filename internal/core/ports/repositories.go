package ports

import (
	"context"
	"errors"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/google/uuid"
)

var (
	ErrUniqueKeyViolation = errors.New("already exists in the repository")
	ErrNotFound           = errors.New("not found in the repository")
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	GetByClientID(ctx context.Context, clientID string) (*domain.Client, error)
	List(ctx context.Context) ([]*domain.Client, error)
	Update(ctx context.Context, client *domain.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	Delete(ctx context.Context, sessionID uuid.UUID) error
}

type AuthorizationCodeRepository interface {
	Create(ctx context.Context, code *domain.AuthorizationCode) error
	GetByCode(ctx context.Context, code string) (*domain.AuthorizationCode, error)
	Delete(ctx context.Context, code string) error
}
