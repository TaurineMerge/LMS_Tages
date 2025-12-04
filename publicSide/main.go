package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

// --- Data Structures ---

// Course represents a course in the system.
type Course struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`       // easy, medium, hard
	Visibility  string `json:"visibility"`  // draft, public, private
	CategoryID  string `json:"category_id"`
}

// Lesson represents a lesson within a course.
type Lesson struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CourseID string `json:"course_id"`
}

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Group routes under /api/v1
	api := app.Group("/api/v1")

	// --- Courses Endpoints ---
	coursesAPI := api.Group("/courses")
	coursesAPI.Get("/", getCourses)
	coursesAPI.Get("/:id", getCourseByID)

	// --- Lessons Endpoints ---
	lessonsAPI := api.Group("/lessons")
	lessonsAPI.Get("/", getLessons)
	lessonsAPI.Get("/:id", getLessonByID)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

// --- Handler Functions ---

// getCourses handles GET /api/v1/courses
func getCourses(c fiber.Ctx) error {
	// Dummy data
	courses := []Course{
		{ID: "f1b3e6c1-79e0-4cb2-8e23-c0ea680ce621", Title: "Introduction to Go", Description: "A beginner's course on Go programming", Level: "easy", Visibility: "public", CategoryID: "c1a9e7f8-6b4e-4f1a-8c1d-0d9e8c7b6a5d"},
		{ID: "a2c4e8d2-80f1-5db3-9f34-d1fb791df732", Title: "Advanced Go", Description: "An advanced course on Go programming", Level: "hard", Visibility: "private", CategoryID: "c1a9e7f8-6b4e-4f1a-8c1d-0d9e8c7b6a5d"},
	}
	// TODO: Implement filtering based on query params: level, visibility, category_id
	return c.JSON(courses)
}

// getCourseByID handles GET /api/v1/courses/:id
func getCourseByID(c fiber.Ctx) error {
	// Dummy data
	course := Course{
		ID:          c.Params("id"),
		Title:       "Introduction to Go",
		Description: "A beginner's course on Go programming",
		Level:       "easy",
		Visibility:  "public",
		CategoryID:  "c1a9e7f8-6b4e-4f1a-8c1d-0d9e8c7b6a5d",
	}
	return c.JSON(course)
}

// getLessons handles GET /api/v1/lessons
func getLessons(c fiber.Ctx) error {
	// Dummy data
	lessons := []Lesson{
		{ID: "d2f3b4a5-c6d7-4e8f-9a0b-1c2d3e4f5a6b", Title: "Lesson 1: Hello World", Content: "This is the first lesson.", CourseID: "f1b3e6c1-79e0-4cb2-8e23-c0ea680ce621"},
		{ID: "b3c4d5e6-f7g8-4h9i-j0k1-l2m3n4o5p6q7", Title: "Lesson 2: Variables", Content: "This is the second lesson.", CourseID: "f1b3e6c1-79e0-4cb2-8e23-c0ea680ce621"},
	}
	// TODO: Implement filtering based on query param: course_id
	return c.JSON(lessons)
}

// getLessonByID handles GET /api/v1/lessons/:id
func getLessonByID(c fiber.Ctx) error {
	// Dummy data
	lesson := Lesson{
		ID:       c.Params("id"),
		Title:    "Lesson 1: Hello World",
		Content:  "This is the first lesson.",
		CourseID: "f1b3e6c1-79e0-4cb2-8e23-c0ea680ce621",
	}
	return c.JSON(lesson)
}