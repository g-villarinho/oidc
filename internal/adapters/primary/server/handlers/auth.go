package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/models"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/response"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService   *services.AuthService
	cookieHandler *CookieHandler
	logger        *slog.Logger
}

func NewAuthHandler(authService *services.AuthService, cookieHandler *CookieHandler, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		cookieHandler: cookieHandler,
		logger:        logger,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	logger := h.logger.With("handler", "Login")

	var payload models.LoginPayload
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind login payload", "error", err)
		return response.InvalidBind(c)
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("invalid login payload", "error", err)
		return response.ValidationError(c, err)
	}

	session, user, err := h.authService.Login(c.Request().Context(), payload.Email, payload.Password)
	if err != nil {
		if errors.Is(err, domain.ErrPasswordMismatch) || errors.Is(err, domain.ErrUserNotFound) {
			logger.Warn("invalid login attempt", "error", err)
			return response.Unauthorized(c, "INVALID_CREDENTIALS", "Invalid email or password. Please try again.")
		}

		if errors.Is(err, domain.ErrEmailNotVerified) {
			logger.Warn("login attempt with unverified email", "error", err)
			return response.Forbidden(c, "EMAIL_NOT_VERIFIED", "Email address has not been verified.")
		}

		logger.Error("failed to login user due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to login")
	}

	h.cookieHandler.Set(c, session.ID.String(), session.ExpiresAt)

	response := models.LoginResponse{
		User: models.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RegisterUser(c echo.Context) error {
	logger := h.logger.With("handler", "RegisterUser")

	var payload models.RegisterPayload
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind register payload", "error", err)
		return response.InvalidBind(c)
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("invalid register payload", "error", err)
		return response.ValidationError(c, err)
	}

	if err := h.authService.RegisterUser(c.Request().Context(), payload.Name, payload.Email, payload.Password); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			logger.Warn("attempt to register with an email that is already in use", "error", err)
			return response.BadRequest(c, "EMAIL_ALREADY_IN_USE", "The provided email is already in use. Please use a different email.")
		}

		logger.Error("failed to register user due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to register user")
	}

	return c.NoContent(http.StatusCreated)
}
