package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthorizationCode struct {
	Code                string
	ClientID            string
	UserID              uuid.UUID
	RedirectURI         string
	Scopes              []string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
	CreatedAt           time.Time
}

func NewAuthorizationCode(clientID string, userID uuid.UUID, redirectURI string, scopes []string, codeChallenge, codeChallengeMethod string, expiresAt time.Time) (*AuthorizationCode, error) {
	code, err := generateCode()
	if err != nil {
		return nil, err
	}

	return &AuthorizationCode{
		Code:                code,
		ClientID:            clientID,
		UserID:              userID,
		RedirectURI:         redirectURI,
		Scopes:              scopes,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           expiresAt,
	}, nil
}

func generateCode() (string, error) {
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

	code := base64.RawURLEncoding.EncodeToString(combined)

	return code, nil
}

func (ac *AuthorizationCode) IsExpired() bool {
	return time.Now().After(ac.ExpiresAt)
}
