package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type OAuthService interface {
	VerifyAuthorization(ctx context.Context, params domain.AuthorizeParams) error
	CreateAuthorizationCode(ctx context.Context, userID uuid.UUID, params domain.AuthorizeParams) (*domain.AuthorizationCode, error)
	ExchangeToken(ctx context.Context, params domain.ExchangeTokenParams) (*TokenResponse, error)
}

type OAuthServiceImpl struct {
	clientRepository            ports.ClientRepository
	authorizationCodeRepository ports.AuthorizationCodeRepository
	tokenService                TokenService
	userRepository              ports.UserRepository
	config                      *config.Config
}

func NewOAuthService(
	clientRepository ports.ClientRepository,
	authorizationCodeRepository ports.AuthorizationCodeRepository,
	tokenService TokenService,
	userRepository ports.UserRepository,
	config *config.Config,
) OAuthService {
	return &OAuthServiceImpl{
		clientRepository:            clientRepository,
		authorizationCodeRepository: authorizationCodeRepository,
		tokenService:                tokenService,
		userRepository:              userRepository,
		config:                      config,
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
		params.Nonce,
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

func (s *OAuthServiceImpl) ExchangeToken(ctx context.Context, params domain.ExchangeTokenParams) (*TokenResponse, error) {
	switch params.GrantType {
	case "authorization_code":
		return s.exchangeAuthorizationCode(ctx, params)
	// case "refresh_token":
	// 	return s.exchangeRefreshToken(ctx, params)
	default:
		return nil, domain.ErrUnsupportedResponseType
	}
}

func (s *OAuthServiceImpl) exchangeAuthorizationCode(ctx context.Context, params domain.ExchangeTokenParams) (*TokenResponse, error) {
	authorizationCode, err := s.authorizationCodeRepository.GetByCode(ctx, params.Code)
	if err != nil {
		if err == ports.ErrNotFound {
			return nil, domain.ErrInvalidAuthorizationCode
		}

		return nil, fmt.Errorf("get authorization code: %w", err)
	}

	if authorizationCode.Used {
		//TODO: Revoke all tokens issued with this code's grant
		return nil, domain.ErrAuthorizationCodeAlreadyUsed
	}

	if authorizationCode.IsExpired() {
		return nil, domain.ErrAuthorizationCodeExpired
	}

	if authorizationCode.ClientID != params.ClientID {
		return nil, domain.ErrUnauthorizedClient
	}

	if authorizationCode.RedirectURI != params.RedirectURI {
		return nil, domain.ErrInvalidRedirectURI
	}

	if !authorizationCode.IsValidPKCE(params.CodeVerifier) {
		return nil, domain.ErrInvalidPKCEVerification
	}

	if err := s.authorizationCodeRepository.MarkAsUsed(ctx, authorizationCode.Code); err != nil {
		return nil, fmt.Errorf("mark authorization code as used: %w", err)
	}

	tokenResponse, err := s.tokenService.CreateTokens(ctx, authorizationCode.ToCreateTokenParams())
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}

	return tokenResponse, nil
}

// func (s *OAuthServiceImpl) exchangeRefreshToken(ctx context.Context, params domain.ExchangeTokenParams) (*TokenResponse, error) {
// 	tokenResponse, err := s.tokenService.RefreshTokens(
// 		ctx,
// 		params.RefreshToken,
// 		s.config.JWT.AccessTokenDuration,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("refresh tokens: %w", err)
// 	}

// 	return tokenResponse, nil
// }
