// Package handlers содержит HTTP-обработчики для всех маршрутов приложения.
// Этот пакет предоставляет обработчики для health check, категорий, курсов и уроков.
package handlers

import (
	"context"

	"adminPanel/database"
	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HealthHandler - обработчик для health check запросов.
// Предоставляет информацию о состоянии приложения и базы данных.
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler создает новый обработчик health check.
//
// Параметры:
//   - db: указатель на соединение с базой данных
//
// Возвращает:
//   - *HealthHandler: указатель на новый экземпляр HealthHandler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// RegisterRoutes регистрирует маршруты для health check.
//
// Параметры:
//   - router: маршрутизатор Fiber для регистрации маршрутов
func (h *HealthHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/health", h.HealthCheck)
	router.Get("/health/db", h.DBHealthCheck)
}

// HealthCheck возвращает статус приложения.
// Этот метод не требует аутентификации и доступен для всех.
//
// HTTP маршрут: GET /health
//
// Возвращает:
//   - JSON с общим статусом приложения и версией
//   - HTTP статус 200 OK
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	ctx := c.UserContext()
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

// DBHealthCheck проверяет состояние подключения к базе данных.
// Этот метод не требует аутентификации и доступен для всех.
//
// HTTP маршрут: GET /health/db
//
// Возвращает:
//   - JSON со статусом приложения и состоянием базы данных
//   - HTTP статус 200 OK, если база данных доступна
//   - HTTP статус 503 Service Unavailable, если база данных недоступна
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
