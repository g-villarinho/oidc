package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/models"
	"github.com/g-villarinho/oidc-server/internal/adapters/primary/server/response"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/g-villarinho/oidc-server/internal/core/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	clientService *services.ClientService
	logger        *slog.Logger
}

func NewClientHandler(clientService *services.ClientService, logger *slog.Logger) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		logger:        logger,
	}
}

func (h *ClientHandler) CreateClient(c echo.Context) error {
	logger := h.logger.With("handler", "CreateClient")

	var payload models.CreateClientRequest
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind create client payload", "error", err)
		return response.InvalidBind(c)
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("invalid create client payload", "error", err)
		return response.ValidationError(c, err)
	}

	client, err := h.clientService.CreateClient(
		c.Request().Context(),
		payload.ClientID,
		payload.ClientSecret,
		payload.ClientName,
		payload.RedirectURIs,
		payload.GrantTypes,
		payload.ResponseTypes,
		payload.Scope,
		payload.LogoURL,
	)
	if err != nil {
		if errors.Is(err, ports.ErrUniqueKeyViolation) {
			logger.Warn("attempt to create client with duplicate client_id", "error", err)
			return response.ConflictError(c, "CLIENT_ALREADY_EXISTS", "A client with this client_id already exists")
		}

		logger.Error("failed to create client due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to create client")
	}

	clientResponse := models.ClientResponse{
		ID:            client.ID.String(),
		ClientID:      client.ClientID,
		ClientName:    client.ClientName,
		RedirectURIs:  client.RedirectURIs,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scope:         client.Scope,
		LogoURL:       client.LogoURL,
		CreatedAt:     client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, clientResponse)
}

func (h *ClientHandler) GetClientByID(c echo.Context) error {
	logger := h.logger.With("handler", "GetClientByID")

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warn("invalid client ID format", "id", idParam, "error", err)
		return response.BadRequest(c, "INVALID_CLIENT_ID", "Invalid client ID format")
	}

	client, err := h.clientService.GetClientByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("client not found", "id", id)
			return response.NotFound(c, "CLIENT_NOT_FOUND", "Client not found")
		}

		logger.Error("failed to get client due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to get client")
	}

	clientResponse := models.ClientResponse{
		ID:            client.ID.String(),
		ClientID:      client.ClientID,
		ClientName:    client.ClientName,
		RedirectURIs:  client.RedirectURIs,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scope:         client.Scope,
		LogoURL:       client.LogoURL,
		CreatedAt:     client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, clientResponse)
}

func (h *ClientHandler) ListClients(c echo.Context) error {
	logger := h.logger.With("handler", "ListClients")

	clients, err := h.clientService.ListClients(c.Request().Context())
	if err != nil {
		logger.Error("failed to list clients due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to list clients")
	}

	clientResponses := make([]models.ClientResponse, 0, len(clients))
	for _, client := range clients {
		clientResponses = append(clientResponses, models.ClientResponse{
			ID:            client.ID.String(),
			ClientID:      client.ClientID,
			ClientName:    client.ClientName,
			RedirectURIs:  client.RedirectURIs,
			GrantTypes:    client.GrantTypes,
			ResponseTypes: client.ResponseTypes,
			Scope:         client.Scope,
			LogoURL:       client.LogoURL,
			CreatedAt:     client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:     client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	response := models.ClientListResponse{
		Clients: clientResponses,
		Total:   len(clientResponses),
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) UpdateClient(c echo.Context) error {
	logger := h.logger.With("handler", "UpdateClient")

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warn("invalid client ID format", "id", idParam, "error", err)
		return response.BadRequest(c, "INVALID_CLIENT_ID", "Invalid client ID format")
	}

	var payload models.UpdateClientRequest
	if err := c.Bind(&payload); err != nil {
		logger.Error("failed to bind update client payload", "error", err)
		return response.InvalidBind(c)
	}

	if err := c.Validate(&payload); err != nil {
		logger.Error("invalid update client payload", "error", err)
		return response.ValidationError(c, err)
	}

	client, err := h.clientService.UpdateClient(
		c.Request().Context(),
		id,
		payload.ClientName,
		payload.RedirectURIs,
		payload.GrantTypes,
		payload.ResponseTypes,
		payload.Scope,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn("client not found for update", "id", id)
			return response.NotFound(c, "CLIENT_NOT_FOUND", "Client not found")
		}

		logger.Error("failed to update client due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to update client")
	}

	clientResponse := models.ClientResponse{
		ID:            client.ID.String(),
		ClientID:      client.ClientID,
		ClientName:    client.ClientName,
		RedirectURIs:  client.RedirectURIs,
		GrantTypes:    client.GrantTypes,
		ResponseTypes: client.ResponseTypes,
		Scope:         client.Scope,
		LogoURL:       client.LogoURL,
		CreatedAt:     client.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     client.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusOK, clientResponse)
}

func (h *ClientHandler) DeleteClient(c echo.Context) error {
	logger := h.logger.With("handler", "DeleteClient")

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		logger.Warn("invalid client ID format", "id", idParam, "error", err)
		return response.BadRequest(c, "INVALID_CLIENT_ID", "Invalid client ID format")
	}

	err = h.clientService.DeleteClient(c.Request().Context(), id)
	if err != nil {
		logger.Error("failed to delete client due to internal error", "error", err)
		return response.InternalServerError(c, "Failed to delete client")
	}

	return c.NoContent(http.StatusNoContent)
}
