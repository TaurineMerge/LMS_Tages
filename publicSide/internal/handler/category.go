// Package handler contains the HTTP handlers for the application.
// It is responsible for parsing requests, calling services, and formatting responses.
package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/request"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CategoryHandler handles HTTP requests related to categories.
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler creates a new instance of a category handler.
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// RegisterRoutes registers the routes for category-related endpoints.
func (h *CategoryHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	categoriesRouter := router.Group("/categories")
	categoriesRouter.Get("/", h.GetAllCategories)
	
	categoriesIdRouter := categoriesRouter.Group("/:" + apiconst.PathVariableCategoryID)
	categoriesIdRouter.Get("/", h.GetCategoryByID)

	return categoriesIdRouter
}

// GetAllCategories handles the request to get a paginated list of all categories.
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	categories, pagination, err := h.service.GetAll(c.UserContext(), query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCategoriesData{
			Items:      categories,
			Pagination: pagination,
		},
	})
}

// GetCategoryByID handles the request to get a single category by its ID.
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.PathVariableCategoryID)
	}

	category, err := h.service.GetByID(c.UserContext(), categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   category,
	})
}
