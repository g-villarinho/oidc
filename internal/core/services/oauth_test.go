package services

import (
	"context"
	"errors"
	"testing"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestValidateAuthorizationClient(t *testing.T) {
	t.Run("should validate client successfully when all parameters are valid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		responseType := "code"
		scopes := []string{"openid", "profile"}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: responseType,
			Scopes:       scopes,
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI, "https://example.com/other"},
			ResponseTypes: []string{"code", "token"},
			Scopes:        []string{"openid", "profile", "email"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, expectedClient.ClientID, client.ClientID)
		assert.Equal(t, expectedClient.ClientName, client.ClientName)
	})

	t.Run("should return error when client not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "nonexistent-client"
		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  "https://example.com/callback",
			ResponseType: "code",
			Scopes:       []string{"openid"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(nil, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.ErrorIs(t, err, domain.ErrClientNotFound)
	})

	t.Run("should return error when client repository fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  "https://example.com/callback",
			ResponseType: "code",
			Scopes:       []string{"openid"},
		}

		expectedError := errors.New("database connection error")

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(nil, expectedError)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "validate authorization client")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should return error when redirect URI is not registered", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		invalidRedirectURI := "https://malicious.com/callback"

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  invalidRedirectURI,
			ResponseType: "code",
			Scopes:       []string{"openid"},
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{"https://example.com/callback"},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.ErrorIs(t, err, domain.ErrInvalidRedirectURI)
	})

	t.Run("should return error when response type is not supported", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		unsupportedResponseType := "token"

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: unsupportedResponseType,
			Scopes:       []string{"openid"},
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.ErrorIs(t, err, domain.ErrUnsupportedResponseType)
	})

	t.Run("should return error when requested scopes are not supported", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		unsupportedScopes := []string{"openid", "admin"}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       unsupportedScopes,
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.Error(t, err)
		assert.Nil(t, client)
		assert.ErrorIs(t, err, domain.ErrInvalidScope)
	})

	t.Run("should validate client successfully when requesting subset of scopes", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid"}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       scopes,
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile", "email"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, expectedClient.ClientID, client.ClientID)
	})

	t.Run("should validate client successfully with multiple response types", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       []string{"openid"},
		}

		expectedClient := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code", "token", "id_token"},
			Scopes:        []string{"openid", "profile", "email"},
		}

		mockClientRepository := mocks.NewClientRepositoryMock(t)
		mockClientRepository.EXPECT().
			GetByClientID(ctx, clientID).
			Return(expectedClient, nil)

		authService := &AuthorizationServiceImpl{
			clientRepository: mockClientRepository,
		}

		// Act
		client, err := authService.ValidateAuthorizationClient(ctx, params)

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, expectedClient.ClientID, client.ClientID)
	})
}

func TestAuthorize(t *testing.T) {
	t.Run("should create authorization code successfully when valid parameters are provided", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid", "profile"}
		codeChallenge := "test-challenge"
		codeChallengeMethod := "S256"

		client := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile", "email"},
		}

		params := domain.AuthorizeParams{
			ClientID:            clientID,
			RedirectURI:         redirectURI,
			ResponseType:        "code",
			Scopes:              scopes,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
		}

		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Return(nil)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, code)
	})

	t.Run("should return error when authorization code repository fails to create code", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid", "profile"}

		client := &domain.Client{
			ID:            uuid.New(),
			ClientID:      clientID,
			ClientSecret:  "secret",
			ClientName:    "Test Client",
			RedirectURIs:  []string{redirectURI},
			ResponseTypes: []string{"code"},
			Scopes:        []string{"openid", "profile", "email"},
		}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       scopes,
		}

		expectedError := errors.New("database connection error")

		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Return(expectedError)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.Error(t, err)
		assert.Empty(t, code)
		assert.Contains(t, err.Error(), "save authorization code")
		assert.Contains(t, err.Error(), expectedError.Error())
	})

	t.Run("should create authorization code with PKCE parameters", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid"}
		codeChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
		codeChallengeMethod := "S256"

		client := &domain.Client{
			ID:           uuid.New(),
			ClientID:     clientID,
			ClientSecret: "secret",
			ClientName:   "Test Client",
			RedirectURIs: []string{redirectURI},
			Scopes:       []string{"openid", "profile"},
		}

		params := domain.AuthorizeParams{
			ClientID:            clientID,
			RedirectURI:         redirectURI,
			ResponseType:        "code",
			Scopes:              scopes,
			CodeChallenge:       codeChallenge,
			CodeChallengeMethod: codeChallengeMethod,
		}

		var capturedAuthCode *domain.AuthorizationCode
		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Run(func(ctx context.Context, authCode *domain.AuthorizationCode) {
				capturedAuthCode = authCode
			}).
			Return(nil)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.NotNil(t, capturedAuthCode)
		assert.Equal(t, codeChallenge, capturedAuthCode.CodeChallenge)
		assert.Equal(t, codeChallengeMethod, capturedAuthCode.CodeChallengeMethod)
		assert.Equal(t, userID, capturedAuthCode.UserID)
		assert.Equal(t, clientID, capturedAuthCode.ClientID)
		assert.Equal(t, redirectURI, capturedAuthCode.RedirectURI)
		assert.Equal(t, scopes, capturedAuthCode.Scopes)
	})

	t.Run("should create authorization code without PKCE parameters", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid", "profile"}

		client := &domain.Client{
			ID:           uuid.New(),
			ClientID:     clientID,
			ClientSecret: "secret",
			ClientName:   "Test Client",
			RedirectURIs: []string{redirectURI},
			Scopes:       []string{"openid", "profile"},
		}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       scopes,
		}

		var capturedAuthCode *domain.AuthorizationCode
		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Run(func(ctx context.Context, authCode *domain.AuthorizationCode) {
				capturedAuthCode = authCode
			}).
			Return(nil)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.NotNil(t, capturedAuthCode)
		assert.Empty(t, capturedAuthCode.CodeChallenge)
		assert.Empty(t, capturedAuthCode.CodeChallengeMethod)
		assert.Equal(t, userID, capturedAuthCode.UserID)
		assert.Equal(t, clientID, capturedAuthCode.ClientID)
	})

	t.Run("should create authorization code with multiple scopes", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid", "profile", "email"}

		client := &domain.Client{
			ID:           uuid.New(),
			ClientID:     clientID,
			ClientSecret: "secret",
			ClientName:   "Test Client",
			RedirectURIs: []string{redirectURI},
			Scopes:       []string{"openid", "profile", "email"},
		}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       scopes,
		}

		var capturedAuthCode *domain.AuthorizationCode
		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Run(func(ctx context.Context, authCode *domain.AuthorizationCode) {
				capturedAuthCode = authCode
			}).
			Return(nil)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.NotNil(t, capturedAuthCode)
		assert.Equal(t, scopes, capturedAuthCode.Scopes)
		assert.Len(t, capturedAuthCode.Scopes, 3)
	})

	t.Run("should return same code as created authorization code", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		userID := uuid.New()
		clientID := "test-client-id"
		redirectURI := "https://example.com/callback"
		scopes := []string{"openid"}

		client := &domain.Client{
			ID:           uuid.New(),
			ClientID:     clientID,
			ClientSecret: "secret",
			ClientName:   "Test Client",
			RedirectURIs: []string{redirectURI},
			Scopes:       []string{"openid"},
		}

		params := domain.AuthorizeParams{
			ClientID:     clientID,
			RedirectURI:  redirectURI,
			ResponseType: "code",
			Scopes:       scopes,
		}

		var capturedAuthCode *domain.AuthorizationCode
		mockAuthorizationCodeRepository := mocks.NewAuthorizationCodeRepositoryMock(t)
		mockAuthorizationCodeRepository.EXPECT().
			Create(ctx, mock.AnythingOfType("*domain.AuthorizationCode")).
			Run(func(ctx context.Context, authCode *domain.AuthorizationCode) {
				capturedAuthCode = authCode
			}).
			Return(nil)

		authService := &AuthorizationServiceImpl{
			authorizationCodeRepository: mockAuthorizationCodeRepository,
		}

		// Act
		code, err := authService.Authorize(ctx, userID, client, params)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, code)
		assert.NotNil(t, capturedAuthCode)
		assert.Equal(t, capturedAuthCode.Code, code)
	})
}