//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/google/uuid"
)

// TestUserBuilder provides a fluent API for creating test users
type TestUserBuilder struct {
	user *domain.User
}

// NewTestUser creates a new user builder with default values
func NewTestUser() *TestUserBuilder {
	id, _ := uuid.NewV7()
	now := time.Now().UTC()

	return &TestUserBuilder{
		user: &domain.User{
			ID:            id,
			Name:          "Test User",
			Email:         "test@example.com",
			PasswordHash:  "$argon2id$v=19$m=65536,t=3,p=2$hash",
			EmailVerified: false,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
}

func (b *TestUserBuilder) WithID(id uuid.UUID) *TestUserBuilder {
	b.user.ID = id
	return b
}

func (b *TestUserBuilder) WithName(name string) *TestUserBuilder {
	b.user.Name = name
	return b
}

func (b *TestUserBuilder) WithEmail(email string) *TestUserBuilder {
	b.user.Email = email
	return b
}

func (b *TestUserBuilder) WithPasswordHash(hash string) *TestUserBuilder {
	b.user.PasswordHash = hash
	return b
}

func (b *TestUserBuilder) WithEmailVerified(verified bool) *TestUserBuilder {
	b.user.EmailVerified = verified
	return b
}

func (b *TestUserBuilder) Build() *domain.User {
	return b.user
}

// TestClientBuilder provides a fluent API for creating test clients
type TestClientBuilder struct {
	client *domain.Client
}

// NewTestClient creates a new client builder with default values
func NewTestClient() *TestClientBuilder {
	id, _ := uuid.NewV7()
	now := time.Now().UTC()

	return &TestClientBuilder{
		client: &domain.Client{
			ID:            id,
			ClientID:      "test-client-id",
			ClientSecret:  "test-client-secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{"http://localhost:3000/callback"},
			GrantTypes:    []string{"authorization_code", "refresh_token"},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile", "email"},
			LogoURL:       "https://example.com/logo.png",
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}
}

func (b *TestClientBuilder) WithID(id uuid.UUID) *TestClientBuilder {
	b.client.ID = id
	return b
}

func (b *TestClientBuilder) WithClientID(clientID string) *TestClientBuilder {
	b.client.ClientID = clientID
	return b
}

func (b *TestClientBuilder) WithClientSecret(secret string) *TestClientBuilder {
	b.client.ClientSecret = secret
	return b
}

func (b *TestClientBuilder) WithClientName(name string) *TestClientBuilder {
	b.client.ClientName = name
	return b
}

func (b *TestClientBuilder) WithRedirectURIs(uris []string) *TestClientBuilder {
	b.client.RedirectURIs = uris
	return b
}

func (b *TestClientBuilder) WithGrantTypes(types []string) *TestClientBuilder {
	b.client.GrantTypes = types
	return b
}

func (b *TestClientBuilder) WithResponseTypes(types []string) *TestClientBuilder {
	b.client.ResponseTypes = types
	return b
}

func (b *TestClientBuilder) WithScopes(scopes []string) *TestClientBuilder {
	b.client.Scopes = scopes
	return b
}

func (b *TestClientBuilder) WithLogoURL(url string) *TestClientBuilder {
	b.client.LogoURL = url
	return b
}

func (b *TestClientBuilder) Build() *domain.Client {
	return b.client
}

// TestAuthorizationCodeBuilder provides a fluent API for creating test auth codes
type TestAuthorizationCodeBuilder struct {
	code *domain.AuthorizationCode
}

// NewTestAuthorizationCode creates a new auth code builder with default values
func NewTestAuthorizationCode(clientID string, userID uuid.UUID) *TestAuthorizationCodeBuilder {
	return &TestAuthorizationCodeBuilder{
		code: &domain.AuthorizationCode{
			Code:                "test-auth-code-12345",
			ClientID:            clientID,
			UserID:              userID,
			RedirectURI:         "http://localhost:3000/callback",
			Scopes:              []string{"openid", "profile"},
			CodeChallenge:       "challenge123",
			CodeChallengeMethod: "S256",
			ExpiresAt:           time.Now().UTC().Add(10 * time.Minute),
			CreatedAt:           time.Now().UTC(),
		},
	}
}

func (b *TestAuthorizationCodeBuilder) WithCode(code string) *TestAuthorizationCodeBuilder {
	b.code.Code = code
	return b
}

func (b *TestAuthorizationCodeBuilder) WithRedirectURI(uri string) *TestAuthorizationCodeBuilder {
	b.code.RedirectURI = uri
	return b
}

func (b *TestAuthorizationCodeBuilder) WithScopes(scopes []string) *TestAuthorizationCodeBuilder {
	b.code.Scopes = scopes
	return b
}

func (b *TestAuthorizationCodeBuilder) WithCodeChallenge(challenge string) *TestAuthorizationCodeBuilder {
	b.code.CodeChallenge = challenge
	return b
}

func (b *TestAuthorizationCodeBuilder) WithCodeChallengeMethod(method string) *TestAuthorizationCodeBuilder {
	b.code.CodeChallengeMethod = method
	return b
}

func (b *TestAuthorizationCodeBuilder) WithExpiresAt(t time.Time) *TestAuthorizationCodeBuilder {
	b.code.ExpiresAt = t
	return b
}

func (b *TestAuthorizationCodeBuilder) Build() *domain.AuthorizationCode {
	return b.code
}

// MustCreateUser is a helper to create a user in the database for tests
func MustCreateUser(t *testing.T, db *TestDB, user *domain.User) {
	t.Helper()
	ctx := context.Background()

	query := `
		INSERT INTO users (id, email, password_hash, name, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Name,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
}

// MustCreateClient is a helper to create a client in the database for tests
func MustCreateClient(t *testing.T, db *TestDB, client *domain.Client) {
	t.Helper()
	ctx := context.Background()

	query := `
		INSERT INTO oauth_clients (id, client_id, client_secret, client_name, redirect_uris, grant_types, response_types, scopes, logo_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := db.Pool.Exec(ctx, query,
		client.ID,
		client.ClientID,
		client.ClientSecret,
		client.ClientName,
		client.RedirectURIs,
		client.GrantTypes,
		client.ResponseTypes,
		client.Scopes,
		client.LogoURL,
		client.CreatedAt,
		client.UpdatedAt,
	)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}
}
