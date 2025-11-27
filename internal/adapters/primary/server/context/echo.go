package context

import (
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/labstack/echo/v4"
)

const (
	sessionKey = "session-key"
	userKey    = "user-key"
)

type EchoContext struct {
}

func NewEchoContext() *EchoContext {
	return &EchoContext{}
}

func (c *EchoContext) SetSession(ectx echo.Context, session domain.Session) {
	ectx.Set(sessionKey, session)
}

func (c *EchoContext) GetSession(ectx echo.Context) *domain.Session {
	session, ok := ectx.Get(sessionKey).(domain.Session)
	if !ok {
		return nil
	}

	return &session
}

func (c *EchoContext) SetUser(ectx echo.Context, user domain.User) {
	ectx.Set(userKey, user)
}

func (c *EchoContext) GetUser(ectx echo.Context) *domain.User {
	user, ok := ectx.Get(userKey).(domain.User)
	if !ok {
		return nil
	}

	return &user
}
