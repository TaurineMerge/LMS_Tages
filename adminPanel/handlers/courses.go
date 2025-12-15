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

// CourseHandler - HTTP обработчик для операций с курсами
//
// Обработчик предоставляет REST API для управления курсами:
//   - GET /categories/:category_id/courses - получение курсов категории
//   - POST /categories/:category_id/courses - создание курса
//   - GET /categories/:category_id/courses/:id - получение курса
//   - PUT /categories/:category_id/courses/:id - обновление курса
//   - DELETE /categories/:category_id/courses/:id - удаление курса
//
// Особенности:
//   - Фильтрация курсов по уровню и видимости
//   - Пагинация результатов
//   - Валидация входных данных
//   - Интеграция с OpenTelemetry для трассировки
//   - Стандартизированный формат ответов
type CourseHandler struct {
	courseService *services.CourseService
}

// NewCourseHandler создает новый HTTP обработчик для курсов
//
// Параметры:
//   - courseService: сервис для работы с курсами
//
// Возвращает:
//   - *CourseHandler: указатель на новый обработчик
func NewCourseHandler(courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// RegisterRoutes регистрирует маршруты для курсов
//
// Регистрирует все необходимые маршруты в указанном роутере.
// Все маршруты включают в себя category_id в пути.
//
// Параметры:
//   - router: роутер Fiber для регистрации маршрутов
func (h *CourseHandler) RegisterRoutes(router fiber.Router) {
	courses := router.Group("/categories/:category_id/courses")

	courses.Get("/", h.getCourses)
	courses.Post("/", h.createCourse)
	courses.Get("/:course_id", h.getCourse)
	courses.Put("/:course_id", h.updateCourse)
	courses.Delete("/:course_id", h.deleteCourse)
}

// getCourses обрабатывает GET /categories/:category_id/courses
//
// Возвращает список курсов категории с возможностью фильтрации
// по уровню сложности и видимости, а также с пагинацией.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
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
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}
	filter := models.CourseFilter{
		CategoryID: categoryID,
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter.Page = page
	filter.Limit = limit

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

	span.AddEvent("handler.getCourses.end",
		trace.WithAttributes(
			attribute.Int("response.count", len(result.Data.Items)),
			attribute.String("response.status", "success"),
		))

	return c.JSON(result)
}

// createCourse обрабатывает POST /categories/:category_id/courses
//
// Создает новый курс в указанной категории с валидацией данных.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
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
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid category ID format",
			},
		})
	}

	var input models.CourseCreate

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

	input.CategoryID = categoryID

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

	span.AddEvent("handler.createCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.Status(201).JSON(course)
}

// getCourse обрабатывает GET /categories/:category_id/courses/:id
//
// Возвращает курс по уникальному идентификатору.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
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

	span.AddEvent("handler.getCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// updateCourse обрабатывает PUT /categories/:category_id/courses/:id
//
// Обновляет существующий курс с валидацией данных.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
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
		return c.Status(400).JSON(models.ErrorResponse{
			Status: "error",
			Error: models.ErrorDetails{
				Code:    "INVALID_UUID",
				Message: "Invalid ID format",
			},
		})
	}
	var input models.CourseUpdate

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

	span.AddEvent("handler.updateCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", course.Data.ID),
			attribute.String("course.title", course.Data.Title),
			attribute.String("response.status", "success"),
		))

	return c.JSON(course)
}

// deleteCourse обрабатывает DELETE /categories/:category_id/courses/:id
//
// Удаляет курс по уникальному идентификатору.
//
// Параметры:
//   - c: контекст Fiber
//
// Возвращает:
//   - error: ошибка выполнения (если есть)
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

	span.AddEvent("handler.deleteCourse.end",
		trace.WithAttributes(
			attribute.String("course.id", id),
			attribute.String("response.status", "success"),
		))

	return c.SendStatus(204)
}

// isValidLevel проверяет валидность уровня сложности
//
// Проверяет, что уровень сложности является одним из допустимых:
// "hard", "medium", "easy" (регистронезависимо).
//
// Параметры:
//   - level: уровень сложности для проверки
//
// Возвращает:
//   - bool: true, если уровень валиден
func isValidLevel(level string) bool {
	switch strings.ToLower(level) {
	case "hard", "medium", "easy":
		return true
	default:
		return false
	}
}

// isValidVisibility проверяет валидность видимости
//
// Проверяет, что видимость является одной из допустимых:
// "draft", "public", "private" (регистронезависимо).
//
// Параметры:
//   - visibility: видимость для проверки
//
// Возвращает:
//   - bool: true, если видимость валидна
func isValidVisibility(visibility string) bool {
	switch strings.ToLower(visibility) {
	case "draft", "public", "private":
		return true
	default:
		return false
	}
}
