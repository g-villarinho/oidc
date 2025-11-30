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
	"github.com/g-villarinho/oidc-server/pkg/oauth"
	"github.com/labstack/echo/v4"
)

type OAuthHandler struct {
	service services.OAuthService
	context *context.EchoContext
	logger  *slog.Logger
	url     config.URL
}

func NewOAuthHandler(
	service services.OAuthService,
	context *context.EchoContext,
	logger *slog.Logger,
	config *config.Config,
) *OAuthHandler {
	return &OAuthHandler{
		service: service,
		context: context,
		logger:  logger.With("handler", "authorization"),
		url:     config.URL,
	}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	logger := h.logger.With("method", "Authorize")

	var payload models.AuthorizePayload
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind authorize payload", "error", err)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("validate authorize payload", "error", err)
		// TODO: redirecionar para redirect_uri com erro em vez de 400
		// seguindo spec OAuth2 (se redirect_uri for válido)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	if err := h.service.VerifyAuthorization(c.Request().Context(), payload.ToAuthorizeParams()); err != nil {
		logger.Error("failed to validate authorization client", "error", err)
		// TODO: tratar erros específicos de redirect_uri diferente
		// se redirect_uri for inválido, não pode redirecionar com erro
		return c.String(http.StatusBadRequest, "invalid client authorization")
	}

	session := h.context.GetSession(c)
	if session == nil {
		logger.Info("no active session, redirecting to login")

		continueURLParams := models.ToContinueURLParams(payload)

		continueURL := oauth.GenerateContinueURL(h.url.APIBaseURL, continueURLParams)
		loginParams := url.Values{}
		loginParams.Set("continue", continueURL)

		loginURL := fmt.Sprintf("%s/login?%s", h.url.AppBaseURL, loginParams.Encode())

		return c.Redirect(http.StatusFound, loginURL)
	}

	code, err := h.service.Authorize(c.Request().Context(), session.UserID, payload.ToAuthorizeParams())
	if err != nil {
		logger.Error("authorize client", "error", err)
		return c.String(http.StatusInternalServerError, "failed to authorize client")
	}

	redirectParams := url.Values{}
	redirectParams.Set("code", code)
	if payload.State != "" {
		redirectParams.Set("state", payload.State)
	}

	redirectURI := fmt.Sprintf("%s?%s", payload.RedirectURI, redirectParams.Encode())

	return c.Redirect(http.StatusFound, redirectURI)
}

func (h *OAuthHandler) Token(c echo.Context) error {
	return c.String(200, "Token issued successfully")
}
