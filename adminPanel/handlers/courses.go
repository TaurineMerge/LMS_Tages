package handlers

import (
	"strconv"
	"strings"

	"adminPanel/handlers/dto/request"
	"adminPanel/handlers/dto/response"
	"adminPanel/middleware"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CourseHandler обрабатывает HTTP-запросы для курсов.
// Содержит сервис для бизнес-логики и методы для маршрутов.
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler создает новый экземпляр CourseHandler.
// Принимает сервис курсов.
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// RegisterRoutes регистрирует маршруты для курсов.
// Создает группу /categories/:category_id/courses и привязывает методы к маршрутам.
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	courses := router.Group("/categories/:category_id/courses")

	courses.Get("/", h.getCourses)
	courses.Post("/", middleware.ValidateJSONSchema("course-create.json"), h.createCourse)
	courses.Get("/:course_id", h.getCourse)
	courses.Put("/:course_id", middleware.ValidateJSONSchema("course-update.json"), h.updateCourse)
	courses.Delete("/:course_id", h.deleteCourse)
}

// getCourses обрабатывает GET /categories/:category_id/courses.
// Возвращает список курсов для категории с фильтрами и пагинацией.
func (h *CourseHandler) getCourses(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}
	filter := request.CourseFilter{
		CategoryID: categoryID,
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter.Page = page
	filter.Limit = limit

	result, err := h.courseService.GetCourses(ctx, filter)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.getCourses.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(result.Data.Items)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(result)
}

// createCourse обрабатывает POST /categories/:category_id/courses.
// Создает новый курс для категории на основе JSON в теле запроса.
func (h *CourseHandler) createCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)
	span.AddEvent("handler.createCourse.start",
		trace.WithAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("category.id", c.Params("category_id")),
		))

	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	var input request.CourseCreate

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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_JSON",
				Message: "Invalid request body",
			},
		})
	}

	input.CategoryID = categoryID

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Visibility must be one of: draft, public, private",
			},
		})
	}

	course, err := h.courseService.CreateCourse(ctx, input)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.createCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.Status(201).JSON(course)
}

// getCourse обрабатывает GET /categories/:category_id/courses/:course_id.
// Возвращает курс по ID в категории.
func (h *CourseHandler) getCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	course, err := h.courseService.GetCourse(ctx, categoryID, id)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.getCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// updateCourse обрабатывает PUT /categories/:category_id/courses/:course_id.
// Обновляет курс по ID в категории на основе JSON в теле запроса.
func (h *CourseHandler) updateCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}
	var input request.CourseUpdate

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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_JSON",
				Message: "Invalid request body",
			},
		})
	}

	if input.CategoryID == "" {
		input.CategoryID = categoryID
	}

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Level must be one of: hard, medium, easy",
			},
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "VALIDATION_ERROR",
				Message: "Visibility must be one of: draft, public, private",
			},
		})
	}

	course, err := h.courseService.UpdateCourse(ctx, categoryID, id, input)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.updateCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// deleteCourse обрабатывает DELETE /categories/:category_id/courses/:course_id.
// Удаляет курс по ID в категории.
func (h *CourseHandler) deleteCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
		return c.Status(400).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}

	err := h.courseService.DeleteCourse(ctx, categoryID, id)
	if err != nil {
		if appErr, ok := err.(*middleware.AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    appErr.Code,
					Message: appErr.Message,
				},
			})
		}
		return c.Status(500).JSON(response.ErrorResponse{
			Status: "error",
			Error: response.ErrorDetails{
				Code:    "SERVER_ERROR",
				Message: "Internal server error",
			},
		})
	}

	span.AddEvent("handler.deleteCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", id),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
}

// isValidLevel проверяет, является ли уровень сложности допустимым.
// Допустимые значения: hard, medium, easy.
func isValidLevel(level string) bool {
	switch strings.ToLower(level) {
	case "hard", "medium", "easy":
		return true
	default:
		return false
	}
}

// isValidVisibility проверяет, является ли видимость допустимой.
// Допустимые значения: draft, public, private.
func isValidVisibility(visibility string) bool {
	switch strings.ToLower(visibility) {
	case "draft", "public", "private":
		return true
	default:
		return false
	}
}
