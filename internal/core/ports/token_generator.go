package ports

import (
	"context"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/google/uuid"
)

type TokenGenerator interface {
	GenerateAccessToken(ctx context.Context, userID uuid.UUID, clientID string, scopes []string) (string, error)
	GenerateRefreshToken(ctx context.Context) (string, error)
	GenerateIDToken(ctx context.Context, user *domain.User, clientID, nonce string, scopes []string) (string, error)
}
