// Пакет web содержит обработчики для веб-интерфейса админ-панели.
package web

import (
	"adminPanel/handlers/dto/request"
	"adminPanel/models"
	"adminPanel/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

// LessonView представляет урок для отображения в веб-интерфейсе.
type LessonView struct {
	ID        string
	CourseID  string
	Title     string
	Content   string
	Number    int
	CreatedAt string
	UpdatedAt string
}

// LessonWebHandler обрабатывает веб-страницы для управления уроками.
type LessonWebHandler struct {
	lessonService   *services.LessonService
	courseService   *services.CourseService
	categoryService *services.CategoryService
}

// NewLessonWebHandler создает новый обработчик веб-страниц уроков.
func NewLessonWebHandler(
	lessonService *services.LessonService,
	courseService *services.CourseService,
	categoryService *services.CategoryService,
) *LessonWebHandler {
	return &LessonWebHandler{
		lessonService:   lessonService,
		courseService:   courseService,
		categoryService: categoryService,
	}
}

// RenderLessonsEditor отображает страницу со списком уроков курса.
func (h *LessonWebHandler) RenderLessonsEditor(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/lessons-editor", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	course, err := h.courseService.GetCourse(ctx, categoryID, courseID)
	if err != nil {
		return c.Status(404).Render("pages/lessons-editor", fiber.Map{
			"title":        "Курс не найден",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Курс с указанным ID не найден",
		}, "layouts/main")
	}

	queryParams := models.QueryList{
		Page:  1,
		Limit: 100,
	}

	lessonsResp, err := h.lessonService.GetLessons(ctx, courseID, queryParams)
	if err != nil {
		return c.Status(500).Render("pages/lessons-editor", fiber.Map{
			"title":        "Уроки курса: " + course.Data.Title,
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"error":        "Ошибка загрузки уроков",
		}, "layouts/main")
	}

	lessonViews := make([]LessonView, 0, len(lessonsResp.Data.Items))
	for i, lesson := range lessonsResp.Data.Items {
		log.Printf("[DEBUG] Lesson #%d: ID=%s, Title=%s", i+1, lesson.ID, lesson.Title)
		lessonViews = append(lessonViews, LessonView{
			ID:        lesson.ID,
			CourseID:  lesson.CourseID,
			Title:     lesson.Title,
			Number:    i + 1,
			CreatedAt: formatDateTime(lesson.CreatedAt),
			UpdatedAt: formatDateTime(lesson.UpdatedAt),
		})
	}

	return c.Render("pages/lessons-editor", fiber.Map{
		"title":        "Уроки курса: " + course.Data.Title,
		"categoryID":   categoryID,
		"categoryName": category.Title,
		"courseID":     courseID,
		"courseName":   course.Data.Title,
		"lessons":      lessonViews,
		"lessonsCount": len(lessonViews),
	}, "layouts/main")
}

// RenderNewLessonForm отображает форму создания нового урока.
func (h *LessonWebHandler) RenderNewLessonForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/lesson-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	course, err := h.courseService.GetCourse(ctx, categoryID, courseID)
	if err != nil {
		return c.Status(404).Render("pages/lesson-form", fiber.Map{
			"title":        "Курс не найден",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Курс с указанным ID не найден",
		}, "layouts/main")
	}

	return c.Render("pages/lesson-form", fiber.Map{
		"title":        "Новый урок",
		"categoryID":   categoryID,
		"categoryName": category.Title,
		"courseID":     courseID,
		"courseName":   course.Data.Title,
	}, "layouts/main")
}

// RenderEditLessonForm отображает форму редактирования урока.
func (h *LessonWebHandler) RenderEditLessonForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/lesson-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	course, err := h.courseService.GetCourse(ctx, categoryID, courseID)
	if err != nil {
		return c.Status(404).Render("pages/lesson-form", fiber.Map{
			"title":        "Курс не найден",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"error":        "Курс с указанным ID не найден",
		}, "layouts/main")
	}

	lesson, err := h.lessonService.GetLesson(ctx, lessonID, courseID)
	if err != nil {
		return c.Status(404).Render("pages/lesson-form", fiber.Map{
			"title":        "Урок не найден",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"error":        "Урок с указанным ID не найден",
		}, "layouts/main")
	}

	lessonView := LessonView{
		ID:        lesson.Data.ID,
		CourseID:  lesson.Data.CourseID,
		Title:     lesson.Data.Title,
		Content:   lesson.Data.Content,
		CreatedAt: formatDateTime(lesson.Data.CreatedAt),
		UpdatedAt: formatDateTime(lesson.Data.UpdatedAt),
	}

	return c.Render("pages/lesson-form", fiber.Map{
		"title":        "Редактировать урок",
		"categoryID":   categoryID,
		"categoryName": category.Title,
		"courseID":     courseID,
		"courseName":   course.Data.Title,
		"lesson":       lessonView,
	}, "layouts/main")
}

// CreateLesson обрабатывает создание нового урока из формы.
func (h *LessonWebHandler) CreateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")

	title := c.FormValue("title")
	content := c.FormValue("content")

	log.Printf("[DEBUG] CreateLesson: title=%s, content length=%d", title, len(content))

	if title == "" {
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)

		return c.Status(400).Render("pages/lesson-form", fiber.Map{
			"title":        "Новый урок",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"error":        "Название урока не может быть пустым",
		}, "layouts/main")
	}

	input := request.LessonCreate{
		Title:   title,
		Content: content,
	}

	_, err := h.lessonService.CreateLesson(ctx, courseID, input)
	if err != nil {
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)

		return c.Status(400).Render("pages/lesson-form", fiber.Map{
			"title":        "Новый урок",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"error":        "Ошибка создания урока: " + err.Error(),
		}, "layouts/main")
	}

	return c.Redirect("/admin/categories/" + categoryID + "/courses/" + courseID + "/lessons")
}

// UpdateLesson обрабатывает обновление урока из формы.
func (h *LessonWebHandler) UpdateLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	title := c.FormValue("title")
	content := c.FormValue("content")

	log.Printf("[DEBUG] UpdateLesson: lessonID=%s, title=%s, content length=%d", lessonID, title, len(content))
	log.Printf("[DEBUG] Content first 100 chars: %s", content[:min(100, len(content))])

	if title == "" {
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)
		lesson, _ := h.lessonService.GetLesson(ctx, lessonID, courseID)

		var lessonView *LessonView
		if lesson != nil {
			lessonView = &LessonView{
				ID:        lesson.Data.ID,
				CourseID:  lesson.Data.CourseID,
				Title:     lesson.Data.Title,
				Content:   lesson.Data.Content,
				CreatedAt: formatDateTime(lesson.Data.CreatedAt),
				UpdatedAt: formatDateTime(lesson.Data.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/lesson-form", fiber.Map{
			"title":        "Редактировать урок",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"lesson":       lessonView,
			"error":        "Название урока не может быть пустым",
		}, "layouts/main")
	}

	input := request.LessonUpdate{
		Title:   title,
		Content: content,
	}

	_, err := h.lessonService.UpdateLesson(ctx, lessonID, courseID, input)
	if err != nil {
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		course, _ := h.courseService.GetCourse(ctx, categoryID, courseID)
		lesson, _ := h.lessonService.GetLesson(ctx, lessonID, courseID)

		var lessonView *LessonView
		if lesson != nil {
			lessonView = &LessonView{
				ID:        lesson.Data.ID,
				CourseID:  lesson.Data.CourseID,
				Title:     lesson.Data.Title,
				Content:   lesson.Data.Content,
				CreatedAt: formatDateTime(lesson.Data.CreatedAt),
				UpdatedAt: formatDateTime(lesson.Data.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/lesson-form", fiber.Map{
			"title":        "Редактировать урок",
			"categoryID":   categoryID,
			"categoryName": category.Title,
			"courseID":     courseID,
			"courseName":   course.Data.Title,
			"lesson":       lessonView,
			"error":        "Ошибка обновления урока: " + err.Error(),
		}, "layouts/main")
	}

	return c.Redirect("/admin/categories/" + categoryID + "/courses/" + courseID + "/lessons")
}

// DeleteLesson обрабатывает удаление урока.
func (h *LessonWebHandler) DeleteLesson(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("category_id")
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	err := h.lessonService.DeleteLesson(ctx, lessonID, courseID)
	if err != nil {
		return c.Redirect("/admin/categories/" + categoryID + "/courses/" + courseID + "/lessons")
	}

	return c.Redirect("/admin/categories/" + categoryID + "/courses/" + courseID + "/lessons")
}
