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
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/router"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/database"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/template"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/tracing"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment
	if err := godotenv.Load(); err != nil {
		slog.Warn("Not loaded .env file, using environment variables", "error", err)
	}

	// 2. Initialize Configuration
	cfg, err := config.New(
		config.WithDBFromEnv(),
		config.WithCORSFromEnv(),
		config.WithPortFromEnv(),
		config.WithTracingFromEnv(),
		config.WithLogLevelFromEnv(),
		config.WithDevFromEnv(),
	)
	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	// 3. Initialize Logger
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

	slog.Info("Application starting", "DEV_MODE", cfg.Dev)

	// 4. Initialize Tracer
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

	// 5. Initialize Database
	dbPool, err := database.NewConnection(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	slog.Info("Database connection pool established")

	// 6. Initialize Services and Handlers
	lessonRepo := repository.NewLessonRepository(dbPool)
	categoryRepo := repository.NewCategoryRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)

	lessonService := service.NewLessonService(lessonRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo, categoryRepo)

	// 7. Initialize Fiber App and Global Middleware
	engine := template.NewEngine(cfg)
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
		Views:        engine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))
	app.Use(otelfiber.Middleware())
	app.Use(middleware.RequestResponseLogger())

	// 8. Setup Routes
	webRouter := &router.WebRouter{
		Config:              cfg,
		HomeHandler:         web.NewHomeHandler(),
		CategoryPageHandler: web.NewCategoryHandler(categoryService, courseService),
		CoursesHandler:      web.NewCoursesHandler(courseService, lessonService),
		WebLessonHandler:    web.NewLessonHandler(lessonService, courseService),
	}
	webRouter.Setup(app)

	apiRouter := &router.APIRouter{
		APICategoryHandler: v1.NewCategoryHandler(categoryService),
		APICourseHandler:   v1.NewCourseHandler(courseService),
		APILessonHandler:   v1.NewLessonHandler(lessonService),
	}
	apiRouter.Setup(app)

	// 9. Start Server
	slog.Info("Starting server", "address", cfg.Port)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
