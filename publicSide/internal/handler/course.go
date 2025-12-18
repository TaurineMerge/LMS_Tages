// Package handler contains the HTTP handlers for the application.
// It is responsible for parsing requests, calling services, and formatting responses.
package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CourseHandler handles HTTP API requests related to courses.
type CourseHandler struct {
	service service.CourseService
}

// NewCourseHandler creates a new instance of a course handler.
func NewCourseHandler(s service.CourseService) *CourseHandler {
	return &CourseHandler{service: s}
}

// RegisterRoutes registers the routes for course-related API endpoints.
func (h *CourseHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	router.Get("/", h.GetCoursesByCategoryID)
	router.Get(apiconst.PathCourse, h.GetCourseByID)
	return router
}

// GetCoursesByCategoryID handles the API request to get a paginated and filtered list of courses for a category.
func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	// Parse pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	courses, pagination, err := h.service.GetCoursesByCategoryID(c.UserContext(), categoryID, page, limit)
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
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.ParamCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCategoryID)
	}

	courseID := c.Params(apiconst.ParamCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.ParamCourseID)
	}

	course, err := h.service.GetCourseByID(c.UserContext(), courseID)
	if err != nil {
		return err
	}

	// Verify the course belongs to the specified category
	if course.CategoryID != categoryID {
		return apperrors.NewNotFoundError("Course not found in specified category")
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   course,
	})
}
