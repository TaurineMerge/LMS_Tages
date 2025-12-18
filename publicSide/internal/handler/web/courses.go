package web

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CoursesHandler handles web requests for the courses page.
type CoursesHandler struct {
	courseService service.CourseService
}

// NewCoursesHandler creates a new instance of CoursesHandler.
func NewCoursesHandler(courseService service.CourseService) *CoursesHandler {
	return &CoursesHandler{
		courseService: courseService,
	}
}

// RenderCourses renders the courses page with filters and sorting.
func (h *CoursesHandler) RenderCourses(c *fiber.Ctx) error {
	// Get category ID from URL params
	categoryID := c.Params(apiconst.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("pages/error", fiber.Map{
			"title":   "Invalid Category",
			"message": "The category ID is invalid.",
		}, "layouts/main")
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
		return c.Status(fiber.StatusNotFound).Render("pages/error", fiber.Map{
			"title":   "Category Not Found",
			"message": "The requested category could not be found.",
		}, "layouts/main")
	}
	slog.Info("Category found", "categoryId", categoryID, "title", category.Title)

	// Get courses for this category with pagination and filters
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(
		c.UserContext(),
		categoryID,
		page,
		limit,
		level,
		sortBy,
	)

	if err != nil {
		slog.Error("Failed to get courses", "categoryId", categoryID, "error", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"title":   "Error Loading Courses",
			"message": "An error occurred while loading the courses.",
		}, "layouts/main")
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

	return c.Render("pages/courses", data, "layouts/main")
}
