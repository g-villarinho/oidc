package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type AuthorizationService struct {
	clientRepository ports.ClientRepository
}

func NewAuthorizationService(clientRepository ports.ClientRepository) *AuthorizationService {
	return &AuthorizationService{
		clientRepository: clientRepository,
	}
}

func (s *AuthorizationService) ValidateAuthorizationClient(ctx context.Context, params domain.AuthorizeParams) (*domain.Client, error) {
	client, err := s.clientRepository.GetByClientID(ctx, params.ClientID)
	if err != nil {
		return nil, fmt.Errorf("validate authorization client: %w", err)
	}

	if client == nil {
		return nil, domain.ErrClientNotFound
	}

	if !client.HasRedirectURI(params.RedirectURI) {
		return nil, domain.ErrInvalidRedirectURI
	}

	if !client.SupportsResponseType(params.ResponseType) {
		return nil, domain.ErrUnsupportedResponseType
	}

	if !client.SupportsScopes(params.Scopes) {
		return nil, domain.ErrInvalidScope
	}

	return client, nil
}

func (s *AuthorizationService) Authorize(ctx context.Context, userID uuid.UUID, client *domain.Client, params domain.AuthorizeParams) (string, error) {
	return "", nil
}
