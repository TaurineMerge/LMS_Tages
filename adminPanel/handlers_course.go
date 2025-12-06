package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ============ COURSES ============

// GET /courses
func getCourses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	level := c.Query("level")
	visibility := c.Query("visibility")
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	if categoryID != "" && !isValidUUID(categoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()
	result, err := GetCoursesService(ctx, page, limit, level, visibility, categoryID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to get courses"))
	}

	return c.JSON(toPaginatedCoursesDTO(result))
}

// POST /courses
func createCourse(c *fiber.Ctx) error {
	var input CourseCreateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title == "" {
		return c.Status(400).JSON(newValidationError("Title is required"))
	}

	if input.CategoryID == "" {
		return c.Status(400).JSON(newValidationError("Category ID is required"))
	}

	if !isValidUUID(input.CategoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(newValidationError("Level must be one of: hard, medium, easy"))
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(newValidationError("Visibility must be one of: draft, public, private"))
	}

	if input.Level == "" {
		input.Level = "medium"
	}
	if input.Visibility == "" {
		input.Visibility = "draft"
	}

	ctx := context.Background()
	course, err := CreateCourseService(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Category not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to create course"))
	}

	return c.Status(201).JSON(toCourseDTO(course))
}

// GET /courses/:course_id
func getCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(newBadRequestError("Invalid course ID format"))
	}

	ctx := context.Background()
	course, err := GetCourseService(ctx, courseID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Course not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to get course"))
	}

	return c.JSON(toCourseDTO(course))
}

// PUT /courses/:course_id
func updateCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(newBadRequestError("Invalid course ID format"))
	}

	var input CourseUpdateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title != "" && len(input.Title) > 255 {
		return c.Status(400).JSON(newValidationError("Title must be less than 255 characters"))
	}

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(newValidationError("Level must be one of: hard, medium, easy"))
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(newValidationError("Visibility must be one of: draft, public, private"))
	}

	if input.CategoryID != "" && !isValidUUID(input.CategoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()
	course, err := UpdateCourseService(ctx, courseID, input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			if input.CategoryID != "" {
				return c.Status(404).JSON(newNotFoundError("Category not found"))
			}
			return c.Status(404).JSON(newNotFoundError("Course not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to update course"))
	}

	return c.JSON(toCourseDTO(course))
}

// DELETE /courses/:course_id
func deleteCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(newBadRequestError("Invalid course ID format"))
	}

	ctx := context.Background()
	rows, err := DeleteCourseService(ctx, courseID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to delete course"))
	}

	if rows == 0 {
		return c.Status(404).JSON(newNotFoundError("Course not found"))
	}

	return c.SendStatus(204)
}
