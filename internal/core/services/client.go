package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type ClientService struct {
	clientRepository ports.ClientRepository
	hasher           ports.Hasher
}

func NewClientService(clientRepository ports.ClientRepository, hasher ports.Hasher) *ClientService {
	return &ClientService{
		clientRepository: clientRepository,
		hasher:           hasher,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, params domain.CreateClientParams) (*domain.Client, error) {
	clientID := uuid.New().String()

	clientSecret, err := s.hasher.Hash(ctx, uuid.New().String())
	if err != nil {
		return nil, fmt.Errorf("hash client secret: %w", err)
	}

	client, err := domain.NewClient(clientID, clientSecret, params.ClientName, params.RedirectURIs, params.GrantTypes, params.ResponseTypes, params.Scopes, params.LogoURL)
	if err != nil {
		return nil, fmt.Errorf("create client domain: %w", err)
	}

	if err := s.clientRepository.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return client, nil
}

func (s *ClientService) GetClientByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	client, err := s.clientRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get client by ID: %w", err)
	}

	return client, nil
}

func (s *ClientService) ListClients(ctx context.Context) ([]*domain.Client, error) {
	clients, err := s.clientRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list clients: %w", err)
	}

	return clients, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, id uuid.UUID, clientName string, redirectURIs, grantTypes, responseTypes, scopes []string) (*domain.Client, error) {
	client, err := s.clientRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get client for update: %w", err)
	}

	client.ClientName = clientName
	client.RedirectURIs = redirectURIs
	client.GrantTypes = grantTypes
	client.ResponseTypes = responseTypes
	client.Scopes = scopes

	if err := s.clientRepository.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("update client: %w", err)
	}

	updatedClient, err := s.clientRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get updated client: %w", err)
	}

	return updatedClient, nil
}

func (s *ClientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	if err := s.clientRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete client: %w", err)
	}

	return nil
}
