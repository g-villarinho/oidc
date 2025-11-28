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

type CookieService interface {
	SignSessionCookie(ctx context.Context, value string) string
	VerifySessionCookie(ctx context.Context, signedValue string) (string, error)
}

type CookieServiceImpl struct {
	secret []byte
}

func NewCookieService(config *config.Config) CookieService {
	return &CookieServiceImpl{
		secret: []byte(config.Session.Secret),
	}
}

func (s *CookieServiceImpl) SignSessionCookie(ctx context.Context, value string) string {
	return sign(value, s.secret)
}

func (s *CookieServiceImpl) VerifySessionCookie(ctx context.Context, signedValue string) (string, error) {
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
