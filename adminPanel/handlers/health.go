package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"adminPanel/database"
	"adminPanel/models"
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

// HealthCheck godoc
// @Summary Health check
// @Description Проверка работоспособности сервиса
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(models.HealthResponse{
		Status:   "healthy",
		Version:  "1.0.0",
	})
}

// DBHealthCheck godoc
// @Summary Database health check
// @Description Проверка подключения к базе данных
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Failure 503 {object} models.ErrorResponse
// @Router /health/db [get]
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