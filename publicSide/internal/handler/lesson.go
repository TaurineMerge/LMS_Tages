package handler

import (
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(s service.LessonService) *LessonHandler {
	return &LessonHandler{service: s}
}

func (h *LessonHandler) GetLessonsByCourseID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if _, err := uuid.Parse(courseID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid course ID format"},
		})
	}

	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_REQUEST", Message: "Invalid query parameters"},
		})
	}

	lessons, pagination, err := h.service.GetAllByCourseID(c.Context(), courseID, query.Page, query.Limit)
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
		Data: dto.PaginatedLessonsData{
			Items:      lessons,
			Pagination: pagination,
		},
	})
}

func (h *LessonHandler) GetLessonByID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if _, err := uuid.Parse(courseID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid course ID format"},
		})
	}

	lessonID := c.Params("lesson_id")
	if _, err := uuid.Parse(lessonID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "INVALID_ID", Message: "Invalid lesson ID format"},
		})
	}

	// We could also validate that the lesson belongs to the course, but the swagger doesn't require it
	// and our mock repo doesn't support that check easily.

	lesson, err := h.service.GetByID(c.Context(), lessonID)
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

	// Ensure the lesson belongs to the course
	if lesson.CourseID != courseID {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Status: "error",
			Error:  dto.ErrorDetail{Code: "NOT_FOUND", Message: "lesson not found in this course"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: "success",
		Data:   lesson,
	})
}
