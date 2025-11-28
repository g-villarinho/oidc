//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/repositories"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should create client successfully", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("my-test-client").
			WithClientName("My Test Client").
			Build()

		err := repo.Create(ctx, client)

		require.NoError(t, err)

		// Verify client was created
		found, err := repo.GetByClientID(ctx, client.ClientID)
		require.NoError(t, err)
		assert.Equal(t, client.ClientID, found.ClientID)
		assert.Equal(t, client.ClientName, found.ClientName)
		assert.Equal(t, client.ClientSecret, found.ClientSecret)
		assert.Equal(t, client.RedirectURIs, found.RedirectURIs)
		assert.Equal(t, client.GrantTypes, found.GrantTypes)
		assert.Equal(t, client.ResponseTypes, found.ResponseTypes)
		assert.Equal(t, client.Scopes, found.Scopes)
		assert.Equal(t, client.LogoURL, found.LogoURL)
	})

	t.Run("should return error when client_id already exists", func(t *testing.T) {
		db.TruncateTables(t)

		client1 := NewTestClient().
			WithClientID("duplicate-client").
			Build()

		client2 := NewTestClient().
			WithClientID("duplicate-client").
			Build()

		err := repo.Create(ctx, client1)
		require.NoError(t, err)

		err = repo.Create(ctx, client2)

		require.Error(t, err)
	})
}

func TestClientRepository_GetByClientID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return client when client_id exists", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("findme-client").
			WithClientName("Find Me Client").
			Build()

		err := repo.Create(ctx, client)
		require.NoError(t, err)

		found, err := repo.GetByClientID(ctx, "findme-client")

		require.NoError(t, err)
		assert.Equal(t, client.ID, found.ID)
		assert.Equal(t, client.ClientID, found.ClientID)
		assert.Equal(t, client.ClientName, found.ClientName)
	})

	t.Run("should return ErrNotFound when client_id does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		found, err := repo.GetByClientID(ctx, "nonexistent-client")

		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrNotFound)
		assert.Nil(t, found)
	})
}

func TestClientRepository_GetByID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return client when ID exists", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("byid-client").
			Build()

		err := repo.Create(ctx, client)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, client.ID)

		require.NoError(t, err)
		assert.Equal(t, client.ID, found.ID)
		assert.Equal(t, client.ClientID, found.ClientID)
	})

	t.Run("should return ErrNotFound when ID does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		randomID := uuid.New()

		found, err := repo.GetByID(ctx, randomID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrNotFound)
		assert.Nil(t, found)
	})
}

func TestClientRepository_List(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return empty list when no clients exist", func(t *testing.T) {
		db.TruncateTables(t)

		clients, err := repo.List(ctx)

		require.NoError(t, err)
		assert.Empty(t, clients)
	})

	t.Run("should return all clients", func(t *testing.T) {
		db.TruncateTables(t)

		client1 := NewTestClient().
			WithClientID("client-1").
			WithClientName("Client 1").
			Build()

		client2 := NewTestClient().
			WithClientID("client-2").
			WithClientName("Client 2").
			Build()

		client3 := NewTestClient().
			WithClientID("client-3").
			WithClientName("Client 3").
			Build()

		require.NoError(t, repo.Create(ctx, client1))
		require.NoError(t, repo.Create(ctx, client2))
		require.NoError(t, repo.Create(ctx, client3))

		clients, err := repo.List(ctx)

		require.NoError(t, err)
		assert.Len(t, clients, 3)
	})
}

func TestClientRepository_Update(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should update client successfully", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("update-client").
			WithClientName("Original Name").
			WithScopes([]string{"openid"}).
			Build()

		err := repo.Create(ctx, client)
		require.NoError(t, err)

		// Update client
		client.ClientName = "Updated Name"
		client.Scopes = []string{"openid", "profile", "email"}
		client.RedirectURIs = []string{"http://localhost:4000/callback"}

		err = repo.Update(ctx, client)
		require.NoError(t, err)

		// Verify update
		found, err := repo.GetByID(ctx, client.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.ClientName)
		assert.Equal(t, []string{"openid", "profile", "email"}, found.Scopes)
		assert.Equal(t, []string{"http://localhost:4000/callback"}, found.RedirectURIs)
	})
}

func TestClientRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewClientRepository(db.Pool)
	ctx := context.Background()

	t.Run("should delete client successfully", func(t *testing.T) {
		db.TruncateTables(t)

		client := NewTestClient().
			WithClientID("delete-client").
			Build()

		err := repo.Create(ctx, client)
		require.NoError(t, err)

		// Verify client exists
		_, err = repo.GetByID(ctx, client.ID)
		require.NoError(t, err)

		// Delete client
		err = repo.Delete(ctx, client.ID)
		require.NoError(t, err)

		// Verify client was deleted
		found, err := repo.GetByID(ctx, client.ID)
		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrNotFound)
		assert.Nil(t, found)
	})

	t.Run("should not return error when deleting non-existent client", func(t *testing.T) {
		db.TruncateTables(t)

		randomID := uuid.New()

		err := repo.Delete(ctx, randomID)

		// PostgreSQL DELETE does not return error for non-existent rows
		require.NoError(t, err)
	})
}
