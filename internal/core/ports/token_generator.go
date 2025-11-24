package ports

import (
	"context"
	"time"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateAccessToken(ctx context.Context, userID uuid.UUID, client *domain.Client, ttl time.Duration) (string, error)
	GenerateRefreshToken(ctx context.Context) (string, error)
	GenerateIDToken(ctx context.Context, user *domain.User, client *domain.Client, nonce string, ttl time.Duration) (string, error)
}
