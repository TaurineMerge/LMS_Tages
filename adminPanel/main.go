package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Education Platform API",
	})

	courses := app.Group("/courses")
	courses.Get("/", getCourses)
	courses.Get("/:id", getCourseByID)
	courses.Post("/", createCourse)
	courses.Put("/:id", updateCourse)
	courses.Delete("/:id", deleteCourse)

	lessons := app.Group("/lessons")
	lessons.Get("/", getLessons)
	lessons.Get("/:id", getLessonByID)
	lessons.Post("/", createLesson)
	lessons.Put("/:id", updateLesson)
	lessons.Delete("/:id", deleteLesson)

	log.Fatal(app.Listen("0.0.0.0:4000"))
}