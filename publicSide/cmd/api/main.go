package main

import (
	"log"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func main() {
	// Create a new Fiber app with a custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:9090",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Serve the `doc` folder so swagger.json is reachable at `/doc/swagger.json`
	app.Static("/doc", "./doc/swagger")

	apiV1 := app.Group("/api/v1")

	// Setup swagger (serve UI at /api/v1/swagger/* and point it to the served JSON)
	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	// Initialize repositories
	categoryRepo := repository.NewCategoryMemoryRepository()
	courseRepo := repository.NewCourseMemoryRepository()
	lessonRepo := repository.NewLessonMemoryRepository()

	// Initialize services
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo)
	lessonService := service.NewLessonService(lessonRepo)

	// Initialize handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	courseHandler := handler.NewCourseHandler(courseService)
	lessonHandler := handler.NewLessonHandler(lessonService)

	// Register routes
	handler.RegisterRoutes(apiV1, categoryHandler, courseHandler, lessonHandler)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
