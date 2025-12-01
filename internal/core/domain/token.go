package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

const (
	TokenTypeBearer = "Bearer"
)

var (
	ErrTokenExpired   = errors.New("token expired")
	ErrTokenRevoked   = errors.New("token revoked")
	ErrTokenNotFound  = errors.New("token not found")
	ErrInvalidToken   = errors.New("invalid token")
	ErrRefreshExpired = errors.New("refresh token expired")
	ErrNoRefreshToken = errors.New("no refresh token available")
)

type Token struct {
	ID                    uuid.UUID
	AccessTokenHash       string
	RefreshTokenHash      string
	AuthorizationCode     *string
	ClientID              string
	UserID                uuid.UUID
	Scopes                []string
	TokenType             string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Revoked               bool
	RevokedAt             *time.Time
	RevokedReason         *string
	CreatedAt             time.Time
	LastUsedAt            *time.Time
}

func NewToken(
	accessToken string,
	refreshToken string,
	authorizationCode *string,
	clientID string,
	userID uuid.UUID,
	scopes []string,
	accessTokenExpiresIn time.Duration,
	refreshTokenExpiresIn time.Duration,
) (*Token, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	accessTokenHash := hashToken(accessToken)
	refreshTokenHash := hashToken(refreshToken)

	return &Token{
		ID:                    id,
		AccessTokenHash:       accessTokenHash,
		RefreshTokenHash:      refreshTokenHash,
		AuthorizationCode:     authorizationCode,
		ClientID:              clientID,
		UserID:                userID,
		Scopes:                scopes,
		TokenType:             TokenTypeBearer,
		AccessTokenExpiresAt:  now.Add(accessTokenExpiresIn),
		RefreshTokenExpiresAt: now.Add(refreshTokenExpiresIn),
		Revoked:               false,
		CreatedAt:             now,
	}, nil
}

type CreateTokenParams struct {
	UserID            uuid.UUID
	ClientID          string
	Scopes            []string
	AuthorizationCode *string
	Nonce             string
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (t *Token) IsAccessTokenExpired() bool {
	return time.Now().UTC().After(t.AccessTokenExpiresAt)
}

func (t *Token) IsRefreshTokenExpired() bool {
	return time.Now().UTC().After(t.RefreshTokenExpiresAt)
}

func (t *Token) IsRevoked() bool {
	return t.Revoked
}

func (t *Token) IsValid() bool {
	return !t.IsRevoked() && !t.IsAccessTokenExpired()
}

func (t *Token) CanRefresh() bool {
	return !t.IsRefreshTokenExpired() && !t.IsRevoked()
}

func (t *Token) Revoke(reason string) {
	now := time.Now().UTC()
	t.Revoked = true
	t.RevokedAt = &now
	t.RevokedReason = &reason
}

func (t *Token) MarkAsUsed() {
	now := time.Now().UTC()
	t.LastUsedAt = &now
}

func (t *Token) HasScope(scope string) bool {
	return slices.Contains(t.Scopes, scope)
}

func (t *Token) HasAllScopes(scopes []string) bool {
	for _, scope := range scopes {
		if !t.HasScope(scope) {
			return false
		}
	}
	return true
}

func (t *Token) ValidateAccessToken(accessToken string) bool {
	return t.AccessTokenHash == hashToken(accessToken)
}

func (t *Token) ValidateRefreshToken(refreshToken string) bool {
	return t.RefreshTokenHash == hashToken(refreshToken)
}
