// // docs/swagger.go
// package docs

// import (
// 	"embed"
// 	"io/fs"
// 	"log"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/swagger"
// )

// //go:embed swagger.json
// var swaggerJSON embed.FS

// // SetupSwagger настраивает Swagger с использованием статического файла
// func SetupSwagger(router fiber.Router) {  // Изменили тип параметра
// 	// Swagger UI
// 	router.Get("/swagger/*", swagger.New(swagger.Config{
// 		URL:         "/swagger/doc.json",
// 		DeepLinking: true,
// 	}))
	
// 	// OpenAPI спецификация из файла
// 	router.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
// 		data, err := fs.ReadFile(swaggerJSON, "swagger.json")
// 		if err != nil {
// 			log.Printf("Failed to read swagger.json: %v", err)
// 			return c.Status(500).JSON(fiber.Map{
// 				"error": "Failed to read API documentation",
// 			})
// 		}
		
// 		c.Set("Content-Type", "application/json")
// 		return c.Send(data)
// 	})
// }