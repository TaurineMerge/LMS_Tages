package main

import (
	"embed"
	"io/fs"
	"log"
	"strings"
	"time"

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

	// Маршрут для главной страницы админ-панели
	app.Get("/admin", func(c *fiber.Ctx) error {
		html := `
		<!DOCTYPE html>
		<html lang="ru">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Админ-панель LMS TAGES</title>
			<style>
				body { font-family: Arial, sans-serif; padding: 20px; background: #f5f5f5; }
				.container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; }
				h1 { color: #333; }
				.nav { background: #2c3e50; padding: 15px; border-radius: 5px; margin: 20px 0; }
				.nav a { color: white; text-decoration: none; margin-right: 20px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Админ-панель LMS TAGES</h1>
				<p>Система успешно запущена!</p>
				
				<div class="nav">
					<a href="/admin/swagger/">📚 Swagger UI</a>
					<a href="/health">🏥 Health Check</a>
					<a href="http://localhost:3000" target="_blank">🌐 Публичный сайт</a>
				</div>
				
				<div style="margin-top: 30px;">
					<h3>Доступные сервисы:</h3>
					<ul>
						<li>Публичный сайт: <a href="http://localhost:3000" target="_blank">http://localhost:3000</a></li>
						<li>Админ-панель: <a href="/admin">http://localhost:4000/admin</a></li>
						<li>Swagger API: <a href="/admin/swagger/">http://localhost:4000/admin/swagger/</a></li>
						<li>Health Check: <a href="/health">http://localhost:4000/health</a></li>
					</ul>
				</div>
				
				<p style="margin-top: 30px; color: #666;">
					Время сервера: ` + time.Now().Format("15:04:05 02.01.2006") + `
				</p>
			</div>
		</body>
		</html>
		`

		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	})

	// Публичные маршруты (без префикса /admin)
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	// Админские маршруты с префиксом /admin
	adminGroup := app.Group("/admin")

	adminGroup.Use(middleware.AuthMiddleware())

	// Swagger под префиксом /admin
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

	adminGroup.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/admin/swagger/doc.json",
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

	// Запуск сервера
	log.Printf("🚀 Server starting on %s", settings.APIAddress)
	log.Printf("📊 Admin Panel: http://localhost%s/admin", settings.APIAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/admin/swagger/", settings.APIAddress)

	// Инициализация аутентификации
	if err := middleware.InitAuth(); err != nil {
		log.Printf("⚠️  Failed to initialize auth: %v", err)
	}

	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
