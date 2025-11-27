package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewHealthHandler(db *pgxpool.Pool, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: redis,
	}
}

type HealthResponse struct {
	Status string `json:"status"`
}

type ReadinessResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func (h *HealthHandler) Liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: "ok",
	})
}

func (h *HealthHandler) Readiness(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	services := make(map[string]string)
	allHealthy := true

	if err := h.db.Ping(ctx); err != nil {
		services["postgres"] = "unhealthy"
		allHealthy = false
	} else {
		services["postgres"] = "healthy"
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		services["redis"] = "unhealthy"
		allHealthy = false
	} else {
		services["redis"] = "healthy"
	}

	status := "ok"
	httpStatus := http.StatusOK
	if !allHealthy {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	return c.JSON(httpStatus, ReadinessResponse{
		Status:   status,
		Services: services,
	})
}
