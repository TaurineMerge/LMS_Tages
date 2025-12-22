// Package main is the entry point for the publicSide service.
// It initializes the configuration, database, router, and starts the HTTP server.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/router"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/database"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/logger"
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
	slog.SetDefault(logger.Setup(cfg.Log.Level))
	slog.Info("Application starting", "DEV_MODE", cfg.App.Dev)

	// 4. Initialize Tracer
	tracer, err := tracing.New(&cfg.Otel)
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tracer.Close()

	// 5. Initialize Database
	dbPool, err := database.NewConnection(&cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	slog.Info("Database connection pool established")

	// 6. Initialize Services
	lessonRepo := repository.NewLessonRepository(dbPool)
	categoryRepo := repository.NewCategoryRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)

	lessonService := service.NewLessonService(lessonRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo, categoryRepo)

	// 7. Initialize Fiber App and Global Middleware
	engine := template.NewEngine(&cfg.App)
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: middleware.CommonErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))
	app.Use(otelfiber.Middleware())
	app.Use(middleware.RequestResponseLogger())

	// 8. Setup Routes
	webRouter := &router.WebRouter{
		Config:              &cfg.App,
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
	slog.Info("Starting server", "address", cfg.Server.Port)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
