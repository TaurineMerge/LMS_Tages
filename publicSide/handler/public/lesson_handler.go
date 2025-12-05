package public

import (
	"strconv"

	// Импорты ВАШЕГО проекта

	"github.com/TaurineMerge/LMS_Tages/publicSide/service"

	// Внешние зависимости
	"github.com/gofiber/fiber/v3"
)

// LessonHandler handles public lesson requests
type LessonHandler struct {
	lessonService service.LessonService
}

// NewLessonHandler creates new lesson handler
func NewLessonHandler(lessonService service.LessonService) *LessonHandler {
	return &LessonHandler{lessonService: lessonService}
}

// RegisterRoutes registers lesson routes
func (h *LessonHandler) RegisterRoutes(app *fiber.App) {
	public := app.Group("/api/v1/public/lessons")

	// @Summary Get all published lessons
	// @Description Get list of all published lessons
	// @Tags Public Lessons
	// @Accept json
	// @Produce json
	// @Param course_id query int false "Filter by course ID"
	// @Success 200 {array} entity.Lesson
	// @Router /public/lessons [get]
	public.Get("/", h.GetAllLessons)

	// @Summary Get lesson by ID
	// @Description Get published lesson by ID with course info
	// @Tags Public Lessons
	// @Accept json
	// @Produce json
	// @Param id path int true "Lesson ID"
	// @Success 200 {object} entity.Lesson
	// @Router /public/lessons/{id} [get]
	public.Get("/:id", h.GetLessonByID)
}

// GetAllLessons returns all published lessons
// @Summary Get all published lessons
// @Description Get list of all published lessons
// @Tags Public Lessons
// @Accept json
// @Produce json
// @Param course_id query int false "Filter by course ID"
// @Success 200 {array} entity.Lesson
// @Router /api/v1/public/lessons [get]
func (h *LessonHandler) GetAllLessons(c fiber.Ctx) error {
	var courseID int64 = 0

	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if id, err := strconv.ParseInt(courseIDStr, 10, 64); err == nil {
			courseID = id
		}
	}

	lessons, err := h.lessonService.GetAllPublished(c.Context(), courseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch lessons",
		})
	}

	return c.JSON(lessons)
}

// GetLessonByID returns lesson by ID with course info
// @Summary Get lesson by ID
// @Description Get published lesson by ID with course info
// @Tags Public Lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID"
// @Success 200 {object} entity.Lesson
// @Router /api/v1/public/lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid lesson ID",
		})
	}

	lesson, err := h.lessonService.GetWithCourse(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Lesson not found",
		})
	}

	return c.JSON(lesson)
}
