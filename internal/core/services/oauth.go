package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type OAuthService interface {
	VerifyAuthorization(ctx context.Context, params domain.AuthorizeParams) error
	CreateAuthorizationCode(ctx context.Context, userID uuid.UUID, params domain.AuthorizeParams) (*domain.AuthorizationCode, error)
	ExchangeToken(ctx context.Context, params domain.ExchangeTokenParams) error
}

type OAuthServiceImpl struct {
	clientRepository            ports.ClientRepository
	authorizationCodeRepository ports.AuthorizationCodeRepository
}

func NewOAuthService(clientRepository ports.ClientRepository, authorizationCodeRepository ports.AuthorizationCodeRepository) OAuthService {
	return &OAuthServiceImpl{
		clientRepository:            clientRepository,
		authorizationCodeRepository: authorizationCodeRepository,
	}
}

func (s *OAuthServiceImpl) VerifyAuthorization(ctx context.Context, params domain.AuthorizeParams) error {
	client, err := s.clientRepository.GetByClientID(ctx, params.ClientID)
	if err != nil {
		return fmt.Errorf("get validated OAuth client: %w", err)
	}

	if client == nil {
		return domain.ErrClientNotFound
	}

	if !client.HasRedirectURI(params.RedirectURI) {
		return domain.ErrInvalidRedirectURI
	}

	if !client.SupportsResponseType(params.ResponseType) {
		return domain.ErrUnsupportedResponseType
	}

	if !client.SupportsScopes(params.Scopes) {
		return domain.ErrInvalidScope
	}

	return nil
}

func (s *OAuthServiceImpl) CreateAuthorizationCode(ctx context.Context, userID uuid.UUID, params domain.AuthorizeParams) (*domain.AuthorizationCode, error) {
	authorizationCode, err := domain.NewAuthorizationCode(
		params.ClientID,
		userID,
		params.RedirectURI,
		params.Scopes,
		params.CodeChallenge,
		params.CodeChallengeMethod,
	)

	if err != nil {
		return nil, fmt.Errorf("create authorization code: %w", err)
	}

	if err := s.authorizationCodeRepository.Create(ctx, authorizationCode); err != nil {
		return nil, fmt.Errorf("save authorization code: %w", err)
	}

	return authorizationCode, nil
}

func (s *OAuthServiceImpl) ExchangeToken(ctx context.Context, params domain.ExchangeTokenParams) error {
	authorizationCode, err := s.authorizationCodeRepository.GetByCode(ctx, params.Code)
	if err != nil {
		if err == ports.ErrNotFound {
			return domain.ErrInvalidAuthorizationCode
		}

		return fmt.Errorf("get authorization code: %w", err)
	}

	if authorizationCode.Used {
		//TODO: Revoke all tokens issued with this code's grant
		return domain.ErrAuthorizationCodeAlreadyUsed
	}

	if authorizationCode.IsExpired() {
		return domain.ErrAuthorizationCodeExpired
	}

	// Validation of client pkce, client_id redirect_uri should be done here

	if err := s.authorizationCodeRepository.MarkAsUsed(ctx, authorizationCode.Code); err != nil {
		return fmt.Errorf("mark authorization code as used: %w", err)
	}

	return nil
}
