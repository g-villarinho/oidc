package server

import (
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/handlers"
	"github.com/labstack/echo/v4"
)

func registerClientRoutes(e *echo.Group, clientHandler *handlers.ClientHandler) {
	clientsV1Group := e.Group("/v1/clients")
	clientsV1Group.POST("", clientHandler.CreateClient)
	clientsV1Group.GET("", clientHandler.ListClients)
	clientsV1Group.GET("/:id", clientHandler.GetClientByID)
	clientsV1Group.PUT("/:id", clientHandler.UpdateClient)
	clientsV1Group.DELETE("/:id", clientHandler.DeleteClient)
}

func registerAuthRoutes(e *echo.Group, authHandler *handlers.AuthHandler) {
	authV1Group := e.Group("/v1/auth")
	authV1Group.POST("/login", authHandler.Login)
	authV1Group.POST("/register", authHandler.RegisterUser)
}
