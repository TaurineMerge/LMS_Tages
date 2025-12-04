package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

//go:embed docs/swagger.json
var swaggerJSON embed.FS

func main() {
	app := fiber.New()

	// Middleware CORS
	app.Use(cors.New())

	// Swagger JSON endpoint
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		data, err := fs.ReadFile(swaggerJSON, "docs/swagger.json")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to read swagger.json",
			})
		}
		c.Set("Content-Type", "application/json")
		return c.Send(data)
	})

	// Используем готовый Swagger UI из пакета
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json", // Используем относительный путь
		DeepLinking: true,
	}))

	// Course endpoints
	courseGroup := app.Group("/courses")
	{
		courseGroup.Post("/", createCourse)
		courseGroup.Get("/", getAllCourses)
		courseGroup.Get("/:id", getCourseByID)
		courseGroup.Put("/:id", updateCourse)
		courseGroup.Delete("/:id", deleteCourse)
	}

	// Lesson endpoints
	lessonGroup := app.Group("/lessons")
	{
		lessonGroup.Post("/", createLesson)
		lessonGroup.Get("/", getAllLessons)
		lessonGroup.Get("/:id", getLessonByID)
		lessonGroup.Put("/:id", updateLesson)
		lessonGroup.Delete("/:id", deleteLesson)
		lessonGroup.Get("/course/:course_id", getLessonsByCourseID)
	}

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "API работает нормально",
			"port":    4000,
		})
	})

	// Redirect root to Swagger UI
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/")
	})

	log.Println("🚀 Сервер запущен на http://localhost:4000")
	log.Println("📚 Swagger UI доступен по адресу http://localhost:4000/swagger/")
	log.Println("📄 Swagger JSON доступен по адресу http://localhost:4000/swagger/doc.json")
	log.Fatal(app.Listen(":4000"))
}

// Course handlers
func createCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Курс успешно создан",
	})
}

func getAllCourses(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Список курсов успешно получен",
	})
}

func getCourseByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Курс успешно получен",
		"id":      c.Params("id"),
	})
}

func updateCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Курс успешно обновлен",
		"id":      c.Params("id"),
	})
}

func deleteCourse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Курс успешно удален",
		"id":      c.Params("id"),
	})
}

// Lesson handlers
func createLesson(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Урок успешно создан",
	})
}

func getAllLessons(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Список уроков успешно получен",
	})
}

func getLessonByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Урок успешно получен",
		"id":      c.Params("id"),
	})
}

func updateLesson(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Урок успешно обновлен",
		"id":      c.Params("id"),
	})
}

func deleteLesson(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Урок успешно удален",
		"id":      c.Params("id"),
	})
}

func getLessonsByCourseID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":   "success",
		"message":  "Уроки для курса успешно получены",
		"courseId": c.Params("course_id"),
	})
}
