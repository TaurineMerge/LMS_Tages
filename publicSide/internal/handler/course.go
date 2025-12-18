// Package handler contains the HTTP handlers for the application.
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

// CourseHandler handles HTTP requests related to course pages.
type CourseHandler struct {
	courseService service.CourseService
}

// NewCourseHandler creates a new instance of a course handler.
func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// RegisterRoutes registers the routes for course page endpoints.
func (h *CourseHandler) RegisterRoutes(router fiber.Router) fiber.Router {
	courseRouter := router.Group("/courses")
	courseRouter.Get("/", h.GetCoursesByCategoryID)

	courseIdRouter := courseRouter.Group("/:" + apiconst.PathVariableCourseID)
	courseIdRouter.Get("/", h.GetCourseByID)

	return courseIdRouter
}

// GetCoursesByCategoryID handles the request to get paginated courses for a specific category.
func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.PathVariableCategoryID)
	}

	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	// Get courses for this category with pagination
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(c.UserContext(), categoryID, query.Page, query.Limit)
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
	categoryID := c.Params(apiconst.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.PathVariableCategoryID)
	}

	courseID := c.Params(apiconst.PathVariableCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(apiconst.PathVariableCourseID)
	}

	course, err := h.courseService.GetCourseByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   course,
	})
}
