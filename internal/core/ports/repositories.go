package ports

import (
	"context"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) error
	GetByClientID(ctx context.Context, clientID string) (*domain.Client, error)
}
