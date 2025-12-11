package handlers

import (
	"context"

	"adminPanel/database"
	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

// HealthCheck returns application health status.
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.HealthCheck.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
		))

	return c.JSON(models.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	})
}

// DBHealthCheck verifies database connectivity status.
func (h *HealthHandler) DBHealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.DBHealthCheck.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
		))

	if ctx == nil {
		ctx = context.Background()
	}
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
