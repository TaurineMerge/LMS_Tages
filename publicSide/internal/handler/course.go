// Package handler contains the HTTP handlers for the application.
package handler

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/request"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CourseHandler handles HTTP requests related to course pages.
type CourseHandler struct {
	courseService   service.CourseService
	categoryService service.CategoryService
}

// NewCourseHandler creates a new instance of a course handler.
func NewCourseHandler(courseService service.CourseService, categoryService service.CategoryService) *CourseHandler {
	return &CourseHandler{
		courseService:   courseService,
		categoryService: categoryService,
	}
}

// RegisterRoutes registers the routes for course page endpoints.
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/:categoryId/courses", h.GetCoursesByCategoryID)
	router.Get("/:categoryId/courses/:courseId", h.GetCourseByID)
}

// GetCoursesByCategoryID handles the request to get paginated courses for a specific category.
func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID("categoryId")
	}

	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	// Get category information to validate it exists
	category, err := h.categoryService.GetByID(c.UserContext(), categoryID)
	if err != nil {
		return err
	}

	// Get courses for this category with pagination
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(c.UserContext(), categoryID, query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCoursesData{
			CategoryID:   category.ID,
			CategoryName: category.Title,
			Items:        courses,
			Pagination:   pagination,
		},
	})
}

// GetCourseByID handles the request to get a single course by its ID.
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID("categoryId")
	}

	courseID := c.Params("courseId")
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID("courseId")
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
