package handlers

import (
	"github.com/gofiber/fiber/v2"
	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"
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
	
	categories.Get("/", h.getCategories)           // Полный путь: /admin/api/v1/categories
	categories.Post("/", h.createCategory)         // Полный путь: /admin/api/v1/categories
	categories.Get("/:id", h.getCategory)          // Полный путь: /admin/api/v1/categories/:id
	categories.Put("/:id", h.updateCategory)       // Полный путь: /admin/api/v1/categories/:id
	categories.Delete("/:id", h.deleteCategory)    // Полный путь: /admin/api/v1/categories/:id
	categories.Get("/:id/courses", h.getCategoryCourses) // Полный путь: /admin/api/v1/categories/:id/courses
}

// GetCategories godoc
// @Summary Получить список категорий
// @Description Возвращает список всех категорий
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} models.CategoryListResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories [get]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) getCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetCategories(c.Context())
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	response := models.CategoryListResponse{
		Data:  make([]models.CategoryResponse, 0, len(categories)),
		Total: len(categories),
	}

	for _, category := range categories {
		response.Data = append(response.Data, models.CategoryResponse{Category: category})
	}

	return c.JSON(response)
}

// GetCategory godoc
// @Summary Получить категорию по ID
// @Description Возвращает данные категории по указанному ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID категории" format(uuid)
// @Success 200 {object} models.CategoryResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories/{id} [get]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) getCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	category, err := h.categoryService.GetCategory(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(models.CategoryResponse{Category: *category})
}

// CreateCategory godoc
// @Summary Создать новую категорию
// @Description Создает новую категорию в системе
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body models.CategoryCreate true "Данные категории"
// @Success 201 {object} models.CategoryResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories [post]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) createCategory(c *fiber.Ctx) error {
	// Валидация входных данных
	var input models.CategoryCreate
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}
	
	// Валидация через функцию
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Validation error",
			Code:  "INTERNAL_ERROR",
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Error:  "Validation failed",
			Code:   "VALIDATION_ERROR",
			Errors: validationErrors,
		})
	}
	
	// Создаем категорию
	category, err := h.categoryService.CreateCategory(c.Context(), input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.Status(201).JSON(models.CategoryResponse{Category: *category})
}


// UpdateCategory godoc
// @Summary Обновить данные категории
// @Description Обновляет данные существующей категории
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID категории" format(uuid)
// @Param request body models.CategoryUpdate true "Данные для обновления"
// @Success 200 {object} models.CategoryResponse
// @Failure 400 {object} models.ValidationErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories/{id} [put]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) updateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Валидация входных данных
	var input models.CategoryUpdate
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}
	
	// Валидация через функцию
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Validation error",
			Code:  "INTERNAL_ERROR",
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Error:  "Validation failed",
			Code:   "VALIDATION_ERROR",
			Errors: validationErrors,
		})
	}

	category, err := h.categoryService.UpdateCategory(c.Context(), id, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(models.CategoryResponse{Category: *category})
}

// DeleteCategory godoc
// @Summary Удалить категорию
// @Description Удаляет категорию из системы
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID категории" format(uuid)
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories/{id} [delete]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) deleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	err := h.categoryService.DeleteCategory(c.Context(), id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.SendStatus(204)
}

// GetCategoryCourses godoc
// @Summary Получить курсы категории
// @Description Возвращает список курсов для указанной категории
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "ID категории" format(uuid)
// @Success 200 {array} models.CourseResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/api/v1/categories/{id}/courses [get]  // ← ПОЛНЫЙ ПУТЬ!
func (h *CategoryHandler) getCategoryCourses(c *fiber.Ctx) error {
	id := c.Params("id")
	
	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// TODO: Реализовать через CourseService
	return c.Status(501).JSON(models.ErrorResponse{
		Error: "Not implemented yet",
		Code:  "NOT_IMPLEMENTED",
	})
}

// // Вспомогательная функция для валидации UUID
// func isValidUUID(u string) bool {
// 	if len(u) != 36 {
// 		return false
// 	}
// 	// Простая проверка формата
// 	return true
// }