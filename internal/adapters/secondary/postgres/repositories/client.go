package repositories

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
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
		Scopes:        client.Scopes,
		LogoUrl:       client.LogoURL,
	})

	return err
}

func (r *ClientRepository) GetByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	client, err := r.queries.GetClientByClientID(ctx, clientID)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}
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
		Scopes:        client.Scopes,
		LogoURL:       client.LogoUrl,
	}, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	client, err := r.queries.GetClientByID(ctx, pgUUID)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}
		return nil, fmt.Errorf("get client by ID: %w", err)
	}

	return &domain.Client{
		ID:            client.ID.Bytes,
		ClientID:      client.ClientID,
		ClientSecret:  client.ClientSecret,
		ClientName:    client.ClientName,
		RedirectURIs:  client.RedirectUris,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scopes:        client.Scopes,
		LogoURL:       client.LogoUrl,
		CreatedAt:     client.CreatedAt.Time,
		UpdatedAt:     client.UpdatedAt.Time,
	}, nil
}

func (r *ClientRepository) List(ctx context.Context) ([]*domain.Client, error) {
	clients, err := r.queries.ListClients(ctx)
	if err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}

	result := make([]*domain.Client, 0, len(clients))
	for _, client := range clients {
		result = append(result, &domain.Client{
			ID:            client.ID.Bytes,
			ClientID:      client.ClientID,
			ClientSecret:  client.ClientSecret,
			ClientName:    client.ClientName,
			RedirectURIs:  client.RedirectUris,
			GrantTypes:    client.GrantTypes,
			ResponseTypes: client.ResponseTypes,
			Scopes:        client.Scopes,
			LogoURL:       client.LogoUrl,
			CreatedAt:     client.CreatedAt.Time,
			UpdatedAt:     client.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (r *ClientRepository) Update(ctx context.Context, client *domain.Client) error {
	pgUUID := pgtype.UUID{
		Bytes: client.ID,
		Valid: true,
	}

	_, err := r.queries.UpdateClient(ctx, db.UpdateClientParams{
		ID:            pgUUID,
		ClientName:    client.ClientName,
		RedirectUris:  client.RedirectURIs,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scopes:        client.Scopes,
	})

	if err != nil {
		return fmt.Errorf("update client: %w", err)
	}

	return nil
}

func (r *ClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	err := r.queries.DeleteClient(ctx, pgUUID)
	if err != nil {
		return fmt.Errorf("delete client: %w", err)
	}

	return nil
}
