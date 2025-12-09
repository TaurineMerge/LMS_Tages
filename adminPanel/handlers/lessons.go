package handlers

import (
	// "strings"

	"github.com/gofiber/fiber/v2"
	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"
)

// LessonHandler - обработчик для уроков
type LessonHandler struct {
	lessonService *services.LessonService
}

// NewLessonHandler создает обработчик уроков
func NewLessonHandler(lessonService *services.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: lessonService,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *LessonHandler) RegisterRoutes(router fiber.Router) {
	// Уроки курса
	courses := router.Group("/courses/:course_id")
	
	courses.Get("/lessons", h.getLessons)
	courses.Post("/lessons", h.createLesson)
	courses.Get("/lessons/:id", h.getLesson)
	courses.Put("/lessons/:id", h.updateLesson)
	courses.Delete("/lessons/:id", h.deleteLesson)
}

// GetLessons - получение уроков курса
func (h *LessonHandler) getLessons(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	
	if !isValidUUID(courseID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	lessons, err := h.lessonService.GetLessons(c.Context(), courseID)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(lessons)
}

// GetLesson - получение урока по ID
func (h *LessonHandler) getLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("id")
	
	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	lesson, err := h.lessonService.GetLesson(c.Context(), lessonID, courseID)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(lesson)
}

// CreateLesson - создание урока
func (h *LessonHandler) createLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	
	if !isValidUUID(courseID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Валидация входных данных
	var input models.LessonCreate
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}
	
	// Валидация через middleware
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Validation error",
			Code:  "INTERNAL_ERROR",
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Error:  "Validation failed",
			Code:   "VALIDATION_ERROR",
			Errors: validationErrors,
		})
	}

	// Дополнительная валидация
	if len(input.Title) > 255 {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"title": "Title must be less than 255 characters",
			},
		})
	}

	// Создаем урок
	lesson, err := h.lessonService.CreateLesson(c.Context(), courseID, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.Status(201).JSON(lesson)
}

// UpdateLesson - обновление урока
func (h *LessonHandler) updateLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("id")
	
	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Валидация входных данных
	var input models.LessonUpdate
	
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}
	
	// Валидация через middleware
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Validation error",
			Code:  "INTERNAL_ERROR",
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Error:  "Validation failed",
			Code:   "VALIDATION_ERROR",
			Errors: validationErrors,
		})
	}

	// Дополнительная валидация
	if input.Title != "" && len(input.Title) > 255 {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"title": "Title must be less than 255 characters",
			},
		})
	}

	// Обновляем урок
	lesson, err := h.lessonService.UpdateLesson(c.Context(), lessonID, courseID, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.JSON(lesson)
}

// DeleteLesson - удаление урока
func (h *LessonHandler) deleteLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("id")
	
	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	err := h.lessonService.DeleteLesson(c.Context(), lessonID, courseID)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
	}

	return c.SendStatus(204)
}