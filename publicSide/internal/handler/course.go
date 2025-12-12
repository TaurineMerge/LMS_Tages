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

// CourseHandler handles HTTP requests related to courses.
type CourseHandler struct {
	service service.CourseService
}

// NewCourseHandler creates a new instance of a course handler.
func NewCourseHandler(s service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

// RegisterRoutes registers the routes for course-related endpoints.
func (h *CourseHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	courseRouter := router.Group("/")

	courseRouter.Get("/", h.GetCoursesByCategoryID)
	courseRouter.Get(apiconst.PathCourse, h.GetCourseByID)

	return courseRouter
}

// GetCoursesByCategoryID handles the request to get a paginated list of courses for a category.
func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	var query request.ListQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	courses, pagination, err := h.service.GetAllByCategoryID(c.UserContext(), categoryID, query.Page, query.Limit, query.Sort)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}

// GetCourseByID handles the request to get a single course by its ID.
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	courseID := c.Params(apiconst.ParamCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCourseID)
	}

	course, err := h.service.GetByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   course,
	})
}
