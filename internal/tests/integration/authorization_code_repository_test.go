//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationCodeRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewAuthorizationCodeRepository(db.Pool)
	ctx := context.Background()

	t.Run("should create authorization code successfully", func(t *testing.T) {
		db.TruncateTables(t)

		// Create required dependencies (client and user)
		client := NewTestClient().
			WithClientID("auth-code-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("auth-code-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		// Create authorization code
		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("test-code-123").
			WithScopes([]string{"openid", "profile"}).
			Build()

		err := repo.Create(ctx, authCode)

		require.NoError(t, err)

		// Verify code was created using GetByCode
		found, err := repo.GetByCode(ctx, authCode.Code)
		require.NoError(t, err)
		assert.Equal(t, authCode.Code, found.Code)
		assert.Equal(t, authCode.ClientID, found.ClientID)
		assert.Equal(t, authCode.UserID, found.UserID)
		assert.Equal(t, authCode.RedirectURI, found.RedirectURI)
		assert.Equal(t, authCode.Scopes, found.Scopes)
		assert.Equal(t, authCode.CodeChallenge, found.CodeChallenge)
		assert.Equal(t, authCode.CodeChallengeMethod, found.CodeChallengeMethod)
	})

	t.Run("should return error when client does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		user := NewTestUser().
			WithEmail("orphan-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		// Try to create auth code with non-existent client
		authCode := NewTestAuthorizationCode("non-existent-client", user.ID).
			WithCode("orphan-code").
			Build()

		err := repo.Create(ctx, authCode)

		require.Error(t, err) // Foreign key violation
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("orphan-client").
			Build()
		MustCreateClient(t, db, client)

		// Try to create auth code with non-existent user
		nonExistentUserID := NewTestUser().Build().ID
		authCode := NewTestAuthorizationCode(client.ClientID, nonExistentUserID).
			WithCode("orphan-user-code").
			Build()

		err := repo.Create(ctx, authCode)

		require.Error(t, err) // Foreign key violation
	})

	t.Run("should create authorization code with empty code challenge", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("no-pkce-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("no-pkce-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("no-pkce-code").
			WithCodeChallenge("").
			WithCodeChallengeMethod("").
			Build()

		err := repo.Create(ctx, authCode)

		require.NoError(t, err)

		// Verify directly in database
		var codeChallenge *string
		err = db.Pool.QueryRow(ctx, "SELECT code_challenge FROM authorization_codes WHERE code = $1", authCode.Code).Scan(&codeChallenge)
		require.NoError(t, err)
	})
}

func TestAuthorizationCodeRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewAuthorizationCodeRepository(db.Pool)
	ctx := context.Background()

	t.Run("should delete authorization code successfully", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("delete-code-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("delete-code-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("delete-me-code").
			Build()

		err := repo.Create(ctx, authCode)
		require.NoError(t, err)

		// Verify code exists
		_, err = repo.GetByCode(ctx, authCode.Code)
		require.NoError(t, err)

		// Delete code
		err = repo.Delete(ctx, authCode.Code)
		require.NoError(t, err)

		// Verify code was deleted
		found, err := repo.GetByCode(ctx, authCode.Code)
		require.Error(t, err)
		assert.ErrorIs(t, err, pgx.ErrNoRows)
		assert.Nil(t, found)
	})

	t.Run("should not return error when deleting non-existent code", func(t *testing.T) {
		db.TruncateTables(t)

		err := repo.Delete(ctx, "non-existent-code")

		// PostgreSQL DELETE does not return error for non-existent rows
		require.NoError(t, err)
	})
}

func TestAuthorizationCodeRepository_GetByCode(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewAuthorizationCodeRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return authorization code when it exists", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("get-code-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("get-code-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("get-me-code").
			Build()

		err := repo.Create(ctx, authCode)
		require.NoError(t, err)

		found, err := repo.GetByCode(ctx, "get-me-code")

		require.NoError(t, err)
		assert.Equal(t, authCode.Code, found.Code)
		assert.Equal(t, authCode.ClientID, found.ClientID)
		assert.Equal(t, authCode.UserID, found.UserID)
	})

	t.Run("should return error when code does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		found, err := repo.GetByCode(ctx, "non-existent-code")

		require.Error(t, err)
		assert.ErrorIs(t, err, pgx.ErrNoRows)
		assert.Nil(t, found)
	})

	t.Run("should not return expired authorization code", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("expired-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("expired-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		// Create code that expires in the past
		expiredTime := time.Now().UTC().Add(-1 * time.Hour)

		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("expired-code").
			WithExpiresAt(expiredTime).
			Build()

		err := repo.Create(ctx, authCode)
		require.NoError(t, err)

		// GetByCode should not return expired codes (query has expires_at > NOW())
		found, err := repo.GetByCode(ctx, authCode.Code)
		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestAuthorizationCodeRepository_Expiration(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	ctx := context.Background()

	t.Run("should store and retrieve expiration time correctly", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("expiry-client").
			Build()
		MustCreateClient(t, db, client)

		user := NewTestUser().
			WithEmail("expiry-user@example.com").
			Build()
		MustCreateUser(t, db, user)

		expiresAt := time.Now().UTC().Add(5 * time.Minute)

		authCode := NewTestAuthorizationCode(client.ClientID, user.ID).
			WithCode("expiry-code").
			WithExpiresAt(expiresAt).
			Build()

		repo := repositories.NewAuthorizationCodeRepository(db.Pool)
		err := repo.Create(ctx, authCode)
		require.NoError(t, err)

		// Verify using GetByCode
		found, err := repo.GetByCode(ctx, authCode.Code)
		require.NoError(t, err)

		// Allow 1 second tolerance for time comparison
		assert.WithinDuration(t, expiresAt, found.ExpiresAt, time.Second)
	})
}
