package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Course struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty"`
	CategoryID  string `json:"category_id,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

type Lesson struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	CourseID string                 `json:"course_id"`
	Content  map[string]interface{} `json:"content,omitempty"`
}

type CreateCourseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id" validate:"required,uuid"`
	Visibility  string `json:"visibility,omitempty" validate:"omitempty,oneof=draft public private"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Level       string `json:"level,omitempty" validate:"omitempty,oneof=hard medium easy"`
	CategoryID  string `json:"category_id,omitempty" validate:"omitempty,uuid"`
	Visibility  string `json:"visibility,omitempty" validate:"omitempty,oneof=draft public private"`
}

type CreateLessonRequest struct {
	Title    string                 `json:"title" validate:"required"`
	CourseID string                 `json:"course_id" validate:"required,uuid"`
	Content  map[string]interface{} `json:"content,omitempty"`
}

type UpdateLessonRequest struct {
	Title   string                 `json:"title,omitempty"`
	Content map[string]interface{} `json:"content,omitempty"`
}

// Хранилище в памяти (для примера)
var courses = make(map[string]Course)
var lessons = make(map[string]Lesson)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Education Platform API",
	})

	// Группа для курсов
	coursesGroup := app.Group("/courses")
	{
		// GET /courses - получить все курсы
		coursesGroup.Get("/", getCourses)

		// GET /courses/:id - получить курс по ID
		coursesGroup.Get("/:id", getCourseByID)

		// POST /courses - создать новый курс
		coursesGroup.Post("/", createCourse)

		// PUT /courses/:id - обновить курс
		coursesGroup.Put("/:id", updateCourse)

		// DELETE /courses/:id - удалить курс
		coursesGroup.Delete("/:id", deleteCourse)
	}

	// Группа для уроков
	lessonsGroup := app.Group("/lessons")
	{
		// GET /lessons - получить все уроки
		lessonsGroup.Get("/", getLessons)

		// GET /lessons/:id - получить урок по ID
		lessonsGroup.Get("/:id", getLessonByID)

		// POST /lessons - создать новый урок
		lessonsGroup.Post("/", createLesson)

		// PUT /lessons/:id - обновить урок
		lessonsGroup.Put("/:id", updateLesson)

		// DELETE /lessons/:id - удалить урок
		lessonsGroup.Delete("/:id", deleteLesson)
	}

	// Запуск сервера
	log.Fatal(app.Listen(":3000"))
}

// ==================== Обработчики для курсов ====================

// getCourses - получить все курсы
func getCourses(c fiber.Ctx) error {
	// Простая фильтрация по query параметрам
	level := c.Query("level")
	visibility := c.Query("visibility")
	categoryID := c.Query("category_id")

	// Фильтрация курсов
	filteredCourses := make([]Course, 0)
	for _, course := range courses {
		if level != "" && course.Level != level {
			continue
		}
		if visibility != "" && course.Visibility != visibility {
			continue
		}
		if categoryID != "" && course.CategoryID != categoryID {
			continue
		}
		filteredCourses = append(filteredCourses, course)
	}

	return c.JSON(filteredCourses)
}

// getCourseByID - получить курс по ID
func getCourseByID(c fiber.Ctx) error {
	id := c.Params("id")

	course, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	return c.JSON(course)
}

// createCourse - создать новый курс
func createCourse(c fiber.Ctx) error {
	var req CreateCourseRequest

	// Парсим тело запроса
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация
	if req.Level != "" && req.Level != "hard" && req.Level != "medium" && req.Level != "easy" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Level must be one of: hard, medium, easy",
		})
	}

	if req.Visibility != "" && req.Visibility != "draft" && req.Visibility != "public" && req.Visibility != "private" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visibility must be one of: draft, public, private",
		})
	}

	// Создаем новый курс
	course := Course{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Level:       req.Level,
		CategoryID:  req.CategoryID,
		Visibility:  req.Visibility,
	}

	// Если видимость не указана, ставим draft по умолчанию
	if course.Visibility == "" {
		course.Visibility = "draft"
	}

	// Сохраняем курс
	courses[course.ID] = course

	return c.Status(fiber.StatusCreated).JSON(course)
}

// updateCourse - обновить курс
func updateCourse(c fiber.Ctx) error {
	id := c.Params("id")

	// Проверяем существование курса
	existingCourse, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	var req UpdateCourseRequest

	// Парсим тело запроса
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация
	if req.Level != "" && req.Level != "hard" && req.Level != "medium" && req.Level != "easy" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Level must be one of: hard, medium, easy",
		})
	}

	if req.Visibility != "" && req.Visibility != "draft" && req.Visibility != "public" && req.Visibility != "private" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Visibility must be one of: draft, public, private",
		})
	}

	// Обновляем поля
	if req.Title != "" {
		existingCourse.Title = req.Title
	}
	if req.Description != "" {
		existingCourse.Description = req.Description
	}
	if req.Level != "" {
		existingCourse.Level = req.Level
	}
	if req.CategoryID != "" {
		existingCourse.CategoryID = req.CategoryID
	}
	if req.Visibility != "" {
		existingCourse.Visibility = req.Visibility
	}

	// Сохраняем обновленный курс
	courses[id] = existingCourse

	return c.JSON(existingCourse)
}

// deleteCourse - удалить курс
func deleteCourse(c fiber.Ctx) error {
	id := c.Params("id")

	// Проверяем существование курса
	_, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	// Удаляем курс
	delete(courses, id)

	return c.SendStatus(fiber.StatusNoContent)
}

// ==================== Обработчики для уроков ====================

// getLessons - получить все уроки
func getLessons(c fiber.Ctx) error {
	// Фильтрация по course_id
	courseID := c.Query("course_id")

	filteredLessons := make([]Lesson, 0)
	for _, lesson := range lessons {
		if courseID != "" && lesson.CourseID != courseID {
			continue
		}
		filteredLessons = append(filteredLessons, lesson)
	}

	return c.JSON(filteredLessons)
}

// getLessonByID - получить урок по ID
func getLessonByID(c fiber.Ctx) error {
	id := c.Params("id")

	lesson, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lesson not found",
		})
	}

	return c.JSON(lesson)
}

// createLesson - создать новый урок
func createLesson(c fiber.Ctx) error {
	var req CreateLessonRequest

	// Парсим тело запроса
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Проверяем существование курса
	_, courseExists := courses[req.CourseID]
	if !courseExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	// Создаем новый урок
	lesson := Lesson{
		ID:       uuid.New().String(),
		Title:    req.Title,
		CourseID: req.CourseID,
		Content:  req.Content,
	}

	// Сохраняем урок
	lessons[lesson.ID] = lesson

	return c.Status(fiber.StatusCreated).JSON(lesson)
}

// updateLesson - обновить урок
func updateLesson(c fiber.Ctx) error {
	id := c.Params("id")

	// Проверяем существование урока
	existingLesson, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lesson not found",
		})
	}

	var req UpdateLessonRequest

	// Парсим тело запроса
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Обновляем поля
	if req.Title != "" {
		existingLesson.Title = req.Title
	}
	if req.Content != nil {
		existingLesson.Content = req.Content
	}

	// Сохраняем обновленный урок
	lessons[id] = existingLesson

	return c.JSON(existingLesson)
}

// deleteLesson - удалить урок
func deleteLesson(c fiber.Ctx) error {
	id := c.Params("id")

	// Проверяем существование урока
	_, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lesson not found",
		})
	}

	// Удаляем урок
	delete(lessons, id)

	return c.SendStatus(fiber.StatusNoContent)
}
