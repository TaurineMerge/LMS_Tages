package handlers

import (
	"context"

	"adminPanel/database"
	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler - обработчик для health check
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler создает обработчик health check
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *HealthHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/health", h.HealthCheck)
	router.Get("/health/db", h.DBHealthCheck)
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(models.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	})
}

func (h *HealthHandler) DBHealthCheck(c *fiber.Ctx) error {
	ctx := context.Background()
	err := h.db.Pool.Ping(ctx)

	if err != nil {
		return c.Status(503).JSON(models.HealthResponse{
			Status:   "unhealthy",
			Database: "disconnected",
			Version:  "1.0.0",
		})
	}

	return c.JSON(models.HealthResponse{
		Status:   "healthy",
		Database: "connected",
		Version:  "1.0.0",
	})
}
