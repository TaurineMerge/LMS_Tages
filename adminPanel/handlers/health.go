package handlers

import (
	"context"

	"adminPanel/database"
	"adminPanel/handlers/dto/response"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HealthHandler обрабатывает HTTP-запросы для проверки здоровья приложения.
// Содержит соединение с базой данных для проверки доступности.
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler создает новый экземпляр HealthHandler.
// Принимает соединение с базой данных.
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// RegisterRoutes регистрирует маршруты для проверки здоровья.
// Регистрирует /health и /health/db.
func (h *HealthHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/health", h.HealthCheck)
	router.Get("/health/db", h.DBHealthCheck)
}

// HealthCheck обрабатывает GET /health.
// Возвращает статус здоровья приложения.
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.HealthCheck.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
		))

	return c.JSON(response.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	})
}

// DBHealthCheck обрабатывает GET /health/db.
// Проверяет подключение к базе данных и возвращает статус.
func (h *HealthHandler) DBHealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
		return c.Status(503).JSON(response.HealthResponse{
			Status:   "unhealthy",
			Database: "disconnected",
			Version:  "1.0.0",
		})
	}

	return c.JSON(response.HealthResponse{
		Status:   "healthy",
		Database: "connected",
		Version:  "1.0.0",
	})
}
