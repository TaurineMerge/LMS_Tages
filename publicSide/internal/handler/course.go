package handler

import (
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(s service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid query parameters"},
		})
	}

	courses, pagination, err := h.service.GetAll(c.Context(), query.Page, query.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Internal server error"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: "success",
		Data: dto.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}

func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if _, err := uuid.Parse(courseID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid course ID format"},
		})
	}

	course, err := h.service.GetByID(c.Context(), courseID)
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
		Data:   course,
	})
}

func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	if _, err := uuid.Parse(categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid category ID format"},
		})
	}

	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid query parameters"},
		})
	}

	courses, pagination, err := h.service.GetAllByCategoryID(c.Context(), categoryID, query.Page, query.Limit)
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
		Data: dto.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}
