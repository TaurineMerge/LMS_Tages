package public

import (
	"strconv"

	// Импорты ВАШЕГО проекта
	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"github.com/TaurineMerge/LMS_Tages/publicSide/service"

	// Внешние зависимости
	"github.com/gofiber/fiber/v3"
)

// CourseHandler handles public course requests
type CourseHandler struct {
	courseService service.CourseService
}

// NewCourseHandler creates new course handler
func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// RegisterRoutes registers course routes
func (h *CourseHandler) RegisterRoutes(app *fiber.App) {
	public := app.Group("/api/v1/public/courses")

	// @Summary Get all published courses
	// @Description Get list of all published courses
	// @Tags Public Courses
	// @Accept json
	// @Produce json
	// @Success 200 {array} entity.CourseResponse
	// @Router /public/courses [get]
	public.Get("/", h.GetAllCourses)

	// @Summary Get course by ID
	// @Description Get published course by ID with lessons
	// @Tags Public Courses
	// @Accept json
	// @Produce json
	// @Param id path int true "Course ID"
	// @Success 200 {object} entity.Course
	// @Router /public/courses/{id} [get]
	public.Get("/:id", h.GetCourseByID)
}

// GetAllCourses returns all published courses
// @Summary Get all published courses
// @Description Get list of all published courses
// @Tags Public Courses
// @Accept json
// @Produce json
// @Success 200 {array} entity.CourseResponse
// @Router /api/v1/public/courses [get]
func (h *CourseHandler) GetAllCourses(c fiber.Ctx) error {
	courses, err := h.courseService.GetAllPublished(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch courses",
		})
	}

	// Convert to response format
	response := make([]entity.CourseResponse, len(courses))
	for i, course := range courses {
		response[i] = entity.CourseResponse{
			ID:           course.ID,
			Title:        course.Title,
			Description:  course.Description,
			Slug:         course.Slug,
			LessonsCount: len(course.Lessons),
			CreatedAt:    course.CreatedAt,
		}
	}

	return c.JSON(response)
}

// GetCourseByID returns course by ID with lessons
// @Summary Get course by ID
// @Description Get published course by ID with lessons
// @Tags Public Courses
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} entity.Course
// @Router /api/v1/public/courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid course ID",
		})
	}

	course, err := h.courseService.GetWithLessons(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Course not found",
		})
	}

	return c.JSON(course)
}
