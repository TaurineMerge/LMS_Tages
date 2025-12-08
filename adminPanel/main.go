// main.go
package main

import (
	"embed"
	"log"
	"strings"

	"adminPanel/config"
	"adminPanel/database"
	"adminPanel/handlers"
	"adminPanel/middleware"
	"adminPanel/repositories"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

//go:embed docs/swagger.json
var swaggerJSON embed.FS

func main() {
	// Инициализация конфигурации
	settings := config.NewSettings()

	// Инициализация аутентификации
	if err := middleware.InitAuth(); err != nil {
		log.Printf("⚠️  Failed to initialize auth123: %v", err)
		return
	}

	// Инициализация базы данных
	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		AppName:               "Admin Panel API",
		DisableStartupMessage: false,
	})

	// Глобальные middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(settings.GetCORSOrigins(), ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: settings.CORSAllowCredentials,
		ExposeHeaders:    "Content-Length",
	}))

	// Общий обработчик ошибок
	app.Use(middleware.ErrorHandlerMiddleware())

	// Публичные маршруты (без префикса /admin)
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	app.Static("/doc", "./docs")
	
	// Затем Swagger UI
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
		DeepLinking: true,
		Title:       "Admin Panel API",
		OAuth: &swagger.OAuthConfig{
			ClientId:     settings.ClientId,
			ClientSecret: settings.ClientSecret,
			AppName:      settings.AppName,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}))

	// API маршруты с префиксом /admin/api/v1
	api := app.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	// Инициализация репозиториев
	categoryRepo := repositories.NewCategoryRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	// Инициализация сервисов
	categoryService := services.NewCategoryService(categoryRepo)
	courseService := services.NewCourseService(courseRepo, categoryRepo)
	lessonService := services.NewLessonService(lessonRepo, courseRepo)

	// Инициализация обработчиков
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	courseHandler := handlers.NewCourseHandler(courseService)
	lessonHandler := handlers.NewLessonHandler(lessonService)

	// Регистрация маршрутов
	categoryHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	lessonHandler.RegisterRoutes(api)

	// Запуск сервера
	log.Printf("🚀 Server starting on %s", settings.APIAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/", settings.APIAddress)
	log.Printf("📖 Swagger JSON: http://localhost%s/swagger/doc.json", settings.APIAddress)
	log.Printf("🏥 Health check: http://localhost%s/health", settings.APIAddress)
	
	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
