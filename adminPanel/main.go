package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// ============ MAIN ============

func main() {
	// Инициализация базы данных
	if err := initDB(); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}
	defer dbPool.Close()

	// Инициализация аутентификации (Keycloak JWT). Если переменные окружения
	// не заданы, сервис будет работать без проверки токенов.
	if err := initAuth(); err != nil {
		log.Fatalf("❌ Ошибка инициализации аутентификации: %v", err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Общие middleware
	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Next()
	})

	// Приводим пути вида /admin/api/v1/... к /api/v1/...
	app.Use(stripAdminPrefixMiddleware)

	// Проверка JWT для защищённых эндпоинтов admin panel.
	app.Use(authMiddleware)

	// Swagger UI (использует ваш существующий swagger.json)
	setupSwagger(app)

	// API routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", healthCheck)

	// Categories CRUD
	categories := api.Group("/categories")
	categories.Get("/", getCategories)
	categories.Post("/", createCategory)
	categories.Get("/:category_id", getCategory)
	categories.Put("/:category_id", updateCategory)
	categories.Delete("/:category_id", deleteCategory)
	categories.Get("/:category_id/courses", getCategoryCourses)

	// Courses CRUD
	courses := api.Group("/courses")
	courses.Get("/", getCourses)
	courses.Post("/", createCourse)
	courses.Get("/:course_id", getCourse)
	courses.Put("/:course_id", updateCourse)
	courses.Delete("/:course_id", deleteCourse)
	courses.Get("/:course_id/lessons", getCourseLessons)

	// Lessons CRUD
	lessons := courses.Group("/:course_id/lessons")
	lessons.Get("/", getLessons)
	lessons.Post("/", createLesson)
	lessons.Get("/:lesson_id", getLesson)
	lessons.Put("/:lesson_id", updateLesson)
	lessons.Delete("/:lesson_id", deleteLesson)

	// Start server
	log.Println("🚀 Сервер запущен на :4000")
	log.Println("📚 API доступен по адресу: http://localhost:4000/api")
	log.Println("📄 Swagger UI: http://localhost:4000/swagger/index.html")
	log.Fatal(app.Listen(":4000"))
}


// ============ HANDLERS ============

// (хендлеры вынесены в отдельные файлы handlers_*.go)

