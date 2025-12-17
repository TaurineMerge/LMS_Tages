// Package main is the entry point for the publicSide service.
// It initializes the configuration, database, router, and starts the HTTP server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/database"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/tracing"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file first to get all env variables including log level
	if err := godotenv.Load(); err != nil {
		slog.Warn("Not loaded .env file, using environment variables", "error", err)
	}

	// Initialize config with all options
	cfg, err := config.New(
		config.WithDBFromEnv(),
		config.WithCORSFromEnv(),
		config.WithPortFromEnv(),
		config.WithTracingFromEnv(),
		config.WithLogLevelFromEnv(),
	)
	if err != nil {
		// Use a temporary logger for this initial error, as the main one is not yet set up
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	// Initialize logger with the configured level
	var level slog.Level
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	// Initialize Tracer
	tp, err := tracing.InitTracer(cfg)
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("Error shutting down tracer provider", "error", err)
		}
	}()

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

	app.Use(otelfiber.Middleware())
	app.Use(middleware.RequestResponseLogger())

	app.Static("/doc", "./doc/swagger")

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	// Инициализация репозитория
	lessonRepo := repository.NewLessonRepository(dbPool)
	categoryRepo := repository.NewCategoryRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)

	// Инициализация сервиса
	lessonService := service.NewLessonService(lessonRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo, categoryRepo)

	// Инициализация хэндлеров
	lessonHandler := handler.NewLessonHandler(lessonService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	courseHandler := handler.NewCourseHandler(courseService, categoryService)

	// Установка маршрутов
	categoryRouter := apiV1.Group("/categories")
	courseRouter := categoryRouter.Group(apiconst.PathCategory + "/courses")

	// Регистрация маршрутов
	categoryHandler.RegisterRoutes(categoryRouter)
	courseHandler.RegisterRoutes(categoryRouter)
	lessonHandler.RegisterRoutes(courseRouter)

	slog.Info("Starting server", "address", cfg.Port)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
