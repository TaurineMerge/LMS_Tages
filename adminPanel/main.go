// main.go
package main

import (
	"embed"
	"io/fs"
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

	log.Printf("⚠️  Failed to initialize auth: %v", settings)
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
	app.Use(middleware.TrustProxyMiddleware())

	// Общий обработчик ошибок
	app.Use(middleware.ErrorHandlerMiddleware())

	// Публичные маршруты (без префикса /admin)
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	// Админские маршруты с префиксом /admin
	adminGroup := app.Group("/admin")

	adminGroup.Use(middleware.AuthMiddleware())

	// Swagger под префиксом /admin
	// Сначала маршрут для doc.json (должен быть до swagger UI)
	adminGroup.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		data, err := fs.ReadFile(swaggerJSON, "docs/swagger.json")
		if err != nil {
			log.Printf("Failed to read swagger.json: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read API documentation",
			})
		}

		c.Set("Content-Type", "application/json")
		return c.Send(data)
	})

	// Затем Swagger UI
	adminGroup.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/admin/swagger/doc.json", // Путь должен быть полный
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
	api := adminGroup.Group("/api/v1")

	// Аутентификация для API маршрутов
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

	// Favicon заглушка
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No Content
	})

	// Запуск сервера
	log.Printf("🚀 Server starting on %s", settings.APIAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/admin/swagger/", settings.APIAddress)
	log.Printf("📖 Swagger JSON: http://localhost%s/admin/swagger/doc.json", settings.APIAddress)
	log.Printf("🏥 Health check: http://localhost%s/health", settings.APIAddress)
	// Инициализация аутентификации
	if err := middleware.InitAuth(); err != nil {
		log.Printf("⚠️  Failed to initialize auth123: %v", err)
		return
	}
	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
