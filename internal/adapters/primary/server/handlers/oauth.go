package handlers

import (
	"log/slog"
	"net/http"
	"net/url"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/context"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/models"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/response"
	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/g-villarinho/oidc-server/pkg/oauth"
	"github.com/labstack/echo/v4"
)

type OAuthHandler struct {
	oauthService services.OAuthService
	context      *context.EchoContext
	logger       *slog.Logger
	url          config.URL
}

func NewOAuthHandler(
	oauthService services.OAuthService,
	context *context.EchoContext,
	logger *slog.Logger,
	config *config.Config,
) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		context:      context,
		logger:       logger.With("handler", "authorization"),
		url:          config.URL,
	}
}

func (h *OAuthHandler) Authorize(c echo.Context) error {
	logger := h.logger.With("method", "Authorize")

	session := h.context.GetSession(c)
	if session != nil {
		logger = logger.With("user_id", session.UserID)
	}

	var payload models.AuthorizePayload
	if err := c.Bind(&payload); err != nil {
		logger.Error("error to bind authorize payload", "error", err)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("validate authorize payload", "error", err)
		// TODO: redirecionar para redirect_uri com erro em vez de 400
		// seguindo spec OAuth2 (se redirect_uri for válido)
		return c.String(http.StatusBadRequest, "invalid params authorize")
	}

	if err := h.oauthService.VerifyAuthorization(c.Request().Context(), payload.ToAuthorizeParams()); err != nil {
		logger.Error("error to verify authorization", "error", err)
		// TODO: tratar erros específicos de redirect_uri diferente
		// se redirect_uri for inválido, não pode redirecionar com erro
		return c.String(http.StatusBadRequest, "invalid client authorization")
	}

	if session == nil {
		logger.Info("no active session, redirecting to login")

		continueURL := oauth.GenerateContinueURL(h.url.APIBaseURL, payload.ToContinueURLParams())

		loginURL, err := url.Parse(h.url.AppBaseURL)
		if err != nil {
			logger.Error("error to parse app base URL", "error", err)
			return response.InternalServerError(c, "The authorization workflow could not be completed due to an internal error.")
		}

		loginURL.Path = "/login"
		q := loginURL.Query()
		q.Set("continue", continueURL)
		loginURL.RawQuery = q.Encode()

		return c.Redirect(http.StatusFound, loginURL.String())
	}

	authorizationCode, err := h.oauthService.CreateAuthorizationCode(c.Request().Context(), session.UserID, payload.ToAuthorizeParams())
	if err != nil {
		logger.Error("error to authorize client", "error", err)
		return response.InternalServerError(c, "The authorization workflow could not be completed due to an internal error.")
	}

	redirectURI := oauth.GenerateCallbackURL(payload.RedirectURI, authorizationCode.Code, payload.State)

	return c.Redirect(http.StatusFound, redirectURI)
}

func (h *OAuthHandler) Token(c echo.Context) error {
	logger := h.logger.With("method", "Token")

	logger.Info("Token endpoint called")
	return c.String(200, "Token issued successfully")
}
