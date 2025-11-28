package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	t.Run("should register user successfully when valid data is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		expectedUser := &domain.User{
			ID:            uuid.New(),
			Name:          name,
			Email:         email,
			EmailVerified: false,
		}

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(expectedUser, nil)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.NoError(t, err)
	})

	t.Run("should return error when user service fails to create user", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		expectedError := errors.New("database connection error")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "existing@example.com"
		password := "SecurePassword123!"

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, domain.ErrUserAlreadyExists)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
	})

	t.Run("should return error when empty name is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := ""
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		expectedError := errors.New("name is required")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
	})

	t.Run("should return error when empty email is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := ""
		password := "SecurePassword123!"

		expectedError := errors.New("email is required")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
	})

	t.Run("should return error when empty password is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := ""

		expectedError := errors.New("password is required")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
	})

	t.Run("should return error when invalid email format is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "invalid-email-format"
		password := "SecurePassword123!"

		expectedError := errors.New("invalid email format")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
	})

	t.Run("should return error when weak password is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		name := "John Doe"
		email := "john.doe@example.com"
		password := "123"

		expectedError := errors.New("password too weak")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			CreateUser(ctx, name, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		err := authService.RegisterUser(ctx, name, email, password)

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "register user")
	})
}

func TestLogin(t *testing.T) {
	t.Run("should login successfully when valid credentials are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"
		userID := uuid.New()

		expectedUser := &domain.User{
			ID:            userID,
			Name:          "John Doe",
			Email:         email,
			EmailVerified: true,
		}

		sessionConfig := config.Session{
			Duration: 24 * time.Hour,
		}

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(expectedUser, nil)

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.Session")).
			Return(nil)

		authService := &AuthServiceImpl{
			userService:       mockUserService,
			sessionRepository: mockSessionRepository,
			sessionConfig:     sessionConfig,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("should return error when authentication fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "WrongPassword"

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(nil, domain.ErrPasswordMismatch)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login user")
		assert.ErrorIs(t, err, domain.ErrPasswordMismatch)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "nonexistent@example.com"
		password := "SecurePassword123!"

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(nil, domain.ErrUserNotFound)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login user")
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should create session with configured duration", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		expectedUser := &domain.User{
			ID:            uuid.New(),
			Name:          "John Doe",
			Email:         email,
			EmailVerified: true,
		}

		sessionConfig := config.Session{
			Duration: 2 * time.Hour,
		}

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(expectedUser, nil)

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.Session")).
			Return(nil)

		authService := &AuthServiceImpl{
			userService:       mockUserService,
			sessionRepository: mockSessionRepository,
			sessionConfig:     sessionConfig,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotNil(t, user)
		assert.True(t, session.TTL() <= 2*time.Hour)
		assert.True(t, session.TTL() > 0)
	})

	t.Run("should return error when session repository fails to store session", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := "SecurePassword123!"

		expectedUser := &domain.User{
			ID:            uuid.New(),
			Name:          "John Doe",
			Email:         email,
			EmailVerified: true,
		}

		sessionConfig := config.Session{
			Duration: 24 * time.Hour,
		}

		expectedError := errors.New("redis connection error")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(expectedUser, nil)

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.Session")).
			Return(expectedError)

		authService := &AuthServiceImpl{
			userService:       mockUserService,
			sessionRepository: mockSessionRepository,
			sessionConfig:     sessionConfig,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "store session")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when user email is not verified", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "unverified@example.com"
		password := "SecurePassword123!"

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(nil, domain.ErrEmailNotVerified)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login user")
		assert.ErrorIs(t, err, domain.ErrEmailNotVerified)
	})

	t.Run("should return error when empty email is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := ""
		password := "SecurePassword123!"

		expectedError := errors.New("email is required")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login user")
	})

	t.Run("should return error when empty password is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		email := "john.doe@example.com"
		password := ""

		expectedError := errors.New("password is required")

		mockUserService := mocks.NewUserServiceMock(t)
		mockUserService.EXPECT().
			Authenticate(ctx, email, password).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			userService: mockUserService,
		}

		// Act
		session, user, err := authService.Login(ctx, email, password)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "login user")
	})
}

func TestGetSessionUser(t *testing.T) {
	t.Run("should return user when valid session ID is provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()

		expectedSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		expectedUser := &domain.User{
			ID:            userID,
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			EmailVerified: true,
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(expectedSession, nil)

		mockUserRepository := mocks.NewUserRepositoryMock(t)
		mockUserRepository.EXPECT().
			GetByID(ctx, userID).
			Return(expectedUser, nil)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
			userRepository:    mockUserRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.Name, user.Name)
	})

	t.Run("should return error when session repository fails to get session", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		expectedError := errors.New("database connection error")

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "get session")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when session is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(nil, nil)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrSessionNotFound)
	})

	t.Run("should return error and delete session when session is expired", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()

		expiredSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(expiredSession, nil)
		mockSessionRepository.EXPECT().
			Delete(ctx, sessionID).
			Return(nil)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrSessionExpired)
	})

	t.Run("should return error when session is expired and delete fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()
		deleteError := errors.New("redis connection error")

		expiredSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(expiredSession, nil)
		mockSessionRepository.EXPECT().
			Delete(ctx, sessionID).
			Return(deleteError)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "delete expired session")
		assert.Contains(t, err.Error(), deleteError.Error())
	})

	t.Run("should return error when user repository fails to get user", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()
		expectedError := errors.New("database connection error")

		validSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(validSession, nil)

		mockUserRepository := mocks.NewUserRepositoryMock(t)
		mockUserRepository.EXPECT().
			GetByID(ctx, userID).
			Return(nil, expectedError)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
			userRepository:    mockUserRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "get user by session")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()

		validSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
			CreatedAt: time.Now(),
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(validSession, nil)

		mockUserRepository := mocks.NewUserRepositoryMock(t)
		mockUserRepository.EXPECT().
			GetByID(ctx, userID).
			Return(nil, nil)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
			userRepository:    mockUserRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})

	t.Run("should verify session is not expired when close to expiration", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		sessionID := uuid.New()
		userID := uuid.New()

		// Session expires in 1 minute (still valid)
		almostExpiredSession := &domain.Session{
			ID:        sessionID,
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Minute),
			CreatedAt: time.Now().Add(-23 * time.Hour),
		}

		expectedUser := &domain.User{
			ID:            userID,
			Name:          "John Doe",
			Email:         "john.doe@example.com",
			EmailVerified: true,
		}

		mockSessionRepository := mocks.NewSessionRepositoryMock(t)
		mockSessionRepository.EXPECT().
			GetByID(ctx, sessionID).
			Return(almostExpiredSession, nil)

		mockUserRepository := mocks.NewUserRepositoryMock(t)
		mockUserRepository.EXPECT().
			GetByID(ctx, userID).
			Return(expectedUser, nil)

		authService := &AuthServiceImpl{
			sessionRepository: mockSessionRepository,
			userRepository:    mockUserRepository,
		}

		// Act
		user, err := authService.GetSessionUser(ctx, sessionID)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
	})
}
