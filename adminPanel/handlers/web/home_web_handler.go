package web

import (
	"adminPanel/handlers/dto/request"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
)

// HomeWebHandler обрабатывает главную страницу админ-панели.
type HomeWebHandler struct {
	categoryService *services.CategoryService
	courseService   *services.CourseService
	lessonService   *services.LessonService
}

// NewHomeWebHandler создает новый обработчик главной страницы.
func NewHomeWebHandler(
	categoryService *services.CategoryService,
	courseService *services.CourseService,
	lessonService *services.LessonService,
) *HomeWebHandler {
	return &HomeWebHandler{
		categoryService: categoryService,
		courseService:   courseService,
		lessonService:   lessonService,
	}
}

// RenderHome отображает главную страницу админ-панели со статистикой.
func (h *HomeWebHandler) RenderHome(c *fiber.Ctx) error {
	ctx := c.UserContext()

	categories, err := h.categoryService.GetCategories(ctx)
	if err != nil {
		return c.Status(500).Render("pages/home", fiber.Map{
			"title": "Главная",
			"error": "Ошибка загрузки статистики",
		}, "layouts/main")
	}

	coursesFilter := request.CourseFilter{}
	coursesResp, err := h.courseService.GetCourses(ctx, coursesFilter)
	if err != nil {
		return c.Status(500).Render("pages/home", fiber.Map{
			"title": "Главная",
			"error": "Ошибка загрузки статистики",
		}, "layouts/main")
	}

	lessonsCount := 0

	return c.Render("pages/home", fiber.Map{
		"title":           "Главная",
		"categoriesCount": len(categories),
		"coursesCount":    len(coursesResp.Data.Items),
		"lessonsCount":    lessonsCount,
	}, "layouts/main")
}
