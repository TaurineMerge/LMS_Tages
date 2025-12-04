package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func getLessons(c fiber.Ctx) error {
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

func getLessonByID(c fiber.Ctx) error {
	id := c.Params("id")
	lesson, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lesson not found"})
	}
	return c.JSON(lesson)
}

func createLesson(c fiber.Ctx) error {
	var req CreateLessonRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, courseExists := courses[req.CourseID]
	if !courseExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Course not found"})
	}

	lesson := Lesson{
		ID:       uuid.New().String(),
		Title:    req.Title,
		CourseID: req.CourseID,
		Content:  req.Content,
	}

	lessons[lesson.ID] = lesson
	return c.Status(fiber.StatusCreated).JSON(lesson)
}

func updateLesson(c fiber.Ctx) error {
	id := c.Params("id")
	existingLesson, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lesson not found"})
	}

	var req UpdateLessonRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title != "" {
		existingLesson.Title = req.Title
	}
	if req.Content != nil {
		existingLesson.Content = req.Content
	}

	lessons[id] = existingLesson
	return c.JSON(existingLesson)
}

func deleteLesson(c fiber.Ctx) error {
	id := c.Params("id")
	_, exists := lessons[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lesson not found"})
	}
	delete(lessons, id)
	return c.SendStatus(fiber.StatusNoContent)
}