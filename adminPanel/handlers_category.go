package main

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ============ CATEGORIES ============

// GET /categories
func getCategories(c *fiber.Ctx) error {
	ctx := context.Background()

	categories, err := GetCategoriesService(ctx)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Internal server error"))
	}

	return c.JSON(toCategoryDTOs(categories))
}

// POST /categories
func createCategory(c *fiber.Ctx) error {
	var input CategoryCreateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title == "" {
		return c.Status(400).JSON(newValidationError("Title is required"))
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(newValidationError("Title must be less than 255 characters"))
	}

	ctx := context.Background()
	category, err := CreateCategoryService(ctx, input.Title)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(409).JSON(newConflictError("Category with this title already exists"))
		}
		return c.Status(500).JSON(newInternalError("Failed to create category"))
	}

	return c.Status(201).JSON(toCategoryDTO(category))
}

// GET /categories/:category_id
func getCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()
	category, err := GetCategoryService(ctx, categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Category not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to get category"))
	}

	return c.JSON(toCategoryDTO(category))
}

// PUT /categories/:category_id
func updateCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()

	// Проверяем, что категория существует
	if _, err := GetCategoryService(ctx, categoryID); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Category not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to check category"))
	}

	var input CategoryUpdateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title == "" {
		return c.Status(400).JSON(newValidationError("Title is required"))
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(newValidationError("Title must be less than 255 characters"))
	}

	category, err := UpdateCategoryService(ctx, input.Title, categoryID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to update category"))
	}

	return c.JSON(toCategoryDTO(category))
}

// DELETE /categories/:category_id
func deleteCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()

	courseCount, err := CountCoursesForCategory(ctx, categoryID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to check associated courses"))
	}

	if courseCount > 0 {
		return c.Status(409).JSON(newConflictError("Cannot delete category with associated courses"))
	}

	rowsAffected, err := DeleteCategoryService(ctx, categoryID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to delete category"))
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(newNotFoundError("Category not found"))
	}

	return c.SendStatus(204)
}

// GET /categories/:category_id/courses
func getCategoryCourses(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(newBadRequestError("Invalid category ID format"))
	}

	ctx := context.Background()
	// Проверяем, что категория существует
	if _, err := GetCategoryService(ctx, categoryID); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Category not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to check category"))
	}

	courses, err := GetCategoryCoursesService(ctx, categoryID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to get courses"))
	}

	return c.JSON(toCourseDTOs(courses))
}
