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
	
	// Инициализация базы данных
	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()
	
	// Инициализация аутентификации
	if err := middleware.InitAuth(); err != nil {
		log.Printf("⚠️  Failed to initialize auth: %v", err)
	}
	
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
	
	// Swagger под префиксом /admin
	setupSwagger(adminGroup)
	
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
	
	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// setupSwagger настраивает Swagger документацию
func setupSwagger(router fiber.Router) {
	// Сначала маршрут для doc.json (должен быть до swagger UI)
	router.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
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
	router.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/admin/swagger/doc.json",  // Путь должен быть полный
		DeepLinking: true,
		Title:       "Admin Panel API",
	}))
}
// package main

// import (
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"
// 	// "io/ioutil"  // для чтения файлов
// 	// "path/filepath"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/gofiber/fiber/v2/middleware/recover"
// 	"github.com/gofiber/swagger"

// 	"adminPanel/config"
// 	"adminPanel/database"
// 	"adminPanel/handlers"
// 	"adminPanel/middleware"
// 	"adminPanel/repositories"
// 	"adminPanel/services"
// )

// func main() {
// 	// 1. Загружаем конфигурацию
// 	settings := config.NewSettings()
	
// 	log.Printf("🚀 Starting Admin Panel API")
// 	log.Printf("🌐 Listening on: %s", settings.APIAddress)
// 	log.Printf("🔧 CORS Origins: %s", settings.CORSAllowOrigins)
// 	log.Printf("🔧 CORS Credentials: %v", settings.CORSAllowCredentials)
	
// 	// 2. Инициализируем базу данных
// 	db, err := database.InitDB(settings)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to initialize database: %v", err)
// 	}
// 	defer database.Close()

// 	// 3. Инициализируем аутентификацию
// 	if err := middleware.InitAuth(); err != nil {
// 		log.Fatalf("❌ Failed to initialize authentication: %v", err)
// 	}

// 	// 4. Создаем репозитории
// 	categoryRepo := repositories.NewCategoryRepository(db)
// 	courseRepo := repositories.NewCourseRepository(db)
// 	lessonRepo := repositories.NewLessonRepository(db)

// 	// 5. Создаем сервисы
// 	categoryService := services.NewCategoryService(categoryRepo)
// 	courseService := services.NewCourseService(courseRepo, categoryRepo)
// 	lessonService := services.NewLessonService(lessonRepo, courseRepo)

// 	// 6. Создаем обработчики
// 	healthHandler := handlers.NewHealthHandler(db)
// 	categoryHandler := handlers.NewCategoryHandler(categoryService)
// 	courseHandler := handlers.NewCourseHandler(courseService)
// 	lessonHandler := handlers.NewLessonHandler(lessonService)

// 	// 7. Настраиваем приложение Fiber
// 	app := fiber.New(fiber.Config{
// 		AppName:               "Admin Panel API",
// 		DisableStartupMessage: false,
// 		ReadTimeout:           10 * time.Second,
// 		WriteTimeout:          10 * time.Second,
// 		IdleTimeout:           30 * time.Second,
// 	})

// 	// 8. Middleware
// 	app.Use(middleware.TrustProxyMiddleware())
	
// 	if settings.Debug {
// 		app.Use(logger.New(logger.Config{
// 			Format: "[${time}] ${status} - ${method} ${path}\n",
// 		}))
// 	}
	
// 	app.Use(recover.New())
	
// 	// Настраиваем CORS из конфигурации
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins:     settings.CORSAllowOrigins,
// 		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
// 		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
// 		AllowCredentials: settings.CORSAllowCredentials,
// 	}))
	
// 	app.Use(middleware.ErrorHandlerMiddleware())
// 	app.Use(middleware.AuthMiddleware())

// 	// 9. Health check для контейнера
// 	app.Get("/health", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"status":  "healthy",
// 			"service": "admin-panel",
// 			"version": "1.0.0",
// 		})
// 	})

// 	// 10. Swagger UI
// 	setupSwagger(app)

// 	// 11. Регистрируем маршруты
// 	registerRoutes(app, healthHandler, categoryHandler, courseHandler, lessonHandler)

// 	// 12. Запускаем сервер
// 	go func() {
// 		if err := app.Listen(settings.APIAddress); err != nil {
// 			log.Fatalf("❌ Server failed to start: %v", err)
// 		}
// 	}()

// 	log.Printf("✅ Server started successfully")
// 	log.Printf("📚 API available at: http://localhost%s/api/v1", settings.APIAddress)
// 	log.Printf("📄 Swagger UI: http://localhost%s/swagger/index.html", settings.APIAddress)

// 	// 13. Graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	log.Println("Shutting down server...")

// 	if err := app.Shutdown(); err != nil {
// 		log.Fatalf("❌ Server forced to shutdown: %v", err)
// 	}

// 	log.Println("Server exited gracefully")
// }
// func setupSwagger(app *fiber.App) {
//     // Простая спецификация для тестирования
//     app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
//         return c.JSON(fiber.Map{
//             "openapi": "3.0.0",
//             "info": fiber.Map{"title": "API", "version": "1.0.0"},
//             "paths": fiber.Map{
//                 "/health": fiber.Map{
//                     "get": fiber.Map{
//                         "responses": fiber.Map{
//                             "200": fiber.Map{"description": "OK"},
//                         },
//                     },
//                 },
//             },
//         })
//     })
    
//     app.Get("/swagger/*", swagger.New(swagger.Config{
//         URL: "/swagger/doc.json",
//     }))
// }

// func setupSwagger(app *fiber.App) {
//     // Получаем путь к рабочей директории
//     workDir, _ := os.Getwd()
//     swaggerPath := filepath.Join(workDir, "docs", "swagger.json")
    
//     log.Printf("📁 Looking for swagger.json at: %s", swaggerPath)
    
//     // Проверяем, существует ли файл
//     if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
//         log.Printf("⚠️  Swagger file not found at %s, using default", swaggerPath)
//         // Создаем простую спецификацию
//         setupDefaultSwagger(app)
//         return
//     }
    
//     // Swagger UI будет доступен по /swagger/index.html
//     app.Get("/swagger/*", swagger.New(swagger.Config{
//         URL:          "/swagger/doc.json",  // Исправляем путь!
//         DeepLinking:  true,
//         DocExpansion: "list",
//         Title:        "Admin Panel API Documentation",
//     }))
    
//     // Эндпоинт для отдачи swagger.json
//     app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
//         content, err := ioutil.ReadFile(swaggerPath)
//         if err != nil {
//             log.Printf("❌ Failed to read swagger.json: %v", err)
//             return c.Status(500).JSON(fiber.Map{
//                 "error": "Swagger documentation not found",
//             })
//         }
        
//         log.Printf("✅ Swagger.json loaded successfully (%d bytes)", len(content))
//         c.Set("Content-Type", "application/json; charset=utf-8")
//         return c.Send(content)
//     })
    
//     // Редирект для удобства
//     app.Get("/swagger", func(c *fiber.Ctx) error {
//         return c.Redirect("/swagger/index.html")
//     })
// }

// func setupDefaultSwagger(app *fiber.App) {
//     // Простая спецификация по умолчанию
//     app.Get("/swagger/*", swagger.New(swagger.Config{
//         URL:          "/swagger/doc.json",
//         DeepLinking:  true,
//         DocExpansion: "list",
//     }))
    
//     app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
//         return c.JSON(fiber.Map{
//             "openapi": "3.0.0",
//             "info": fiber.Map{
//                 "title":       "Admin Panel API",
//                 "description": "API для админ-панели системы онлайн образования",
//                 "version":     "1.0.0",
//             },
//             "servers": []fiber.Map{
//                 {
//                     "url":         "/api/v1",
//                     "description": "Admin Panel API Server",
//                 },
//             },
//             "paths": fiber.Map{
//                 "/health": fiber.Map{
//                     "get": fiber.Map{
//                         "tags": []string{"Health"},
//                         "summary": "Проверка здоровья сервиса",
//                         "responses": fiber.Map{
//                             "200": fiber.Map{
//                                 "description": "Сервис работает",
//                             },
//                         },
//                     },
//                 },
//             },
//         })
//     })
// }

// registerRoutes регистрирует все маршруты
func registerRoutes(
	app *fiber.App,
	healthHandler *handlers.HealthHandler,
	categoryHandler *handlers.CategoryHandler,
	courseHandler *handlers.CourseHandler,
	lessonHandler *handlers.LessonHandler,
) {
	// Основная группа API
	apiGroup := app.Group("/api/v1")
	
	// Health endpoints
	healthHandler.RegisterRoutes(apiGroup)
	
	// Category endpoints
	categoryHandler.RegisterRoutes(apiGroup)
	
	// Course endpoints
	courseHandler.RegisterRoutes(apiGroup)
	
	// Lesson endpoints
	lessonHandler.RegisterRoutes(apiGroup)
	
	// Дополнительный health check
	apiGroup.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
			"version":  "1.0.0",
		})
	})
}