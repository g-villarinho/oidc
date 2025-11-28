//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/models"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogin_Endpoint tests the POST /v1/auth/login endpoint
func TestLogin_Endpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Teardown(t)

	server := SetupTestHTTPServer(t, env)

	// Setup handler
	cookieService := services.NewCookieService(server.Config)
	cookieHandler := handlers.NewCookieHandler(cookieService, server.Config)
	authHandler := handlers.NewAuthHandler(server.Services.AuthService, cookieHandler, server.Logger)

	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		env.Reset(t)

		// Create a verified user
		password := "SecurePassword123!"
		passwordHash := MustHashPassword(t, password)

		user := NewTestUser().
			WithEmail("john@example.com").
			WithName("John Doe").
			WithPasswordHash(passwordHash).
			WithEmailVerified(true).
			Build()

		MustCreateUser(t, env.DB, user)

		// Create request payload
		payload := models.LoginPayload{
			Email:    "john@example.com",
			Password: password,
			Continue: "https://app.example.com/callback", // Use full URL (will be rejected and fallback to /)
		}

		body, _ := json.Marshal(payload)

		// Create HTTP request
		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		// Execute handler
		err := authHandler.Login(c)

		// Assertions
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse response
		var response models.LoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, user.Email, response.User.Email)
		assert.Equal(t, user.Name, response.User.Name)
		assert.Equal(t, user.ID.String(), response.User.ID)
		// External URLs without allowedHosts fallback to "/"
		assert.Equal(t, "/", response.Continue)

		// Verify cookie was set
		cookies := rec.Result().Cookies()
		assert.NotEmpty(t, cookies)

		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == server.Config.Session.CookieOptions.Name {
				sessionCookie = cookie
				break
			}
		}
		assert.NotNil(t, sessionCookie)
		assert.NotEmpty(t, sessionCookie.Value)
		assert.True(t, sessionCookie.HttpOnly)
	})

	t.Run("should fail with invalid email", func(t *testing.T) {
		env.Reset(t)

		payload := models.LoginPayload{
			Email:    "nonexistent@example.com",
			Password: "anypassword",
			Continue: "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var errorResp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "INVALID_CREDENTIALS", errorResp["code"])
	})

	t.Run("should fail with invalid password", func(t *testing.T) {
		env.Reset(t)

		// Create a verified user
		password := "CorrectPassword123!"
		passwordHash := MustHashPassword(t, password)

		user := NewTestUser().
			WithEmail("jane@example.com").
			WithPasswordHash(passwordHash).
			WithEmailVerified(true).
			Build()

		MustCreateUser(t, env.DB, user)

		payload := models.LoginPayload{
			Email:    user.Email,
			Password: "WrongPassword456!",
			Continue: "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var errorResp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "INVALID_CREDENTIALS", errorResp["code"])
	})

	t.Run("should fail with unverified email", func(t *testing.T) {
		env.Reset(t)

		// Create an unverified user
		password := "SecurePassword123!"
		passwordHash := MustHashPassword(t, password)

		user := NewTestUser().
			WithEmail("unverified@example.com").
			WithPasswordHash(passwordHash).
			WithEmailVerified(false).
			Build()

		MustCreateUser(t, env.DB, user)

		payload := models.LoginPayload{
			Email:    user.Email,
			Password: password,
			Continue: "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		var errorResp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "EMAIL_NOT_VERIFIED", errorResp["code"])
	})

	t.Run("should fail with missing email", func(t *testing.T) {
		env.Reset(t)

		payload := models.LoginPayload{
			Password: "password123",
			Continue: "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"email":    "not-an-email",
			"password": "password123",
			"continue": "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"email":    "test@example.com",
			"password": "short",
			"continue": "https://app.example.com/callback",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should sanitize redirect URL", func(t *testing.T) {
		env.Reset(t)

		// Create a verified user
		password := "SecurePassword123!"
		passwordHash := MustHashPassword(t, password)

		user := NewTestUser().
			WithEmail("redirect@example.com").
			WithPasswordHash(passwordHash).
			WithEmailVerified(true).
			Build()

		MustCreateUser(t, env.DB, user)

		payload := models.LoginPayload{
			Email:    user.Email,
			Password: password,
			Continue: "javascript:alert('xss')", // Invalid/dangerous redirect
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.LoginResponse
		json.Unmarshal(rec.Body.Bytes(), &response)

		// Should fallback to safe default
		assert.Equal(t, "/", response.Continue)
	})

	t.Run("should fail with malformed JSON", func(t *testing.T) {
		env.Reset(t)

		body := []byte(`{"email": "test@example.com", "password": `)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/login", body)

		err := authHandler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}

// TestRegister_Endpoint tests the POST /v1/auth/register endpoint
func TestRegister_Endpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Teardown(t)

	server := SetupTestHTTPServer(t, env)

	// Setup handler
	cookieService := services.NewCookieService(server.Config)
	cookieHandler := handlers.NewCookieHandler(cookieService, server.Config)
	authHandler := handlers.NewAuthHandler(server.Services.AuthService, cookieHandler, server.Logger)

	t.Run("should register user successfully with valid data", func(t *testing.T) {
		env.Reset(t)

		payload := models.RegisterPayload{
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			Password: "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("should fail with duplicate email", func(t *testing.T) {
		env.Reset(t)

		// Create a user first
		user := NewTestUser().
			WithEmail("duplicate@example.com").
			WithPasswordHash(MustHashPassword(t, "password123")).
			Build()

		MustCreateUser(t, env.DB, user)

		// Try to register with same email
		payload := models.RegisterPayload{
			Name:     "Another User",
			Email:    "duplicate@example.com",
			Password: "DifferentPassword456!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		var errorResp map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.Equal(t, "EMAIL_ALREADY_IN_USE", errorResp["code"])
	})

	t.Run("should fail with missing name", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"email":    "test@example.com",
			"password": "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with missing email", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"name":     "Test User",
			"password": "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with invalid email format", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"name":     "Test User",
			"email":    "not-an-email",
			"password": "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with short password", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "short",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with weak password", func(t *testing.T) {
		env.Reset(t)

		payload := map[string]interface{}{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123", // Weak password (no uppercase, no special chars)
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with malformed JSON", func(t *testing.T) {
		env.Reset(t)

		body := []byte(`{"name": "Test User", "email": `)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should fail with email containing whitespace", func(t *testing.T) {
		env.Reset(t)

		payload := models.RegisterPayload{
			Name:     "Whitespace Test",
			Email:    "  whitespace@example.com  ",
			Password: "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("should handle very long name", func(t *testing.T) {
		env.Reset(t)

		longName := ""
		for i := 0; i < 300; i++ {
			longName += "a"
		}

		payload := models.RegisterPayload{
			Name:     longName,
			Email:    "longname@example.com",
			Password: "StrongPassword123!",
		}

		body, _ := json.Marshal(payload)

		c, rec := MakeRequest(http.MethodPost, "/v1/auth/register", body)

		err := authHandler.RegisterUser(c)

		require.NoError(t, err)
		assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusUnprocessableEntity)
	})
}
