package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type TokenService interface {
	CreateTokens(ctx context.Context, userID uuid.UUID, clientID string, scopes []string, authorizationCode *string, nonce string) (*TokenResponse, error)
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

func (s *TokenServiceImpl) CreateTokens(
	ctx context.Context,
	userID uuid.UUID,
	clientID string,
	scopes []string,
	authorizationCode *string,
	nonce string,
) (*TokenResponse, error) {
	// Generate JWT tokens
	accessToken, err := s.tokenGenerator.GenerateAccessToken(ctx, userID, clientID, scopes)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// Generate ID token if openid scope is present
	var idToken string
	if containsScope(scopes, "openid") {
		user, err := s.userRepository.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("get user for ID token: %w", err)
		}

		idToken, err = s.tokenGenerator.GenerateIDToken(ctx, user, clientID, nonce, scopes)
		if err != nil {
			return nil, fmt.Errorf("generate ID token: %w", err)
		}
	}

	// Create and persist token entity using config durations
	token, err := domain.NewToken(
		accessToken,
		refreshToken,
		authorizationCode,
		clientID,
		userID,
		scopes,
		s.config.JWT.AccessTokenDuration,
		s.config.JWT.RefreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("create token domain: %w", err)
	}

	if err := s.tokenRepository.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("save token: %w", err)
	}

	// Build response
	response := &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    domain.TokenTypeBearer,
		ExpiresIn:    int64(s.config.JWT.AccessTokenDuration.Seconds()),
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}

	if len(scopes) > 0 {
		response.Scope = joinScopes(scopes)
	}

	return response, nil
}

// Helper functions
func containsScope(scopes []string, scope string) bool {
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func joinScopes(scopes []string) string {
	result := ""
	for i, scope := range scopes {
		if i > 0 {
			result += " "
		}
		result += scope
	}
	return result
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
