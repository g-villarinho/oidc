package server

import (
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, clientHandler *handlers.ClientHandler, authHandler *handlers.AuthHandler) {
	// API v1 group
	api := e.Group("/api/v1")

	// Client routes
	clients := api.Group("/clients")
	clients.POST("", clientHandler.CreateClient)
	clients.GET("", clientHandler.ListClients)
	clients.GET("/:id", clientHandler.GetClientByID)
	clients.PUT("/:id", clientHandler.UpdateClient)
	clients.DELETE("/:id", clientHandler.DeleteClient)

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.RegisterUser)
}
