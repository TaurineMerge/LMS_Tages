package handlers

import (
	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CategoryHandler - HTTP обработчик для операций с категориями
//
// Обработчик предоставляет REST API для управления категориями курсов:
//   - GET /categories - получение списка категорий
//   - POST /categories - создание новой категории
//   - GET /categories/:id - получение категории по ID
//   - PUT /categories/:id - обновление категории
//   - DELETE /categories/:id - удаление категории
//
// Особенности:
//   - Валидация входных данных
//   - Интеграция с OpenTelemetry для трассировки
//   - Стандартизированный формат ответов
//   - Централизованная обработка ошибок
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler создает новый HTTP обработчик для категорий
//
// Параметры:
//   - categoryService: сервис для работы с категориями
//
// Возвращает:
//   - *CategoryHandler: указатель на новый обработчик
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// RegisterRoutes регистрирует маршруты для категорий
//
// Регистрирует все необходимые маршруты в указанном роутере.
//
// Параметры:
//   - router: роутер Fiber для регистрации маршрутов
func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	categories.Get("/", h.getCategories)
	categories.Post("/", h.createCategory)
	categories.Get("/:category_id", h.getCategory)
	categories.Put("/:category_id", h.updateCategory)
	categories.Delete("/:category_id", h.deleteCategory)
}

// getCategories обрабатывает GET /categories
//
// Возвращает список всех категорий с пагинацией.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (h *CategoryHandler) getCategories(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода с контекстом трассировки
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCategories.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("http.query", c.Context().QueryArgs().String()),
		))

	categories, err := h.categoryService.GetCategories(ctx)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	response := models.PaginatedCategoriesResponse{
		Status: "success",
	}
	response.Data.Items = categories
	response.Data.Pagination = models.Pagination{
		Total: len(categories),
		Page:  1,
		Limit: len(categories),
		Pages: 1,
	}

	// Логируем успешное завершение
	span.AddEvent("handler.getCategories.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(categories)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(response)
}

// getCategory обрабатывает GET /categories/:id
//
// Возвращает категорию по уникальному идентификатору.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (h *CategoryHandler) getCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	category, err := h.categoryService.GetCategory(ctx, id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.getCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// createCategory обрабатывает POST /categories
//
// Создает новую категорию с валидацией данных.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (h *CategoryHandler) createCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.createCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
		))

	// Валидация входных данных
	var input models.CategoryCreate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.createCategory.request_body",
			trace.WithAttributes(
				attribute.String("request.body", string(body)),
			))
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request body",
			},
		})
	}

	// Валидация через функцию
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Validation error",
			},
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Validation failed",
			},
			Errors: validationErrors,
		})
	}

	// Создаем категорию
	category, err := h.categoryService.CreateCategory(ctx, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.createCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.Status(201).JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// updateCategory обрабатывает PUT /categories/:id
//
// Обновляет существующую категорию с валидацией данных.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (h *CategoryHandler) updateCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.updateCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	// Валидация входных данных
	var input models.CategoryUpdate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.updateCategory.request_body",
			trace.WithAttributes(
				attribute.String("request.body", string(body)),
			))
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request body",
			},
		})
	}

	// Валидация через функцию
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Validation error",
			},
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Validation failed",
			},
			Errors: validationErrors,
		})
	}

	category, err := h.categoryService.UpdateCategory(ctx, id, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.updateCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(models.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// deleteCategory обрабатывает DELETE /categories/:id
//
// Удаляет категорию по уникальному идентификатору.
// Перед удалением проверяет наличие связанных курсов.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
func (h *CategoryHandler) deleteCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.deleteCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	err := h.categoryService.DeleteCategory(ctx, id)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.deleteCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", id),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
}
