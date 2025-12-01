package jwt

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTTokenGenerator struct {
	jwtConfig *config.JWT
}

func NewJWTTokenGenerator(cfg *config.Config) ports.TokenGenerator {
	return &JWTTokenGenerator{
		jwtConfig: &cfg.JWT,
	}
}

func (j *JWTTokenGenerator) GenerateAccessToken(ctx context.Context, userID uuid.UUID, clientID string, scopes []string) (string, error) {
	secret := []byte(j.jwtConfig.Secret)

	claims := jwt.MapClaims{
		"iss":   j.jwtConfig.Issuer,
		"sub":   userID.String(),
		"aud":   clientID,
		"exp":   time.Now().Add(j.jwtConfig.AccessTokenDuration).Unix(),
		"iat":   time.Now().Unix(),
		"scope": strings.Join(scopes, " "),
		"jti":   uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (j *JWTTokenGenerator) GenerateRefreshToken(ctx context.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	entropy := fmt.Sprintf("%s:%s:%d",
		bytes,
		uuid.New().String(),
		time.Now().UnixNano(),
	)

	hash := sha256.Sum256([]byte(entropy))

	combined := append(bytes, hash[:]...)

	token := base64.RawURLEncoding.EncodeToString(combined)

	return token, nil
}

func (j *JWTTokenGenerator) GenerateIDToken(ctx context.Context, user *domain.User, clientID, nonce string, scopes []string) (string, error) {
	secret := []byte(j.jwtConfig.Secret)

	claims := jwt.MapClaims{
		"iss": j.jwtConfig.Issuer,
		"sub": user.ID.String(),
		"aud": clientID,
		"exp": time.Now().Add(j.jwtConfig.IDTokenDuration).Unix(),
		"iat": time.Now().Unix(),
	}

	if nonce != "" {
		claims["nonce"] = nonce
	}

	if slices.Contains(scopes, "email") {
		claims["email"] = user.Email
		claims["email_verified"] = user.EmailVerified
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
