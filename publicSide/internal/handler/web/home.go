// Package web содержит обработчики для рендеринга веб-страниц.
package web

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/gofiber/fiber/v2"
)

// HomeHandler обрабатывает HTTP-запросы для главной страницы.
type HomeHandler struct {
	categoriesService service.CategoryService
	coursesService    service.CourseService
}

// NewHomeHandler создает новый экземпляр HomeHandler.
func NewHomeHandler(categoriesService service.CategoryService, coursesService service.CourseService) *HomeHandler {
	return &HomeHandler{
		categoriesService: categoriesService,
		coursesService:    coursesService,
	}
}

// RenderHome отображает главную страницу.
// Он загружает несколько непустых категорий и для каждой из них —
// небольшое превью курсов для отображения.
func (h *HomeHandler) RenderHome(c *fiber.Ctx) error {
	const COURSE_LIMIT = 5
	ctx := c.UserContext()

	// Получаем несколько категорий для отображения на главной.
	categoriesDTOs, _, err := h.categoriesService.GetAllNotEmpty(ctx, 1, 5)
	if err != nil {
		slog.Error("Failed to get categories for home page", "error", err)
		categoriesDTOs = []response.CategoryDTO{}
	}

	// Для каждой категории загружаем превью курсов.
	categories := make([]viewmodel.CategoryViewModel, 0, len(categoriesDTOs))
	for _, cat := range categoriesDTOs {
		coursesDTOs, coursesPagination, err := h.coursesService.GetCoursesByCategoryID(ctx, cat.ID, 1, COURSE_LIMIT, "", "")
		if err != nil {
			slog.Error("Failed to get courses for category", "categoryID", cat.ID, "error", err)
			coursesDTOs = []response.CourseDTO{}
			coursesPagination = response.Pagination{}
		}
		categories = append(categories, viewmodel.NewCategoryViewModel(cat, coursesDTOs, coursesPagination, COURSE_LIMIT))
	}

	return c.Render("pages/home", fiber.Map{
		"Header":  viewmodel.NewHeader(),
		"User":    viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":    viewmodel.NewMain("Home"),
		"Context": viewmodel.NewHomePageViewModel(categories),
	}, "layouts/main")
}
