package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func getCourses(c fiber.Ctx) error {
	level := c.Query("level")
	visibility := c.Query("visibility")
	categoryID := c.Query("category_id")

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

func getCourseByID(c fiber.Ctx) error {
	id := c.Params("id")
	course, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}
	return c.JSON(course)
}

func createCourse(c fiber.Ctx) error {
	var req CreateCourseRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Level != "" && req.Level != "hard" && req.Level != "medium" && req.Level != "easy" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Level must be one of: hard, medium, easy"})
	}
	if req.Visibility != "" && req.Visibility != "draft" && req.Visibility != "public" && req.Visibility != "private" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Visibility must be one of: draft, public, private"})
	}

	course := Course{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Level:       req.Level,
		CategoryID:  req.CategoryID,
		Visibility:  req.Visibility,
	}
	if course.Visibility == "" {
		course.Visibility = "draft"
	}

	courses[course.ID] = course
	return c.Status(fiber.StatusCreated).JSON(course)
}

func updateCourse(c fiber.Ctx) error {
	id := c.Params("id")
	existingCourse, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}

	var req UpdateCourseRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Level != "" && req.Level != "hard" && req.Level != "medium" && req.Level != "easy" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Level must be one of: hard, medium, easy"})
	}
	if req.Visibility != "" && req.Visibility != "draft" && req.Visibility != "public" && req.Visibility != "private" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Visibility must be one of: draft, public, private"})
	}

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

	courses[id] = existingCourse
	return c.JSON(existingCourse)
}

func deleteCourse(c fiber.Ctx) error {
	id := c.Params("id")
	_, exists := courses[id]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Course not found"})
	}
	delete(courses, id)
	return c.SendStatus(fiber.StatusNoContent)
}