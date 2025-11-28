//go:build integration

package integration

import (
	"bytes"
	"context"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/argon2"
	pgRepo "github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/repositories"
	redisRepo "github.com/g-villarinho/oidc-server/internal/adapters/secondary/redis/repositories"
	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/g-villarinho/oidc-server/pkg/validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TestUserBuilder struct {
	user *domain.User
}

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

type TestClientBuilder struct {
	client *domain.Client
}

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

type TestAuthorizationCodeBuilder struct {
	code *domain.AuthorizationCode
}

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

func MustHashPassword(t *testing.T, password string) string {
	t.Helper()
	ctx := context.Background()

	hasher := NewTestHasher()
	hash, err := hasher.Hash(ctx, password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return hash
}

type TestServices struct {
	UserService services.UserService
	AuthService services.AuthService
}

func NewTestHasher() ports.Hasher {
	return argon2.NewHasher()
}

func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

func NewTestConfig() *config.Config {
	return &config.Config{
		Session: config.Session{
			Duration: 24 * time.Hour,
			Secret:   "test-secret-key-for-integration-tests-min-32-chars",
			CookieOptions: config.CookieOptions{
				Name:   "_oidc.sid",
				Secure: false,
			},
		},
	}
}

func SetupTestServices(t *testing.T, env *TestEnv) *TestServices {
	t.Helper()

	userRepo := pgRepo.NewUserRepository(env.DB.Pool)
	sessionRepo := redisRepo.NewSessionRepository(env.Redis.Client)

	hasher := NewTestHasher()
	logger := NewTestLogger()
	cfg := NewTestConfig()

	userService := services.NewUserService(userRepo, hasher, logger)
	authService := services.NewAuthService(userService, userRepo, sessionRepo, cfg)

	return &TestServices{
		UserService: userService,
		AuthService: authService,
	}
}

type TestHTTPServer struct {
	Services *TestServices
	Config   *config.Config
	Logger   *slog.Logger
	Hasher   ports.Hasher
}

func SetupTestHTTPServer(t *testing.T, env *TestEnv) *TestHTTPServer {
	t.Helper()

	services := SetupTestServices(t, env)
	cfg := NewTestConfig()
	logger := NewTestLogger()
	hasher := NewTestHasher()

	return &TestHTTPServer{
		Services: services,
		Config:   cfg,
		Logger:   logger,
		Hasher:   hasher,
	}
}

// MakeRequest creates an Echo HTTP test request with validation configured.
// Returns the Echo context and response recorder for testing handlers.
func MakeRequest(method, target string, body []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = validation.NewValidator()

	req := httptest.NewRequest(method, target, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}
