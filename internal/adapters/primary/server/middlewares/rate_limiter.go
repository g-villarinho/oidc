package middlewares

import (
	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func RateLimiter(config *config.RateLimit) echo.MiddlewareFunc {
	rate := rate.Limit(config.MaxRequests)
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate))
}
