package handlers

import (
	"fmt"

	"adminPanel/handlers/dto/request"
	"adminPanel/handlers/dto/response"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// LessonHandler - HTTP обработчик для операций с уроками
type LessonHandler struct {
	lessonService *services.LessonService
}

// NewLessonHandler создает новый HTTP обработчик для уроков
func NewLessonHandler(lessonService *services.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: lessonService,
	}
}

// RegisterRoutes регистрирует маршруты для уроков
func (h *LessonHandler) RegisterRoutes(lessons fiber.Router) {
	lessons.Get("/", h.getLessons)
	lessons.Post("/", middleware.ValidateJSONSchema("lesson-create.json"), h.createLesson)
	lessons.Get("/:lesson_id", h.getLesson)
	lessons.Put("/:lesson_id", middleware.ValidateJSONSchema("lesson-update.json"), h.updateLesson)
	lessons.Delete("/:lesson_id", h.deleteLesson)
}

// getLessons обрабатывает GET .../lessons
func (h *LessonHandler) getLessons(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)

	courseID := c.Params("course_id")
	if !isValidUUID(courseID) {
		return middleware.NewAppError("Invalid course ID format", 400, "INVALID_UUID")
	}

	var queryParams models.QueryList
	if err := c.QueryParser(&queryParams); err != nil {
		return middleware.NewAppError(fmt.Sprintf("Invalid query parameters: %v", err), 400, "VALIDATION_ERROR")
	}

	lessonsResponse, err := h.lessonService.GetLessons(ctx, courseID, queryParams)
	if err != nil {
		return err
	}

	span.AddEvent("handler.getLessons.end", trace.WithAttributes(
		attribute.Int("response.count", len(lessonsResponse.Data.Items)),
	))

	return c.JSON(lessonsResponse)
}

// getLesson обрабатывает GET .../lessons/:lesson_id
func (h *LessonHandler) getLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return middleware.NewAppError("Invalid course or lesson ID format", 400, "INVALID_UUID")
	}

	lesson, err := h.lessonService.GetLesson(ctx, lessonID, courseID)
	if err != nil {
		return err
	}

	return c.JSON(lesson)
}

// createLesson обрабатывает POST .../lessons
func (h *LessonHandler) createLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return middleware.NewAppError("Invalid course ID format", 400, "INVALID_UUID")
	}

	var input request.LessonCreate
	if err := c.BodyParser(&input); err != nil {
		return middleware.NewAppError(fmt.Sprintf("Invalid request body: %v", err), 400, "VALIDATION_ERROR")
	}

	lesson, err := h.lessonService.CreateLesson(ctx, courseID, input)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(lesson)
}

// updateLesson обрабатывает PUT .../lessons/:lesson_id
func (h *LessonHandler) updateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return middleware.NewAppError("Invalid course or lesson ID format", 400, "INVALID_UUID")
	}

	var input request.LessonUpdate
	if err := c.BodyParser(&input); err != nil {
		return middleware.NewAppError(fmt.Sprintf("Invalid request body: %v", err), 400, "VALIDATION_ERROR")
	}

	lesson, err := h.lessonService.UpdateLesson(ctx, lessonID, courseID, input)
	if err != nil {
		return err
	}

	return c.JSON(lesson)
}

// deleteLesson обрабатывает DELETE .../lessons/:lesson_id
func (h *LessonHandler) deleteLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return middleware.NewAppError("Invalid course or lesson ID format", 400, "INVALID_UUID")
	}

	err := h.lessonService.DeleteLesson(ctx, lessonID, courseID)
	if err != nil {
		return err
	}

	return c.JSON(response.StatusOnly{Status: "success"})
}
