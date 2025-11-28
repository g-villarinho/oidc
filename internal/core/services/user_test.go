package services

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/g-villarinho/oidc-server/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Run("should create user successfully when valid data is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, hashedPassword, user.PasswordHash)
		assert.False(t, user.EmailVerified)
		assert.NotEqual(t, "", user.ID.String())
	})

	t.Run("should return error when user already exists by email", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "existing@example.com"
		password := "SecurePassword123!"

		existingUser := &domain.User{
			Email: email,
			Name:  name,
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(existingUser, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should return error when GetByEmail fails with unexpected error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		unexpectedErr := errors.New("database connection failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, unexpectedErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "check existing user")
	})

	t.Run("should return error when password hashing fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		hashErr := errors.New("hashing failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return("", hashErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "hash password")
	})

	t.Run("should return error when repository Create fails with generic error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		createErr := errors.New("database write failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(createErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "persist user")
	})

	t.Run("should return ErrUserAlreadyExists when repository returns unique key violation", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(ports.ErrUniqueKeyViolation)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should handle empty name correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := ""
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "", user.Name)
	})

	t.Run("should handle empty email correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := ""
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "", user.Email)
	})

	t.Run("should handle empty password correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := ""
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, ports.ErrNotFound)

		mockHasher.EXPECT().
			Hash(ctx, password).
			Return(hashedPassword, nil)

		mockUserRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.User")).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		user, err := service.CreateUser(ctx, name, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, hashedPassword, user.PasswordHash)
	})
}

func TestAuthenticate(t *testing.T) {
	t.Run("should authenticate successfully when valid credentials are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		user := &domain.User{
			Email:         email,
			PasswordHash:  hashedPassword,
			EmailVerified: true,
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(ctx, password, hashedPassword).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, email, result.Email)
		assert.Equal(t, hashedPassword, result.PasswordHash)
	})

	t.Run("should return error when GetByEmail fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		dbErr := errors.New("database connection failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, dbErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get user by email")
	})

	t.Run("should return ErrUserNotFound when user does not exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"
		password := "SecurePassword123!"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should return ErrPasswordMismatch when password is incorrect", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "WrongPassword"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		user := &domain.User{
			Email:         email,
			PasswordHash:  hashedPassword,
			EmailVerified: true,
		}

		compareErr := errors.New("password does not match")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(ctx, password, hashedPassword).
			Return(compareErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})

	t.Run("should return ErrEmailNotVerified when email is not verified", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		user := &domain.User{
			Email:         email,
			PasswordHash:  hashedPassword,
			EmailVerified: false,
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(ctx, password, hashedPassword).
			Return(nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrEmailNotVerified)
	})

	t.Run("should handle empty email correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := ""
		password := "SecurePassword123!"

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(nil, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should handle empty password correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := ""
		hashedPassword := "$argon2id$v=19$m=65536,t=3,p=2$hashed"

		user := &domain.User{
			Email:         email,
			PasswordHash:  hashedPassword,
			EmailVerified: true,
		}

		compareErr := errors.New("password does not match")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockHasher.EXPECT().
			Compare(ctx, password, hashedPassword).
			Return(compareErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.Authenticate(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("should return user when valid ID is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		expectedUser := &domain.User{
			ID:            userID,
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			PasswordHash:  "$argon2id$v=19$m=65536,t=3,p=2$hashed",
			EmailVerified: true,
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(expectedUser, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.GetUserByID(ctx, userID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.Name, result.Name)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.PasswordHash, result.PasswordHash)
		assert.Equal(t, expectedUser.EmailVerified, result.EmailVerified)
	})

	t.Run("should return ErrUserNotFound when user does not exist", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.GetUserByID(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should return error when repository GetByID fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()

		dbErr := errors.New("database connection failed")

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, dbErr)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.GetUserByID(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get user by ID")
	})

	t.Run("should handle nil UUID correctly", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.Nil

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(nil, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.GetUserByID(ctx, userID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should return user with unverified email", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		expectedUser := &domain.User{
			ID:            userID,
			Name:          "Jane Doe",
			Email:         "jane.doe@example.com",
			PasswordHash:  "$argon2id$v=19$m=65536,t=3,p=2$hashed",
			EmailVerified: false,
		}

		mockUserRepo := mocks.NewUserRepositoryMock(t)
		mockHasher := mocks.NewHasherMock(t)
		logger := slog.Default()

		mockUserRepo.EXPECT().
			GetByID(ctx, userID).
			Return(expectedUser, nil)

		service := NewUserService(mockUserRepo, mockHasher, logger)

		// Act
		result, err := service.GetUserByID(ctx, userID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.False(t, result.EmailVerified)
	})
}