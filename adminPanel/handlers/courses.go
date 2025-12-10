package handlers

import (
	"strconv"
	"strings"

	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
)

// CourseHandler - обработчик для курсов
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler создает обработчик курсов
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// RegisterRoutes регистрирует маршруты
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	courses := router.Group("/courses")

	courses.Get("/", h.getCourses)
	courses.Post("/", h.createCourse)
	courses.Get("/:id", h.getCourse)
	courses.Put("/:id", h.updateCourse)
	courses.Delete("/:id", h.deleteCourse)
	courses.Get("/:id/lessons", h.getCourseLessons)
}

// GetCourses - получение курсов с фильтрацией
func (h *CourseHandler) getCourses(c *fiber.Ctx) error {
	// Парсим параметры запроса
	filter := models.CourseFilter{
		Level:      c.Query("level"),
		Visibility: c.Query("visibility"),
		CategoryID: c.Query("category_id"),
	}

	// Парсим page и limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter.Page = page
	filter.Limit = limit

	// Валидация
	if filter.CategoryID != "" && !isValidUUID(filter.CategoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Получаем курсы
	result, err := h.courseService.GetCourses(c.Context(), filter)
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

	return c.JSON(result)
}

// CreateCourse - создание курса
func (h *CourseHandler) createCourse(c *fiber.Ctx) error {
	// Валидация входных данных
	var input models.CourseCreate

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

	// Проверка дополнительных условий
	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"level": "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"visibility": "Visibility must be one of: draft, public, private",
			},
		})
	}

	// Устанавливаем значения по умолчанию
	if input.Level == "" {
		input.Level = "medium"
	}
	if input.Visibility == "" {
		input.Visibility = "draft"
	}

	// Создаем курс
	course, err := h.courseService.CreateCourse(c.Context(), input)
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

	return c.Status(201).JSON(course)
}

// GetCourse - получение курса по ID
func (h *CourseHandler) getCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	course, err := h.courseService.GetCourse(c.Context(), id)
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

	return c.JSON(course)
}

// UpdateCourse - обновление курса
func (h *CourseHandler) updateCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Валидация входных данных
	var input models.CourseUpdate

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

	// Проверка дополнительных условий
	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"level": "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(models.ValidationErrorResponse{
			Error: "Validation failed",
			Code:  "VALIDATION_ERROR",
			Errors: map[string]string{
				"visibility": "Visibility must be one of: draft, public, private",
			},
		})
	}

	// Обновляем курс
	course, err := h.courseService.UpdateCourse(c.Context(), id, input)
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

	return c.JSON(course)
}

// DeleteCourse - удаление курса
func (h *CourseHandler) deleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	err := h.courseService.DeleteCourse(c.Context(), id)
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

// GetCourseLessons - получение уроков курса
func (h *CourseHandler) getCourseLessons(c *fiber.Ctx) error {
	id := c.Params("id")

	if !isValidUUID(id) {
		return c.Status(400).JSON(models.ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	// Реализация будет в LessonHandler
	return c.Status(501).JSON(models.ErrorResponse{
		Error: "Use /api/v1/courses/:id/lessons endpoint",
		Code:  "NOT_IMPLEMENTED",
	})
}

// Вспомогательные функции валидации
func isValidLevel(level string) bool {
	switch strings.ToLower(level) {
	case "hard", "medium", "easy":
		return true
	default:
		return false
	}
}

func isValidVisibility(visibility string) bool {
	switch strings.ToLower(visibility) {
	case "draft", "public", "private":
		return true
	default:
		return false
	}
}
