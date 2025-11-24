package repositories

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/jackc/pgx/v5/pgtype"
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
	pgUUID := pgtype.UUID{
		Bytes: client.ID,
		Valid: true,
	}

	_, err := r.queries.CreateClient(ctx, db.CreateClientParams{
		ID:            pgUUID,
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
	client, err := r.queries.GetClientByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("get client by clientID: %w", err)
	}

	return &domain.Client{
		ID:            client.ID.Bytes,
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		ClientName:    client.ClientName,
		RedirectURIs:  client.RedirectUris,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scope:         client.Scope,
	}, nil
}
