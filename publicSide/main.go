package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
)

func main() {

	// 1. Отладочная информация
	wd, _ := os.Getwd()
	fmt.Println("📁 Текущая директория:", wd)

	// 2. Проверяем файлы перед инициализацией
	checkFile := func(path string) bool {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("✅ Найден: %s\n", path)
			return true
		} else {
			fmt.Printf("❌ Не найден: %s (%v)\n", path, err)
			return false
		}
	}

	layoutsOk := checkFile("layouts/main.hbs")
	pagesIndexOk := checkFile("pages/index.hbs")
	pagesCoursesOk := checkFile("pages/courses.hbs")
	_ = checkFile("pages/course.hbs")

	if !layoutsOk || !pagesIndexOk || !pagesCoursesOk {
		log.Fatal("❌ Не найдены необходимые файлы шаблонов")
	}

	// Инициализация движка шаблонов
	engine := handlebars.New(".", ".hbs")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Тестовый маршрут - проверяем что работает
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("Тест работает! Fiber работает!")
	})

	// Тестовый маршрут с простым рендерингом
	app.Get("/test-page", func(c *fiber.Ctx) error {
		return c.Render("pages/index", fiber.Map{
			"Title": "Тестовая страница",
		}, "layouts/main")
	})

	// Middleware
	app.Use(func(c *fiber.Ctx) error {
		// Устанавливаем текущий год для всех шаблонов
		c.Locals("currentYear", time.Now().Year())
		return c.Next()
	})

	// Статические файлы
	//app.Static("/static", "./publicSide/static")
	//app.Static("/assets", "./publicSide/assets")

	// Главная страница
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/index", fiber.Map{
			"Title":          "Tages LMS - Главная",
			"AdminURL":       "http://localhost:4000",
			"SwaggerURL":     "http://localhost:4000/api/docs", // Ссылка на Swagger в админке
			"ShowAdminBtn":   true,
			"ShowSwaggerBtn": true,
		}, "layouts/main")
	})

	// Список курсов
	app.Get("/courses", func(c *fiber.Ctx) error {
		courses := []map[string]interface{}{
			{
				"ID":          1,
				"Title":       "Введение в Go",
				"Description": "Основы языка Go для начинающих",
				"Instructor":  "Иван Иванов",
				"Category":    "Программирование",
				"Rating":      4.8,
				"Students":    1250,
				"Lessons":     12,
				"Duration":    "8 часов",
			},
			{
				"ID":          2,
				"Title":       "Веб-разработка",
				"Description": "Создание веб-приложений",
				"Instructor":  "Мария Петрова",
				"Category":    "Веб",
				"Rating":      4.9,
				"Students":    890,
				"Lessons":     18,
				"Duration":    "15 часов",
			},
			{
				"ID":          3,
				"Title":       "Базы данных",
				"Description": "PostgreSQL и MySQL",
				"Instructor":  "Алексей Сидоров",
				"Category":    "Базы данных",
				"Rating":      4.7,
				"Students":    540,
				"Lessons":     10,
				"Duration":    "10 часов",
			},
		}

		return c.Render("pages/courses", fiber.Map{
			"Title":   "Все курсы",
			"Courses": courses,
		}, "layouts/main")
	})

	// Страница курса
	app.Get("/courses/:id", func(c *fiber.Ctx) error {
		courseID := c.Params("id")

		// Моковый курс
		course := map[string]interface{}{
			"ID":          courseID,
			"Title":       "Курс " + courseID,
			"Description": "Описание курса " + courseID,
			"Instructor":  "Преподаватель",
			"Rating":      4.5,
			"Students":    1000,
			"Lessons":     15,
			"Duration":    "12 часов",
		}

		// Моковые уроки
		lessons := []map[string]interface{}{
			{"ID": 1, "Title": "Урок 1: Введение", "Duration": "45 мин", "IsFree": true},
			{"ID": 2, "Title": "Урок 2: Основы", "Duration": "60 мин", "IsFree": true},
			{"ID": 3, "Title": "Урок 3: Практика", "Duration": "90 мин", "IsFree": false},
		}

		return c.Render("pages/course", fiber.Map{
			"Title":   course["Title"],
			"Course":  course,
			"Lessons": lessons,
		}, "layouts/main")
	})

	// Обработчик 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("public/404", fiber.Map{
			"title": "Страница не найдена",
		})
	})

	fmt.Println("🚀 Запускаем сервер на порту 3000...")

	// Убиваем старый процесс если есть
	fmt.Println("🛑 Проверяем порт 3000...")

	// Запускаем сервер на порту 3000
	log.Printf("🌐 Публичный сайт запущен на http://localhost:3000")
	log.Printf("📊 Главная страница: http://localhost:3000")
	log.Printf("📚 Все курсы: http://localhost:3000/courses")
	log.Println("🔗 Админ-панель: http://localhost:4000")
	log.Println("📚 Swagger: http://localhost:4000/admin/swagger/")

	log.Fatal(app.Listen(":3000"))
}
