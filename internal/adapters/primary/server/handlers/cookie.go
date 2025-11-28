package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/labstack/echo/v4"
)

var (
	ErrCookieNotFound   = errors.New("cookie not found in header")
	ErrInvalidSignature = errors.New("invalid cookie signature")
)

type CookieHandler struct {
	cookieService services.CookieService
	cookieOptions config.CookieOptions
}

func NewCookieHandler(cookieService services.CookieService, config *config.Config) *CookieHandler {
	return &CookieHandler{
		cookieService: cookieService,
		cookieOptions: config.Session.CookieOptions,
	}
}

func (h *CookieHandler) Get(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie(h.cookieOptions.Name)
	if err != nil || cookie == nil {
		return nil, ErrCookieNotFound
	}

	if cookie.Value == "" {
		return nil, ErrCookieNotFound
	}

	value, err := h.cookieService.VerifySessionCookie(c.Request().Context(), cookie.Value)
	if err != nil {
		return nil, ErrInvalidSignature
	}

	cookie.Value = value
	return cookie, nil
}

func (h *CookieHandler) Set(c echo.Context, value string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	signedValue := h.cookieService.SignSessionCookie(c.Request().Context(), value)

	cookie := &http.Cookie{
		Name:     h.cookieOptions.Name,
		Value:    signedValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieOptions.Secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	}

	c.SetCookie(cookie)
}

func (h *CookieHandler) Clear(c echo.Context) {
	cookie := &http.Cookie{
		Name:     h.cookieOptions.Name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   h.cookieOptions.Secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	c.SetCookie(cookie)
}
