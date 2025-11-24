package postgres

import (
	"context"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) ports.ClientRepository {
	return &ClientRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *ClientRepository) Create(ctx context.Context, client *domain.Client) error {
	_, err := r.queries.CreateClient(ctx, db.CreateClientParams{
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		ClientName:    client.ClientName,
		RedirectUris:  client.RedirectURIs,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scope:         client.Scope,
	})

	return err
}

func (r *ClientRepository) GetByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	panic("unimplemented")
}
