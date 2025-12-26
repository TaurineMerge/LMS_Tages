

// publicSide - это публичная часть образовательной платформы LMS Tages.
// Данный файл является точкой входа в приложение.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/clients/testing"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	v1 "github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/router"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/database"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/logger"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/template"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/tracing"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	// --- Загрузка переменных окружения ---
	if err := godotenv.Load(); err != nil {
		slog.Warn("Not loaded .env file, using environment variables", "error", err)
	}

	// --- Конфигурация ---
	cfg, err := config.New(
		config.WithDBFromEnv(),
		config.WithCORSFromEnv(),
		config.WithPortFromEnv(),
		config.WithTracingFromEnv(),
		config.WithLogLevelFromEnv(),
		config.WithDevFromEnv(),
		config.WithOIDCFromEnv(),
		config.WithMinioFromEnv(),
		config.WithTestingFromEnv(),
	)
	if err != nil {
		slog.Error("Failed to initialize config", "error", err)
		os.Exit(1)
	}

	// --- Логирование и Трассировка ---
	slog.SetDefault(logger.Setup(cfg.Log.Level))
	slog.Info("Application starting", "DEV_MODE", cfg.App.Dev)

	tracer, err := tracing.New(&cfg.Otel)
	if err != nil {
		slog.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer tracer.Close()

	// --- Аутентификация (OIDC) ---
	provider, err := oidc.NewProvider(context.Background(), cfg.OIDC.IssuerURL)
	if err != nil {
		slog.Error("Failed to initialize OIDC provider", "error", err)
		os.Exit(1)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.OIDC.ClientID,
		ClientSecret: cfg.OIDC.ClientSecret,
		RedirectURL:  cfg.OIDC.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	authHandler := web.NewAuthHandler(provider, oauth2Config)
	authMiddleware := web.NewAuthMiddleware(provider, cfg.OIDC.ClientID)

	// --- Инициализация зависимостей (DI) ---
	dbPool, err := database.NewConnection(&cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	slog.Info("Database connection pool established")

	s3Service, err := service.NewS3Service(cfg.Minio)
	if err != nil {
		slog.Error("Failed to initialize S3 service", "error", err)
		os.Exit(1)
	}
	slog.Info("S3 service initialized")

	testingClient, err := testing.NewClient(cfg.TestingService.BaseURL, "./doc/schemas/external/testing/get_test_response.json")
	if err != nil {
		slog.Error("Failed to initialize testing client", "error", err)
		os.Exit(1)
	}
	slog.Info("Testing client initialized")

	// Репозитории
	lessonRepo := repository.NewLessonRepository(dbPool)
	categoryRepo := repository.NewCategoryRepository(dbPool)
	courseRepo := repository.NewCourseRepository(dbPool)

	// Сервисы
	lessonService := service.NewLessonService(lessonRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo, categoryRepo, s3Service)
	testService := service.NewTestService(testingClient)
	slog.Info("All services initialized")

	// --- Настройка Fiber ---
	engine := template.NewEngine(&cfg.App)
	app := fiber.New(fiber.Config{
		Views:        engine,
		ErrorHandler: middleware.CommonErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))
	app.Use(otelfiber.Middleware())
	app.Use(middleware.RequestResponseLogger())

	// --- Роутинг ---
	webRouter := &router.WebRouter{
		Config:              &cfg.App,
		HomeHandler:         web.NewHomeHandler(categoryService, courseService),
		CategoryPageHandler: web.NewCategoryHandler(categoryService, courseService),
		CoursesHandler:      web.NewCoursesHandler(courseService, categoryService, lessonService, testService, cfg.TestingService),
		WebLessonHandler:    web.NewLessonHandler(lessonService, courseService, categoryService),
		AuthHandler:         authHandler,
		AuthMiddleware:      authMiddleware,
	}
	webRouter.Setup(app)

	apiRouter := &router.APIRouter{
		APICategoryHandler: v1.NewCategoryHandler(categoryService),
		APICourseHandler:   v1.NewCourseHandler(courseService),
		APILessonHandler:   v1.NewLessonHandler(lessonService),
	}
	apiRouter.Setup(app)

	// --- Запуск сервера ---
	slog.Info("Starting server", "address", cfg.Server.Port)
	if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
