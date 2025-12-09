package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("Error loading .env file, using environment variables", "error", err)
	}

	cfg, err := config.New(
		config.WithDBFromEnv(),
		config.WithCORSFromEnv(),
		config.WithPortFromEnv(),
	)

	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	dbPool, err := database.NewConnection(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	slog.Info("Database connection pool established")

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))

	app.Static("/doc", "./doc/swagger")

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	lessonRepo := repository.NewLessonRepository(dbPool)

	lessonService := service.NewLessonService(lessonRepo)

	lessonHandler := handler.NewLessonHandler(lessonService)


	categoryRouter := apiV1.Group("/categories")
	courseRouter := categoryRouter.Group(apiconst.PathCategory + "/courses")
	
	lessonHandler.RegisterRoutes(courseRouter)
	
	slog.Info("Starting server", "address", cfg.Port)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
