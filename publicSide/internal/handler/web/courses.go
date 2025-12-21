package web

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/breadcrumbs"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CoursesHandler handles web requests for the courses page.
type CoursesHandler struct {
	courseService service.CourseService
	lessonService service.LessonService
}

// NewCoursesHandler creates a new instance of CoursesHandler.
func NewCoursesHandler(courseService service.CourseService, lessonService service.LessonService) *CoursesHandler {
	return &CoursesHandler{
		courseService: courseService,
		lessonService: lessonService,
	}
}



// RenderCourses renders the courses page with filters and sorting.
func (h *CoursesHandler) RenderCourses(c *fiber.Ctx) error {
	vm, err := h.buildCoursesPageViewModel(c)
	if err != nil {
		slog.Error("Failed to build courses page view model", "error", err)
		// Render a generic error page, or let the global error handler manage it
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"title":   "Error",
			"message": "Could not load the courses page.",
		}, "layouts/main")
	}
	return c.Render("pages/courses", vm, "layouts/main")
}

// RenderCoursePage renders the individual course page with course details.
func (h *CoursesHandler) RenderCoursePage(c *fiber.Ctx) error {
	vm, err := h.buildCoursePageViewModel(c)
	if err != nil {
		slog.Error("Failed to build course page view model", "error", err)
		return c.Status(fiber.StatusInternalServerError).Render("pages/error", fiber.Map{
			"title":   "Error",
			"message": "Could not load the course page.",
		}, "layouts/main")
	}
	return c.Render("pages/course", vm, "layouts/main")
}

func (h *CoursesHandler) buildCoursesPageViewModel(c *fiber.Ctx) (viewmodel.CoursesPageViewModel, error) {
	vm := viewmodel.CoursesPageViewModel{}
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		// This should be handled by a validation middleware ideally
		return vm, err
	}
	vm.CategoryID = categoryID

	// Parse query parameters
	vm.CurrentPage = c.QueryInt("page", 1)
	vm.Level = c.Query("level", "all")
	vm.SortBy = c.Query("sort_by", "updated_desc")
	limit := c.QueryInt("limit", 28)

	// Get category information
	category, err := h.courseService.GetCategoryByID(c.UserContext(), categoryID)
	if err != nil {
		return vm, err
	}
	vm.CategoryTitle = category.Title

	// Build breadcrumbs
	vm.PageHeader = viewmodel.PageHeaderViewModel{
		Title:       category.Title,
		Breadcrumbs: breadcrumbs.ForCoursesPage(*category),
	}

	// Get courses
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(
		c.UserContext(), categoryID, vm.CurrentPage, limit, vm.Level, vm.SortBy,
	)
	if err != nil {
		return vm, err
	}
	vm.Pagination = pagination

	// Transform courses to view model
	vm.Courses = make([]viewmodel.CourseView, len(courses))
	for i, course := range courses {
		levelRu := "Средний"
		switch course.Level {
		case "easy":
			levelRu = "Легкий"
		case "hard":
			levelRu = "Сложный"
		}
		vm.Courses[i] = viewmodel.CourseView{
			ID:          course.ID,
			Title:       course.Title,
			Description: course.Description,
			Level:       course.Level,
			LevelRu:     levelRu,
		}
	}

	return vm, nil
}

func (h *CoursesHandler) buildCoursePageViewModel(c *fiber.Ctx) (viewmodel.CoursePageViewModel, error) {
	vm := viewmodel.CoursePageViewModel{}
	categoryID := c.Params(routing.PathVariableCategoryID)
	courseID := c.Params(routing.PathVariableCourseID)

	if _, err := uuid.Parse(categoryID); err != nil {
		return vm, err
	}
	if _, err := uuid.Parse(courseID); err != nil {
		return vm, err
	}
	vm.CategoryID = categoryID

	// Get category and course info
	category, err := h.courseService.GetCategoryByID(c.UserContext(), categoryID)
	if err != nil {
		return vm, err
	}
	vm.CategoryTitle = category.Title

	course, err := h.courseService.GetCourseByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return vm, err
	}
	vm.Course = course

	// Build breadcrumbs
	vm.PageHeader = viewmodel.PageHeaderViewModel{
		Title:       course.Title,
		Breadcrumbs: breadcrumbs.ForCoursePage(*category, course),
	}

	// Localized level
	vm.LevelRu = "Средний"
	switch course.Level {
	case "easy":
		vm.LevelRu = "Легкий"
	case "hard":
		vm.LevelRu = "Сложный"
	}

	// First lesson ID
	lessons, _, err := h.lessonService.GetAllByCourseID(c.UserContext(), categoryID, courseID, 1, 1, "created_at")
	if err != nil {
		slog.Warn("Failed to get first lesson for course page", "courseId", courseID, "error", err)
	} else if len(lessons) > 0 {
		vm.FirstLessonID = lessons[0].ID
	}

	return vm, nil
}
