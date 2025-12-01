package services

import (
	"context"
	"fmt"
	"slices"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

type TokenService interface {
	CreateTokens(ctx context.Context, params domain.CreateTokenParams) (*TokenResponse, error)
}

type TokenServiceImpl struct {
	tokenRepository ports.TokenRepository
	tokenGenerator  ports.TokenGenerator
	userRepository  ports.UserRepository
	config          *config.Config
}

func NewTokenService(
	tokenRepository ports.TokenRepository,
	tokenGenerator ports.TokenGenerator,
	userRepository ports.UserRepository,
	cfg *config.Config,
) TokenService {
	return &TokenServiceImpl{
		tokenRepository: tokenRepository,
		tokenGenerator:  tokenGenerator,
		userRepository:  userRepository,
		config:          cfg,
	}
}

func (s *TokenServiceImpl) CreateTokens(ctx context.Context, params domain.CreateTokenParams) (*TokenResponse, error) {
	accessToken, err := s.tokenGenerator.GenerateAccessToken(ctx, params.UserID, params.ClientID, params.Scopes)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	var idToken string
	if slices.Contains(params.Scopes, "openid") {
		user, err := s.userRepository.GetByID(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("get user for ID token: %w", err)
		}

		idToken, err = s.tokenGenerator.GenerateIDToken(ctx, user, params.ClientID, params.Nonce, params.Scopes)
		if err != nil {
			return nil, fmt.Errorf("generate ID token: %w", err)
		}
	}

	token, err := domain.NewToken(
		accessToken,
		refreshToken,
		params.AuthorizationCode,
		params.ClientID,
		params.UserID,
		params.Scopes,
		s.config.JWT.AccessTokenDuration,
		s.config.JWT.RefreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("create token domain: %w", err)
	}

	if err := s.tokenRepository.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("save token: %w", err)
	}

	response := &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    domain.TokenTypeBearer,
		ExpiresIn:    int64(s.config.JWT.AccessTokenDuration.Seconds()),
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}

	return response, nil
}
