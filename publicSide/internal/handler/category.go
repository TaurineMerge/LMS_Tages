package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("")
	}

	categories, pagination, err := h.service.GetAll(c.Context(), query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data: dto.PaginatedCategoriesData{
			Items:      categories,
			Pagination: pagination,
		},
	})
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID()
	}

	category, err := h.service.GetByID(c.Context(), categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data:   category,
	})
}