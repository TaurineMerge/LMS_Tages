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

// CourseAPIHandler handles HTTP API requests related to courses.
type CourseAPIHandler struct {
	service service.CourseService
}

// NewCourseAPIHandler creates a new instance of a course API handler.
func NewCourseAPIHandler(s service.CourseService) *CourseAPIHandler {
	return &CourseAPIHandler{service: s}
}

// RegisterRoutes registers the routes for course-related API endpoints.
func (h *CourseAPIHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	router.Get("/", h.GetCoursesByCategoryID)
	router.Get(apiconst.PathCourse, h.GetCourseByID)
	return router
}

// GetCoursesByCategoryID handles the API request to get a paginated and filtered list of courses for a category.
func (h *CourseAPIHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	var query request.CourseQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	// Set defaults if not provided
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.Level == "" {
		query.Level = "all"
	}
	if query.SortBy == "" {
		query.SortBy = "updated_desc"
	}

	courses, pagination, err := h.service.GetCoursesByCategoryID(c.UserContext(), categoryID, query.Page, query.Limit, query.Level, query.SortBy)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCoursesData{
			Items:      courses,
			Pagination: *pagination,
		},
	})
}

// GetCourseByID handles the API request to get a single course by its ID.
func (h *CourseAPIHandler) GetCourseByID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	courseID := c.Params(apiconst.ParamCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCourseID)
	}

	// For now, we'll return a not implemented error since we don't have a GetByID method in the service
	// This can be implemented later if needed
	return apperrors.NewNotFoundError("Course detail endpoint not yet implemented")
}
