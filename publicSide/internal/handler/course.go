// Package handler contains the HTTP handlers for the application.
package handler

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
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
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/categories/:categoryId/courses", h.ShowCourses)
}

// ShowCourses renders the courses page for a specific category.
func (h *CourseHandler) ShowCourses(c *fiber.Ctx) error {
	categoryID := c.Params("categoryId")
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID("categoryId")
	}

	// Parse pagination parameters from query
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 28)

	// Parse filter and sort parameters
	level := c.Query("level", "all")             // all, easy, medium, hard
	sortBy := c.Query("sort_by", "updated_desc") // updated_desc, updated_asc, created_desc, created_asc

	// Get category information
	category, err := h.courseService.GetCategoryByID(c.UserContext(), categoryID)
	if err != nil {
		slog.Error("Failed to get category", "categoryId", categoryID, "error", err)
		return err
	}
	slog.Info("Category found", "categoryId", categoryID, "title", category.Title)

	// Get courses for this category with pagination and filters
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(c.UserContext(), categoryID, page, limit, level, sortBy)

	if err != nil {
		slog.Error("Failed to get courses", "categoryId", categoryID, "error", err)
		return err
	}
	slog.Info("Courses retrieved", "categoryId", categoryID, "count", len(courses), "page", page, "total", pagination.Total)

	// Transform courses to add localized level names
	type CourseView struct {
		ID          string
		Title       string
		Description string
		Level       string
		LevelRu     string
	}

	var coursesView []CourseView
	for _, course := range courses {
		levelRu := "Средний"
		switch course.Level {
		case "easy":
			levelRu = "Легкий"
		case "hard":
			levelRu = "Сложный"
		}

		coursesView = append(coursesView, CourseView{
			ID:          course.ID,
			Title:       course.Title,
			Description: course.Description,
			Level:       course.Level,
			LevelRu:     levelRu,
		})
	}

	// Prepare data for template
	data := fiber.Map{
		"CategoryTitle": category.Title,
		"CategoryID":    category.ID,
		"Courses":       coursesView,
		"Pagination":    pagination,
		"CurrentPage":   page,
		"Level":         level,
		"SortBy":        sortBy,
	}

	slog.Info("Rendering template", "template", "pages/courses", "coursesCount", len(coursesView))

	// Render the template with layout
	return c.Render("pages/courses", data, "layouts/main")
}
