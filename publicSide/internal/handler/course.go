package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
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
		return apperrors.NewInvalidRequest("")
	}

	courses, pagination, err := h.service.GetAll(c.Context(), query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data: dto.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}

func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID()
	}

	course, err := h.service.GetByID(c.Context(), courseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data:   course,
	})
}

func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID()
	}

	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("")
	}

	courses, pagination, err := h.service.GetAllByCategoryID(c.Context(), categoryID, query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data: dto.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}