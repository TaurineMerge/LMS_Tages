package web

import (
	"adminPanel/handlers/dto/request"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
)

// CourseView представляет курс для отображения
type CourseView struct {
	ID          string
	CategoryID  string
	Title       string
	Description string
	Level       string
	LevelRu     string
	Visible     bool
	CreatedAt   string
	UpdatedAt   string
}

// levelToRussian преобразует уровень сложности в русский текст
func levelToRussian(level string) string {
	switch level {
	case "easy":
		return "Лёгкий"
	case "medium":
		return "Средний"
	case "hard":
		return "Сложный"
	default:
		return level
	}
}

// CourseWebHandler обрабатывает веб-страницы для управления курсами
type CourseWebHandler struct {
	courseService   *services.CourseService
	categoryService *services.CategoryService
}

// NewCourseWebHandler создает новый обработчик веб-страниц курсов
func NewCourseWebHandler(courseService *services.CourseService, categoryService *services.CategoryService) *CourseWebHandler {
	return &CourseWebHandler{
		courseService:   courseService,
		categoryService: categoryService,
	}
}

// RenderCoursesEditor отображает страницу со списком курсов категории
func (h *CourseWebHandler) RenderCoursesEditor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")

	// Получаем параметры фильтрации
	levelFilter := c.Query("level", "all")
	visibilityFilter := c.Query("visibility", "all")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/courses-editor", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	// Получаем курсы категории
	filter := request.CourseFilter{
		CategoryID: categoryID,
	}
	// Применяем фильтр по уровню
	if levelFilter != "all" {
		filter.Level = levelFilter
	}
	// Применяем фильтр по видимости
	if visibilityFilter != "all" {
		filter.Visibility = visibilityFilter
	}

	coursesResp, err := h.courseService.GetCourses(ctx, filter)
	if err != nil {
		return c.Status(500).Render("pages/courses-editor", fiber.Map{
			"title":        "Курсы категории: " + category.Title,
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Ошибка загрузки курсов",
		}, "layouts/main")
	}

	// Получаем общее количество курсов без фильтров для отображения total
	totalFilter := request.CourseFilter{CategoryID: categoryID}
	totalResp, _ := h.courseService.GetCourses(ctx, totalFilter)
	totalCount := 0
	if totalResp != nil {
		totalCount = len(totalResp.Data.Items)
	}

	// Преобразуем в CourseView
	courseViews := make([]CourseView, 0, len(coursesResp.Data.Items))
	for _, course := range coursesResp.Data.Items {
		courseViews = append(courseViews, CourseView{
			ID:          course.ID,
			CategoryID:  course.CategoryID,
			Title:       course.Title,
			Description: course.Description,
			Level:       course.Level,
			LevelRu:     levelToRussian(course.Level),
			Visible:     course.Visibility == "public",
			CreatedAt:   formatDateTime(course.CreatedAt),
			UpdatedAt:   formatDateTime(course.UpdatedAt),
		})
	}

	return c.Render("pages/courses-editor", fiber.Map{
		"title":            "Курсы категории: " + category.Title,
		"categoryID":       categoryID,
		"categoryName":     category.Title,
		"courses":          courseViews,
		"coursesCount":     totalCount,
		"levelFilter":      levelFilter,
		"visibilityFilter": visibilityFilter,
	}, "layouts/main")
}

// RenderNewCourseForm отображает форму создания нового курса
func (h *CourseWebHandler) RenderNewCourseForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/course-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	return c.Render("pages/course-form", fiber.Map{
		"title":        "Новый курс",
		"categoryID":   categoryID,
		"categoryName": category.Title,
	}, "layouts/main")
}

// RenderEditCourseForm отображает форму редактирования курса
func (h *CourseWebHandler) RenderEditCourseForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/course-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	// Получаем курс
	course, err := h.courseService.GetCourse(ctx, categoryID, courseID)
	if err != nil {
		return c.Status(404).Render("pages/course-form", fiber.Map{
			"title":        "Курс не найден",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Курс с указанным ID не найден",
		}, "layouts/main")
	}

	courseView := CourseView{
		ID:          course.Data.ID,
		CategoryID:  course.Data.CategoryID,
		Title:       course.Data.Title,
		Description: course.Data.Description,
		Level:       course.Data.Level,
		Visible:     course.Data.Visibility == "public",
		CreatedAt:   formatDateTime(course.Data.CreatedAt),
		UpdatedAt:   formatDateTime(course.Data.UpdatedAt),
	}

	return c.Render("pages/course-form", fiber.Map{
		"title":        "Редактировать курс",
		"categoryID":   categoryID,
		"categoryName": category.Title,
		"course":       courseView,
	}, "layouts/main")
}

// CreateCourse обрабатывает создание нового курса
func (h *CourseWebHandler) CreateCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/course-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	level := c.FormValue("level")
	visibleStr := c.FormValue("visible")

	if title == "" {
		return c.Status(400).Render("pages/course-form", fiber.Map{
			"title":        "Новый курс",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Название курса не может быть пустым",
		}, "layouts/main")
	}

	visibility := "draft"
	if visibleStr == "on" {
		visibility = "public"
	}

	// Создаем курс
	input := request.CourseCreate{
		Title:       title,
		Description: description,
		Level:       level,
		CategoryID:  categoryID,
		Visibility:  visibility,
	}

	_, err = h.courseService.CreateCourse(ctx, input)
	if err != nil {
		return c.Status(400).Render("pages/course-form", fiber.Map{
			"title":        "Новый курс",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Ошибка создания курса: " + err.Error(),
		}, "layouts/main")
	}

	// Перенаправляем на список курсов
	return c.Redirect("/admin/categories/" + categoryID + "/courses")
}

// UpdateCourse обрабатывает обновление курса
func (h *CourseWebHandler) UpdateCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/course-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	title := c.FormValue("title")
	description := c.FormValue("description")
	level := c.FormValue("level")
	visibleStr := c.FormValue("visible")

	if title == "" {
		// Получаем курс для отображения в форме
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)
		var courseView *CourseView
		if course != nil {
			courseView = &CourseView{
				ID:          course.Data.ID,
				CategoryID:  course.Data.CategoryID,
				Title:       course.Data.Title,
				Description: course.Data.Description,
				Level:       course.Data.Level,
				Visible:     course.Data.Visibility == "public",
				CreatedAt:   formatDateTime(course.Data.CreatedAt),
				UpdatedAt:   formatDateTime(course.Data.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/course-form", fiber.Map{
			"title":        "Редактировать курс",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"course":       courseView,
			"error":        "Название курса не может быть пустым",
		}, "layouts/main")
	}

	visibility := "draft"
	if visibleStr == "on" {
		visibility = "public"
	}

	// Обновляем курс
	input := request.CourseUpdate{
		Title:       title,
		Description: description,
		Level:       level,
		CategoryID:  categoryID,
		Visibility:  visibility,
	}

	_, err = h.courseService.UpdateCourse(ctx, categoryID, courseID, input)
	if err != nil {
		// Получаем курс для отображения в форме
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)
		var courseView *CourseView
		if course != nil {
			courseView = &CourseView{
				ID:          course.Data.ID,
				CategoryID:  course.Data.CategoryID,
				Title:       course.Data.Title,
				Description: course.Data.Description,
				Level:       course.Data.Level,
				Visible:     course.Data.Visibility == "public",
				CreatedAt:   formatDateTime(course.Data.CreatedAt),
				UpdatedAt:   formatDateTime(course.Data.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/course-form", fiber.Map{
			"title":        "Редактировать курс",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"course":       courseView,
			"error":        "Ошибка обновления курса: " + err.Error(),
		}, "layouts/main")
	}

	// Перенаправляем на список курсов
	return c.Redirect("/admin/categories/" + categoryID + "/courses")
}

// DeleteCourse обрабатывает удаление курса
func (h *CourseWebHandler) DeleteCourse(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	err := h.courseService.DeleteCourse(ctx, categoryID, courseID)
	if err != nil {
		// Можно добавить flash-сообщение об ошибке
		return c.Redirect("/admin/categories/" + categoryID + "/courses")
	}

	// Перенаправляем на список курсов
	return c.Redirect("/admin/categories/" + categoryID + "/courses")
}
