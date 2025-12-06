package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Health check
func healthCheck(c *fiber.Ctx) error {
	ctx := context.Background()
	err := dbPool.Ping(ctx)
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
	}

	return c.JSON(toHealthResponseDTO(HealthResponse{
		Status:   "healthy",
		Database: dbStatus,
		Version:  "1.0.0",
	}))
}
