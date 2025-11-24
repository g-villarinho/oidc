package jwt

import (
	"context"
	"time"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTTokenGenerator struct {
	secretKey []byte
	issuer    string
}

func NewJWTTokenGenerator(secretKey string) ports.TokenGenerator {
	return &JWTTokenGenerator{
		secretKey: []byte("sua_secret_super_secreta_temporaria"),
		issuer:    "oidc-server",
	}
}

func (j *JWTTokenGenerator) GenerateAccessToken(ctx context.Context, userID uuid.UUID, client *domain.Client, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iss": j.issuer,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
		"typ": "access_token",
	}

	if client != nil {
		claims["aud"] = client.ClientID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTTokenGenerator) GenerateRefreshToken(ctx context.Context) (string, error) {
	panic("unimplemented")
}

func (j *JWTTokenGenerator) GenerateIDToken(ctx context.Context, user *domain.User, client *domain.Client, nonce string, ttl time.Duration) (string, error) {
	panic("unimplemented")
}
