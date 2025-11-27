package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/context"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/models"
	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/g-villarinho/oidc-server/pkg/security"
	"github.com/labstack/echo/v4"
)

type AuthorizationHandler struct {
	service *services.AuthorizationService
	context *context.EchoContext
	logger  *slog.Logger
	url     config.URL
}

func NewAuthorizationHandler(
	service *services.AuthorizationService,
	context *context.EchoContext,
	logger *slog.Logger,
	config *config.Config,
) *AuthorizationHandler {
	return &AuthorizationHandler{
		service: service,
		context: context,
		logger:  logger.With("handler", "authorization"),
		url:     config.URL,
	}
}

func (h *AuthorizationHandler) Authorize(c echo.Context) error {
	logger := h.logger.With("method", "Authorize")

	var payload models.AuthorizePayload
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind authorize payload", "error", err)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("validate authorize payload", "error", err)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	client, err := h.service.ValidateAuthorizationClient(c.Request().Context(), payload.ToAuthorizeParams())
	if err != nil {
		logger.Error("failed to validate authorization client", "error", err)
		return c.String(http.StatusBadRequest, "invalid client authorization")
	}

	session := h.context.GetSession(c)
	if session == nil {
		logger.Info("no active session, redirecting to login")

		continueURLParams := models.ToContinueURLParams(payload)

		continueURL := security.GenerateContinueURL(fmt.Sprintf("%s/authorize", h.url.AppBaseURL), continueURLParams)
		loginParams := url.Values{}
		loginParams.Set("continue", continueURL)

		loginURL := fmt.Sprintf("%s/login?%s", h.url.AppBaseURL, loginParams.Encode())

		return c.Redirect(http.StatusFound, loginURL)
	}

	code, err := h.service.Authorize(c.Request().Context(), session.UserID, client, payload.ToAuthorizeParams())
	if err != nil {
		logger.Error("authorize client", "error", err)
		return c.String(http.StatusInternalServerError, "failed to authorize client")
	}

	redirectParams := url.Values{}
	redirectParams.Set("code", code)
	redirectParams.Set("state", payload.State)

	redirectURI := fmt.Sprintf("%s?%s", payload.RedirectURI, redirectParams.Encode())

	return c.Redirect(http.StatusFound, redirectURI)
}

func (h *AuthorizationHandler) Token(c echo.Context) error {
	return c.String(200, "Token issued successfully")
}
