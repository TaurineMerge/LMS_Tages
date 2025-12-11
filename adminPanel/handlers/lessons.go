package handlers

import (
	// "strings"

	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	lessons := router.Group("/categories/:category_id/courses/:course_id")

	lessons.Get("/lessons", h.getLessons)
	lessons.Post("/lessons", h.createLesson)
	lessons.Get("/lessons/:lesson_id", h.getLesson)
	lessons.Put("/lessons/:lesson_id", h.updateLesson)
	lessons.Delete("/lessons/:lesson_id", h.deleteLesson)
}

// GetLessons - получение уроков курса
func (h *LessonHandler) getLessons(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getLessons.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
			attribute.String("http.query", c.Context().QueryArgs().String()),
		))

	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid course or category ID format",
			},
		})
	}

	lessons, pagination, err := h.lessonService.GetLessons(ctx, categoryID, courseID, c.Query("page"), c.Query("limit"))
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	response := models.LessonListResponse{
		Status: "success",
	}
	response.Data.Items = lessons
	response.Data.Pagination = pagination

	// Логируем успешное завершение
	span.AddEvent("handler.getLessons.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(response.Data.Items)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(response)
}

// GetLesson - получение урока по ID
func (h *LessonHandler) getLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getLesson.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
			attribute.String("lesson.id", c.Params("lesson_id")),
		))

	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	lesson, err := h.lessonService.GetLesson(ctx, lessonID, courseID, categoryID)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.getLesson.end",
		trace.WithAttributes(
			attribute.String("lesson.id", lesson.Data.ID),
			attribute.String("lesson.title", lesson.Data.Title),
			attribute.String("response.status", "success"),
		))

	lesson.Status = "success"
	return c.JSON(lesson)
}

// CreateLesson - создание урока
func (h *LessonHandler) createLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.createLesson.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
		))

	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid course or category ID format",
			},
		})
	}

	// Валидация входных данных
	var input models.LessonCreate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.createLesson.request_body",
			trace.WithAttributes(
				attribute.String("request.body", string(body)),
			))
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request body",
			},
		})
	}

	input.CategoryID = categoryID

	// Валидация через middleware
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Validation error",
			},
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Validation failed",
			},
			Errors: validationErrors,
		})
	}

	// Дополнительная валидация
	if len(input.Title) > 255 {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Title must be less than 255 characters",
			},
		})
	}

	// Создаем урок
	lesson, err := h.lessonService.CreateLesson(ctx, courseID, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.createLesson.end",
		trace.WithAttributes(
			attribute.String("lesson.id", lesson.Data.ID),
			attribute.String("lesson.title", lesson.Data.Title),
			attribute.String("response.status", "success"),
		))

	lesson.Status = "success"
	return c.Status(201).JSON(lesson)
}

// UpdateLesson - обновление урока
func (h *LessonHandler) updateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.updateLesson.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
			attribute.String("lesson.id", c.Params("lesson_id")),
		))

	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	// Валидация входных данных
	var input models.LessonUpdate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.updateLesson.request_body",
			trace.WithAttributes(
				attribute.String("request.body", string(body)),
			))
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid request body",
			},
		})
	}

	if input.CategoryID == "" {
		input.CategoryID = categoryID
	}

	// Валидация через middleware
	if validationErrors, err := middleware.ValidateStruct(&input); err != nil {
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Validation error",
			},
		})
	} else if len(validationErrors) > 0 {
		return c.Status(422).JSON(models.ValidationErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Validation failed",
			},
			Errors: validationErrors,
		})
	}

	// Дополнительная валидация
	if input.Title != "" && len(input.Title) > 255 {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Title must be less than 255 characters",
			},
		})
	}

	// Обновляем урок
	lesson, err := h.lessonService.UpdateLesson(ctx, lessonID, courseID, categoryID, input)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.updateLesson.end",
		trace.WithAttributes(
			attribute.String("lesson.id", lesson.Data.ID),
			attribute.String("lesson.title", lesson.Data.Title),
			attribute.String("response.status", "success"),
		))

	lesson.Status = "success"
	return c.JSON(lesson)
}

// DeleteLesson - удаление урока
func (h *LessonHandler) deleteLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.deleteLesson.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
			attribute.String("lesson.id", c.Params("lesson_id")),
		))

	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	err := h.lessonService.DeleteLesson(ctx, lessonID, courseID, categoryID)
	if err != nil {
		if appErr, ok := err.(*exceptions.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	// Логируем успешное завершение
	span.AddEvent("handler.deleteLesson.end",
		trace.WithAttributes(
			attribute.String("lesson.id", lessonID),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
}
