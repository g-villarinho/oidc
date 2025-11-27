package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

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

	session := h.context.GetSession(c)
	if session == nil {
		logger.Info("no active session, redirecting to login")

		continueURLParams := models.ToContinueURLParams(payload)

		continueURL := security.GenerateContinueURL(fmt.Sprintf("%s/authorize", h.url.AppBaseURL), continueURLParams)
		loginURL := fmt.Sprintf("%s/login?continue=%s", h.url.AppBaseURL, continueURL)

		return c.Redirect(http.StatusFound, loginURL)
	}

	return c.String(200, "Authorization successful")
}

func (h *AuthorizationHandler) Token(c echo.Context) error {
	return c.String(200, "Token issued successfully")
}
