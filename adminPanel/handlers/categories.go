package handlers

import (
	"adminPanel/handlers/dto/request"
	"adminPanel/handlers/dto/response"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CategoryHandler обрабатывает HTTP-запросы для категорий.
// Содержит сервис для бизнес-логики и методы для маршрутов.
type CategoryHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
// Принимает сервис категорий.
func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// RegisterRoutes регистрирует маршруты для категорий.
// Создает группу /categories и привязывает методы к маршрутам.
func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")

	categories.Get("/", h.getCategories)
	categories.Post("/", middleware.ValidateJSONSchema("category-create.json"), h.createCategory)
	categories.Get("/:category_id", h.getCategory)
	categories.Put("/:category_id", middleware.ValidateJSONSchema("category-update.json"), h.updateCategory)
	categories.Delete("/:category_id", h.deleteCategory)
}

// getCategories обрабатывает GET /categories.
// Возвращает список всех категорий.
func (h *CategoryHandler) getCategories(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCategories.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("http.query", c.Context().QueryArgs().String()),
		))

	categories, err := h.categoryService.GetCategories(ctx)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	resp := response.PaginatedCategoriesResponse{
		Status: "success",
	}
	resp.Data.Items = categories
	resp.Data.Pagination = models.Pagination{
		Total: len(categories),
		Page:  1,
		Limit: len(categories),
		Pages: 1,
	}

	span.AddEvent("handler.getCategories.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(categories)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(resp)
}

// getCategory обрабатывает GET /categories/:category_id.
// Возвращает категорию по ID.
func (h *CategoryHandler) getCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	category, err := h.categoryService.GetCategory(ctx, id)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.getCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(response.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// createCategory обрабатывает POST /categories.
// Создает новую категорию на основе JSON в теле запроса.
func (h *CategoryHandler) createCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.createCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
		))

	var input request.CategoryCreate

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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_JSON",
				Message: "Invalid request body",
			},
		})
	}

	category, err := h.categoryService.CreateCategory(ctx, input)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.createCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.Status(201).JSON(response.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// updateCategory обрабатывает PUT /categories/:category_id.
// Обновляет категорию по ID на основе JSON в теле запроса.
func (h *CategoryHandler) updateCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.updateCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	var input request.CategoryUpdate

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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_JSON",
				Message: "Invalid request body",
			},
		})
	}

	category, err := h.categoryService.UpdateCategory(ctx, id, input)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.updateCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", category.ID),
			attribute.String("category.title", category.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(response.CategoryResponse{
		Status: "success",
		Data:   *category,
	})
}

// deleteCategory обрабатывает DELETE /categories/:category_id.
// Удаляет категорию по ID.
func (h *CategoryHandler) deleteCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.deleteCategory.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	id := c.Params("category_id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	err := h.categoryService.DeleteCategory(ctx, id)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.deleteCategory.end",
		trace.WithAttributes(
			attribute.String("category.id", id),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
}
