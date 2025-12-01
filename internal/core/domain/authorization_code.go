package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	AuthorizationCodeExpiry = 10 * time.Minute
)

type AuthorizationCode struct {
	Code                string
	ClientID            string
	UserID              uuid.UUID
	RedirectURI         string
	Scopes              []string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool
	ExpiresAt           time.Time
	CreatedAt           time.Time
}

func NewAuthorizationCode(clientID string, userID uuid.UUID, redirectURI string, scopes []string, nonce, codeChallenge, codeChallengeMethod string) (*AuthorizationCode, error) {
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
		Nonce:               nonce,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().Add(AuthorizationCodeExpiry),
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

func (ac *AuthorizationCode) IsValidPKCE(verifier string) bool {
	var computedChallenge string

	switch ac.CodeChallengeMethod {
	case "S256":
		hash := sha256.Sum256([]byte(verifier))
		computedChallenge = base64.RawURLEncoding.EncodeToString(hash[:])
	case "plain":
		computedChallenge = verifier
	default:
		return false
	}

	return computedChallenge == ac.CodeChallenge
}

func (ac *AuthorizationCode) ToCreateTokenParams() CreateTokenParams {
	return CreateTokenParams{
		AuthorizationCode: &ac.Code,
		ClientID:          ac.ClientID,
		UserID:            ac.UserID,
		Scopes:            ac.Scopes,
		Nonce:             ac.Nonce,
	}
}
