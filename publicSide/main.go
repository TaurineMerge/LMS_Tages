package main

import (
	"log"

	"github.com/TaurineMerge/LMS_Tages/publicSide/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/database"
	"github.com/TaurineMerge/LMS_Tages/publicSide/handler"
	"github.com/TaurineMerge/LMS_Tages/publicSide/handler/public"
	"github.com/TaurineMerge/LMS_Tages/publicSide/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	courseRepo := repository.NewCourseRepository(db)
	lessonRepo := repository.NewLessonRepository(db)

	// Initialize services
	courseService := service.NewCourseService(courseRepo)
	lessonService := service.NewLessonService(lessonRepo, courseRepo)

	// Initialize handlers
	courseHandler := public.NewCourseHandler(courseService)
	lessonHandler := public.NewLessonHandler(lessonService)
	swaggerHandler := handler.NewSwaggerHandler()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Public API v1.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	// Register routes
	courseHandler.RegisterRoutes(app)
	lessonHandler.RegisterRoutes(app)
	swaggerHandler.RegisterRoutes(app)

	// Root route - serve Swagger UI
	app.Get("/", swaggerHandler.ServeSwaggerUI)

	// Start server
	log.Println("Server starting on port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
