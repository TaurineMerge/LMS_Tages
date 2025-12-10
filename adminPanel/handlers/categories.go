package handlers

import (
	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CategoryHandler - обработчик для категорий
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler создает обработчик категорий
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	categories.Get("/", h.getCategories)                 // GET /api/v1/categories
	categories.Post("/", h.createCategory)               // POST /api/v1/categories
	categories.Get("/:category_id", h.getCategory)       // GET /api/v1/categories/{id}
	categories.Put("/:category_id", h.updateCategory)    // PUT /api/v1/categories/{id}
	categories.Delete("/:category_id", h.deleteCategory) // DELETE /api/v1/categories/{id}
}

func (h *CategoryHandler) getCategories(c *fiber.Ctx) error {
	// Получаем параметры пагинации
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	categories, total, err := h.categoryService.GetCategories(c.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(h.formatError(appErr))
		}
		return c.Status(500).JSON(h.formatServerError(err.Error()))
	}

	// Рассчитываем общее количество страниц
	pages := (total + limit - 1) / limit

	response := models.PaginatedCategoriesResponse{
		Status: "success",
		Data: struct {
			Items      []models.Category `json:"items"`
			Pagination models.Pagination `json:"pagination"`
		}{
			Items: categories,
			Pagination: models.Pagination{
				Total: total,
				Page:  page,
				Limit: limit,
				Pages: pages,
			},
		},
	}

	return c.JSON(response)
}

func (h *CategoryHandler) getCategory(c *fiber.Ctx) error {
	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(h.formatInvalidUUID("Invalid category ID format"))
	}

	category, err := h.categoryService.GetCategory(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(h.formatError(appErr))
		}
		return c.Status(500).JSON(h.formatServerError(err.Error()))
	}

	return c.JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

func (h *CategoryHandler) createCategory(c *fiber.Ctx) error {
	var input models.CategoryCreate

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(h.formatValidationError("Invalid request body"))
	}

	// Валидация
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(h.formatServerError("Validation error: " + err.Error()))
	} else if len(validationErrors) > 0 {
		return c.Status(400).JSON(h.formatValidationErrorWithDetails("Validation failed", validationErrors))
	}

	// Создаем категорию
	category, err := h.categoryService.CreateCategory(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			if appErr.StatusCode == 409 {
				return c.Status(409).JSON(h.formatConflict(appErr.Message))
			}
			return c.Status(appErr.StatusCode).JSON(h.formatError(appErr))
		}
		return c.Status(500).JSON(h.formatServerError(err.Error()))
	}

	return c.Status(201).JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

func (h *CategoryHandler) updateCategory(c *fiber.Ctx) error {
	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(h.formatInvalidUUID("Invalid category ID format"))
	}

	var input models.CategoryUpdate

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(h.formatValidationError("Invalid request body"))
	}

	// Валидация
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(h.formatServerError("Validation error: " + err.Error()))
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(h.formatValidationErrorWithDetails("Validation failed", validationErrors))
	}

	category, err := h.categoryService.UpdateCategory(c.Context(), id, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			if appErr.StatusCode == 409 {
				return c.Status(409).JSON(h.formatConflict(appErr.Message))
			}
			return c.Status(appErr.StatusCode).JSON(h.formatError(appErr))
		}
		return c.Status(500).JSON(h.formatServerError(err.Error()))
	}

	return c.JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

func (h *CategoryHandler) deleteCategory(c *fiber.Ctx) error {
	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(h.formatInvalidUUID("Invalid category ID format"))
	}

	err := h.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			if appErr.StatusCode == 409 {
				return c.Status(409).JSON(h.formatConflict("Category has related courses"))
			}
			return c.Status(appErr.StatusCode).JSON(h.formatError(appErr))
		}
		return c.Status(500).JSON(h.formatServerError(err.Error()))
	}

	return c.JSON(models.StatusOnly{
		Status: "success",
	})
}

// Вспомогательные методы для форматирования ошибок
// Вспомогательные методы для форматирования ошибок
func (h *CategoryHandler) formatError(appErr *exceptions.AppError) models.ErrorResponse {
	return models.ErrorResponse{
		Error: appErr.Message, // Просто строка
		Code:  appErr.Code,    // Просто строка
	}
}

func (h *CategoryHandler) formatServerError(message string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: message,
		Code:  "INTERNAL_SERVER_ERROR",
	}
}

func (h *CategoryHandler) formatInvalidUUID(message string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: message,
		Code:  "INVALID_UUID",
	}
}

func (h *CategoryHandler) formatNotFound(resource string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: resource + " not found",
		Code:  "NOT_FOUND",
	}
}

func (h *CategoryHandler) formatConflict(message string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: message,
		Code:  "CONFLICT",
	}
}

func (h *CategoryHandler) formatValidationError(message string) models.ErrorResponse {
	return models.ErrorResponse{
		Error: message,
		Code:  "VALIDATION_ERROR",
	}
}

func (h *CategoryHandler) formatValidationErrorWithDetails(message string, details map[string]string) models.ValidationErrorResponse {
	return models.ValidationErrorResponse{
		Error:  message,
		Code:   "VALIDATION_ERROR",
		Errors: details,
	}
}
