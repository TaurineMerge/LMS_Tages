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

type LessonHandler struct {
	service service.LessonService
}

func NewLessonHandler(s service.LessonService) *LessonHandler {
	return &LessonHandler{service: s}
}

func (h *LessonHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	lessonRouter := router.Group(apiconst.PathCategory + "/lessons")

	lessonRouter.Get("/", h.GetLessonsByCourseID)
	lessonRouter.Get(apiconst.PathLesson, h.GetLessonByID)

	return lessonRouter
}

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

	lessons, pagination, err := h.service.GetAllByCourseID(c.Context(), categoryID, courseID, query.Page, query.Limit)
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

	lesson, err := h.service.GetByID(c.Context(), categoryID, courseID, lessonID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   lesson,
	})
}