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

func TestUserRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewUserRepository(db.Pool)
	ctx := context.Background()

	t.Run("should create user successfully", func(t *testing.T) {
		db.TruncateTables(t)

		user := NewTestUser().
			WithEmail("john@example.com").
			WithName("John Doe").
			Build()

		err := repo.Create(ctx, user)

		require.NoError(t, err)

		// Verify user was created
		found, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.Name, found.Name)
		assert.Equal(t, user.PasswordHash, found.PasswordHash)
		assert.Equal(t, user.EmailVerified, found.EmailVerified)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		db.TruncateTables(t)

		user1 := NewTestUser().
			WithEmail("duplicate@example.com").
			Build()

		user2 := NewTestUser().
			WithEmail("duplicate@example.com").
			Build()

		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		err = repo.Create(ctx, user2)

		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrUniqueKeyViolation)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewUserRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return user when email exists", func(t *testing.T) {
		db.TruncateTables(t)

		user := NewTestUser().
			WithEmail("findme@example.com").
			WithName("Find Me").
			Build()

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.GetByEmail(ctx, "findme@example.com")

		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.Email, found.Email)
		assert.Equal(t, user.Name, found.Name)
	})

	t.Run("should return ErrNotFound when email does not exist", func(t *testing.T) {
		db.TruncateTables(t)

		found, err := repo.GetByEmail(ctx, "nonexistent@example.com")

		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrNotFound)
		assert.Nil(t, found)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewUserRepository(db.Pool)
	ctx := context.Background()

	t.Run("should return user when ID exists", func(t *testing.T) {
		db.TruncateTables(t)

		user := NewTestUser().
			WithEmail("byid@example.com").
			Build()

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		found, err := repo.GetByID(ctx, user.ID)

		require.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
		assert.Equal(t, user.Email, found.Email)
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

func TestUserRepository_Update(t *testing.T) {
	db := SetupTestDB(t)
	defer db.Teardown(t)

	repo := repositories.NewUserRepository(db.Pool)
	ctx := context.Background()

	t.Run("should update user successfully", func(t *testing.T) {
		db.TruncateTables(t)

		user := NewTestUser().
			WithEmail("update@example.com").
			WithName("Original Name").
			WithEmailVerified(false).
			Build()

		err := repo.Create(ctx, user)
		require.NoError(t, err)

		// Update user
		user.Name = "Updated Name"
		user.EmailVerified = true

		err = repo.Update(ctx, user)
		require.NoError(t, err)

		// Verify update
		found, err := repo.GetByID(ctx, user.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", found.Name)
		assert.True(t, found.EmailVerified)
	})

	t.Run("should return error when updating to existing email", func(t *testing.T) {
		db.TruncateTables(t)

		user1 := NewTestUser().
			WithEmail("first@example.com").
			Build()

		user2 := NewTestUser().
			WithEmail("second@example.com").
			Build()

		err := repo.Create(ctx, user1)
		require.NoError(t, err)

		err = repo.Create(ctx, user2)
		require.NoError(t, err)

		// Try to update user2's email to user1's email
		user2.Email = "first@example.com"

		err = repo.Update(ctx, user2)

		require.Error(t, err)
		assert.ErrorIs(t, err, ports.ErrUniqueKeyViolation)
	})
}
