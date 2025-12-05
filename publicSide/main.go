package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// --- Data Structures based on swager.json ---

// ContentBlock represents a block of content in a lesson (text or image).
type ContentBlock struct {
	ContentType string `json:"content_type"` // "text" or "image"
	Data        string `json:"data,omitempty"`
	URL         string `json:"url,omitempty"`
	Alt         string `json:"alt,omitempty"`
}

// Lesson represents a lesson within a course, as per swagger definition.
type Lesson struct {
	UUID    string         `json:"uuid"`
	Title   string         `json:"title"`
	Content []ContentBlock `json:"content"`
}

// Course represents a course in the system, as per swagger definition.
type Course struct {
	UUID        string   `json:"uuid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Level       string   `json:"level"` // easy, medium, hard
	Lessons     []Lesson `json:"lessons"`
}

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Group routes under /api/v1 as per swagger basePath
	api := app.Group("/api/v1")

	// --- Courses Endpoints ---
	coursesAPI := api.Group("/courses")
	coursesAPI.Get("/", getCourses)
	coursesAPI.Get("/:course_id", getCourseByID)
	coursesAPI.Get("/:course_id/lessons", getLessonsByCourseID)
	coursesAPI.Get("/:course_id/lessons/:lesson_id", getLessonByCourseAndLessonID)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}

// --- Handler Functions ---

// getCourses handles GET /api/v1/courses
func getCourses(c fiber.Ctx) error {
	// Dummy data matching the swagger definition
	courses := []Course{
		{
			UUID:        "f1b3e6c1-79e0-4cb2-8e23-c0ea680ce621",
			Title:       "Introduction to Go",
			Description: "A beginner's course on Go programming",
			Level:       "easy",
			Lessons:     []Lesson{}, // Lessons are empty in a list view for brevity
		},
		{
			UUID:        "a2c4e8d2-80f1-5db3-9f34-d1fb791df732",
			Title:       "Advanced Go",
			Description: "An advanced course on Go programming",
			Level:       "hard",
			Lessons:     []Lesson{},
		},
	}
	// TODO: Implement filtering based on query params: level, visibility, category_id
	return c.JSON(courses)
}

// getCourseByID handles GET /api/v1/courses/:course_id
func getCourseByID(c fiber.Ctx) error {
	courseIDStr := c.Params("course_id")
	if _, err := uuid.Parse(courseIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID format. ID must be a valid UUID.",
		})
	}

	// Simulate "not found"
	if courseIDStr == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course with ID '" + courseIDStr + "' not found.",
		})
	}

	// Dummy data for a successful case
	course := Course{
		UUID:        courseIDStr,
		Title:       "Introduction to Go",
		Description: "A beginner's course on Go programming",
		Level:       "easy",
		Lessons: []Lesson{
			{
				UUID:    "d2f3b4a5-c6d7-4e8f-9a0b-1c2d3e4f5a6b",
				Title:   "Lesson 1: Hello World",
				Content: []ContentBlock{{ContentType: "text", Data: "This is the first lesson content."}},
			},
		},
	}
	return c.JSON(course)
}

// getLessonsByCourseID handles GET /api/v1/courses/:course_id/lessons
func getLessonsByCourseID(c fiber.Ctx) error {
	courseIDStr := c.Params("course_id")
	if _, err := uuid.Parse(courseIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID format. ID must be a valid UUID.",
		})
	}

	if courseIDStr == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course with ID '" + courseIDStr + "' not found.",
		})
	}

	// Dummy data with multiple lessons and empty Content
	lessons := []Lesson{
		{
			UUID:    "d2f3b4a5-c6d7-4e8f-9a0b-1c2d3e4f5a6b",
			Title:   "Lesson 1: Introduction to " + courseIDStr,
			Content: []ContentBlock{}, // Content is explicitly empty
		},
		{
			UUID:    "b3c4d5e6-f7g8-4h9i-j0k1-l2m3n4o5p6q7",
			Title:   "Lesson 2: Basic Concepts for " + courseIDStr,
			Content: []ContentBlock{},
		},
		{
			UUID:    "e6c1b3f1-79e0-4cb2-8e23-c0ea680ce621",
			Title:   "Lesson 3: Advanced Topics in " + courseIDStr,
			Content: []ContentBlock{},
		},
	}
	return c.JSON(lessons)
}

// getLessonByCourseAndLessonID handles GET /api/v1/courses/:course_id/lessons/:lesson_id
func getLessonByCourseAndLessonID(c fiber.Ctx) error {
	courseIDStr := c.Params("course_id")
	if _, err := uuid.Parse(courseIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID format. ID must be a valid UUID.",
		})
	}

	lessonIDStr := c.Params("lesson_id")
	if _, err := uuid.Parse(lessonIDStr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid lesson ID format. ID must be a valid UUID.",
		})
	}

	if courseIDStr == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found."})
	}
	if lessonIDStr == "00000000-0000-0000-0000-000000000000" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lesson not found."})
	}

	// Dummy data
	lesson := Lesson{
		UUID:  lessonIDStr,
		Title: "Details for Lesson " + lessonIDStr,
		Content: []ContentBlock{
			{ContentType: "text", Data: "This is the content for the lesson."},
			{ContentType: "text", Data: "This is the content for the lesson."},
			{ContentType: "text", Data: "This is the content for the lesson."},
		},
	}
	return c.JSON(lesson)
}
