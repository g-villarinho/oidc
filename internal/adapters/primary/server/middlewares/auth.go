package middlewares

import (
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/context"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/response"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	cookieHandler     *handlers.CookieHandler
	sessionRepository ports.SessionRepository
	context           *context.EchoContext
}

func NewAuthMiddleware(cookieHandler *handlers.CookieHandler, sessionRepository ports.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		cookieHandler:     cookieHandler,
		sessionRepository: sessionRepository,
	}
}

func (m *AuthMiddleware) RequireAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := m.cookieHandler.Get(c)
		if err != nil {
			return response.Unauthorized(c, "TOKEN_MISSING", "You need to be logged int to access this resource")
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			return response.Unauthorized(c, "INVALID_TOKEN", "The provided token is invalid")
		}

		session, err := m.sessionRepository.GetByID(c.Request().Context(), sessionID)
		if err != nil || session == nil {
			m.cookieHandler.Clear(c)
			return response.Unauthorized(c, "INVALID_TOKEN", "The provided token is invalid")
		}

		m.context.SetSession(c, *session)
		return next(c)
	}
}

func (m *AuthMiddleware) OptionalAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := m.cookieHandler.Get(c)
		if err != nil {
			return next(c)
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			return next(c)
		}

		session, err := m.sessionRepository.GetByID(c.Request().Context(), sessionID)
		if err != nil || session == nil {
			m.cookieHandler.Clear(c)
			return next(c)
		}

		m.context.SetSession(c, *session)
		return next(c)
	}
}
