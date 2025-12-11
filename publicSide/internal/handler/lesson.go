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

// LessonHandler handles HTTP requests related to lessons.
type LessonHandler struct {
	service service.LessonService
}

// NewLessonHandler creates a new instance of a lesson handler.
func NewLessonHandler(s service.LessonService) *LessonHandler {
	return &LessonHandler{service: s}
}

// RegisterRoutes registers the routes for lesson-related endpoints.
func (h *LessonHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	lessonRouter := router.Group(apiconst.PathCourse + "/lessons")

	lessonRouter.Get("/", h.GetLessonsByCourseID)
	lessonRouter.Get(apiconst.PathLesson, h.GetLessonByID)

	return lessonRouter
}

// GetLessonsByCourseID handles the request to get a paginated list of lessons for a course.
func (h *LessonHandler) GetLessonsByCourseID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}
	courseID := c.Params(apiconst.ParamCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCourseID)
	}

	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	lessons, pagination, err := h.service.GetAllByCourseID(c.UserContext(), categoryID, courseID, query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedLessonsData{
			Items:      lessons,
			Pagination: pagination,
		},
	})
}

// GetLessonByID handles the request to get a single lesson by its ID.
func (h *LessonHandler) GetLessonByID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}
	courseID := c.Params(apiconst.ParamCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCourseID)
	}

	lessonID := c.Params(apiconst.ParamLessonID)
	if _, err := uuid.Parse(lessonID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamLessonID)
	}

	lesson, err := h.service.GetByID(c.UserContext(), categoryID, courseID, lessonID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   lesson,
	})
}
