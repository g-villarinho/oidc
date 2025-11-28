package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/middlewares"
	"github.com/g-villarinho/oidc-server/internal/config"
	"github.com/g-villarinho/oidc-server/pkg/serializer"
	"github.com/g-villarinho/oidc-server/pkg/validation"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"
)

type ServerParams struct {
	dig.In

	Config         *config.Config
	AuthHandler    *handlers.AuthHandler
	ClientHandler  *handlers.ClientHandler
	HealthHandler  *handlers.HealthHandler
	OAuthHandler   *handlers.OAuthHandler
	AuthMiddleware *middlewares.AuthMiddleware
}

type Server struct {
	echo            *echo.Echo
	port            int
	shutdownTimeout time.Duration
}

func NewServer(params ServerParams) *Server {
	e := echo.New()
	port := params.Config.Server.Port
	e.Validator = validation.NewValidator()
	e.JSONSerializer = serializer.NewJSONSerializer()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.BodyLimit("10M"))
	e.Use(middlewares.Cors(params.Config))
	e.Use(middlewares.RateLimiter(&params.Config.RateLimit))

	group := e.Group("/api")
	registerAuthRoutes(group, params.AuthHandler)
	registerClientRoutes(group, params.ClientHandler)
	registerHealthRoutes(e, params.HealthHandler)
	registerOAuthRoutes(group, params.OAuthHandler, params.AuthMiddleware)

	return &Server{
		echo:            e,
		port:            port,
		shutdownTimeout: params.Config.Server.ShutdownTimeout,
	}
}

func (s *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		address := fmt.Sprintf(":%d", s.port)
		if err := s.echo.Start(address); err != nil {
			s.echo.Logger.Info("Shutting down the server")
		}
	}()

	<-quit
	s.echo.Logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.echo.Logger.Info("Server exited gracefully")
	return nil
}
