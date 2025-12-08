package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
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
		return apperrors.NewInvalidUUID()
	}

	var query dto.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("")
	}

	lessons, pagination, err := h.service.GetAllByCourseID(c.Context(), courseID, query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data: dto.PaginatedLessonsData{
			Items:      lessons,
			Pagination: pagination,
		},
	})
}

func (h *LessonHandler) GetLessonByID(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID()
	}

	lessonID := c.Params("lesson_id")
	if _, err := uuid.Parse(lessonID); err != nil {
		return apperrors.NewInvalidUUID()
	}

	lesson, err := h.service.GetByID(c.Context(), lessonID)
	if err != nil {
		return err
	}

	// Ensure the lesson belongs to the course
	if lesson.CourseID != courseID {
		return apperrors.NewNotFound("Lesson in this course")
	}

	return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
		Status: dto.StatusSuccess,
		Data:   lesson,
	})
}