package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// ============ SWAGGER SETUP ============

func setupSwagger(app *fiber.App) {
	// 1. Endpoint для вашего swagger.json
	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			log.Printf("❌ Ошибка чтения swagger.json: %v", err)
			return c.Status(500).JSON(newInternalError("Failed to load Swagger documentation"))
		}

		var swaggerSpec map[string]interface{}
		if err := json.Unmarshal(data, &swaggerSpec); err != nil {
			log.Printf("❌ Ошибка парсинга swagger.json: %v", err)
			return c.Status(500).JSON(newInternalError("Invalid Swagger JSON"))
		}

		// Обновляем basePath и host для вашего API
		swaggerSpec["basePath"] = "/api/v1"
		swaggerSpec["host"] = "localhost:4000"

		return c.JSON(swaggerSpec)
	})

	// 2. Swagger UI от библиотеки (будет использовать наш /swagger.json)
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger.json",
		DeepLinking:  true,
		DocExpansion: "list",
	}))
}
