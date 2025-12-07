package handler

import (
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
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
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid query parameters"},
		})
	}

	categories, pagination, err := h.service.GetAll(c.Context(), query.Page, query.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Internal server error"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: "success",
		Data: dto.PaginatedCategoriesData{
			Items:      categories,
			Pagination: pagination,
		},
	})
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	if _, err := uuid.Parse(categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid category ID format"},
		})
	}

	category, err := h.service.GetByID(c.Context(), categoryID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Status: "error",
				Error:  dto.ErrorDetail{Code: "NOT_FOUND", Message: err.Error()},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Internal server error"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: "success",
		Data:   category,
	})
}
