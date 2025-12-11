package handlers

import (
	"strconv"
	"strings"

	"adminPanel/exceptions"
	"adminPanel/middleware"
	"adminPanel/models"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

// RegisterRoutes регистрирует маршруты вида /categories/:category_id/courses
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	courses := router.Group("/categories/:category_id/courses")

	courses.Get("/", h.getCourses)
	courses.Post("/", h.createCourse)
	courses.Get("/:course_id", h.getCourse)
	courses.Put("/:course_id", h.updateCourse)
	courses.Delete("/:course_id", h.deleteCourse)
}

// GetCourses - получение курсов с фильтрацией
func (h *CourseHandler) getCourses(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCourses.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("http.query", c.Context().QueryArgs().String()),
		))

	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	// Парсим параметры запроса
	filter := models.CourseFilter{
		CategoryID: categoryID,
	}
	// level/visibility фильтры временно отключены (см. swagger)
	// filter.Level = ""
	// filter.Visibility = ""

	// Парсим page и limit
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter.Page = page
	filter.Limit = limit

	// Валидация
	// Получаем курсы
	result, err := h.courseService.GetCourses(ctx, filter)
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
	span.AddEvent("handler.getCourses.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(result.Data.Items)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(result)
}

// CreateCourse - создание курса
func (h *CourseHandler) createCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.createCourse.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	// Валидация входных данных
	var input models.CourseCreate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.createCourse.request_body",
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

	// Привязываем категорию из пути
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

	// Проверка дополнительных условий
	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Visibility must be one of: draft, public, private",
			},
		})
	}

	// Создаем курс
	course, err := h.courseService.CreateCourse(ctx, input)
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
	span.AddEvent("handler.createCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.Status(201).JSON(course)
}

// GetCourse - получение курса по ID
func (h *CourseHandler) getCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.getCourse.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
		))

	categoryID := c.Params("category_id")
	id := c.Params("course_id")

	if !isValidUUID(id) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	course, err := h.courseService.GetCourse(ctx, categoryID, id)
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
	span.AddEvent("handler.getCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// UpdateCourse - обновление курса
func (h *CourseHandler) updateCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.updateCourse.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
		))

	categoryID := c.Params("category_id")
	id := c.Params("course_id")

	if !isValidUUID(id) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	// Валидация входных данных
	var input models.CourseUpdate

	// Логируем тело запроса
	if len(c.Body()) > 0 {
		body := c.Body()
		const maxLoggedBody = 2048
		if len(body) > maxLoggedBody {
			body = body[:maxLoggedBody]
		}
		span.AddEvent("handler.updateCourse.request_body",
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

	// Привязываем категорию из пути
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

	// Проверка дополнительных условий
	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Visibility must be one of: draft, public, private",
			},
		})
	}

	// Обновляем курс
	course, err := h.courseService.UpdateCourse(ctx, categoryID, id, input)
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
	span.AddEvent("handler.updateCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// DeleteCourse - удаление курса
func (h *CourseHandler) deleteCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Логируем вызов метода
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.deleteCourse.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
			attribute.String("course.id", c.Params("course_id")),
		))

	categoryID := c.Params("category_id")
	id := c.Params("course_id")

	if !isValidUUID(id) || !isValidUUID(categoryID) {
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	err := h.courseService.DeleteCourse(ctx, categoryID, id)
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
	span.AddEvent("handler.deleteCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", id),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
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
