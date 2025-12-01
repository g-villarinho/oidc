package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenExpiry  = 1 * time.Hour
	RefreshTokenExpiry = 24 * time.Hour * 30 // 30 days
	TokenTypeBearer    = "Bearer"
)

var (
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenRevoked    = errors.New("token revoked")
	ErrTokenNotFound   = errors.New("token not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrRefreshExpired  = errors.New("refresh token expired")
	ErrNoRefreshToken  = errors.New("no refresh token available")
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

type NewTokenParams struct {
	AccessToken           string
	RefreshToken          string
	AuthorizationCode     *string
	ClientID              string
	UserID                uuid.UUID
	Scopes                []string
	AccessTokenExpiresIn  time.Duration
	RefreshTokenExpiresIn time.Duration
}

func NewToken(params NewTokenParams) (*Token, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	accessTokenHash := hashToken(params.AccessToken)
	refreshTokenHash := hashToken(params.RefreshToken)

	accessTokenExpiresIn := params.AccessTokenExpiresIn
	if accessTokenExpiresIn == 0 {
		accessTokenExpiresIn = AccessTokenExpiry
	}

	refreshTokenExpiresIn := params.RefreshTokenExpiresIn
	if refreshTokenExpiresIn == 0 {
		refreshTokenExpiresIn = RefreshTokenExpiry
	}

	return &Token{
		ID:                    id,
		AccessTokenHash:       accessTokenHash,
		RefreshTokenHash:      refreshTokenHash,
		AuthorizationCode:     params.AuthorizationCode,
		ClientID:              params.ClientID,
		UserID:                params.UserID,
		Scopes:                params.Scopes,
		TokenType:             TokenTypeBearer,
		AccessTokenExpiresAt:  now.Add(accessTokenExpiresIn),
		RefreshTokenExpiresAt: now.Add(refreshTokenExpiresIn),
		Revoked:               false,
		CreatedAt:             now,
	}, nil
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
	for _, s := range t.Scopes {
		if s == scope {
			return true
		}
	}
	return false
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
