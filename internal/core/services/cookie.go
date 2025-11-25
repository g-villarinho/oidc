package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
)

type CookieService struct {
	secret []byte
}

func NewCookieService(config *config.Config) *CookieService {
	return &CookieService{
		secret: []byte(config.Session.Secret),
	}
}

func (s *CookieService) SignSessionCookie(ctx context.Context, value string) string {
	return sign(value, s.secret)
}

func (s *CookieService) VerifySessionCookie(ctx context.Context, signedValue string) (string, error) {
	parts := strings.Split(signedValue, ".")
	if len(parts) != 2 {
		return "", domain.ErrInvalidSessionSignature
	}

	value, signature := parts[0], parts[1]

	expectedSignature := sign(value, s.secret)
	expectedParts := strings.Split(expectedSignature, ".")

	if len(expectedParts) == 2 && hmac.Equal([]byte(signature), []byte(expectedParts[1])) {
		return value, nil
	}

	return "", domain.ErrInvalidSessionSignature
}

func sign(value string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(value))
	signature := hex.EncodeToString(mac.Sum(nil))
	return value + "." + signature
}
