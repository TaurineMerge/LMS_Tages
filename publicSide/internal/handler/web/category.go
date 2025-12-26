// Package web содержит обработчики для рендеринга веб-страниц.
package web

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/request"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// CategoryHandler обрабатывает HTTP-запросы, связанные со страницами категорий.
type CategoryHandler struct {
	categoriesService service.CategoryService
	coursesService    service.CourseService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
func NewCategoryHandler(categoriesService service.CategoryService, coursesService service.CourseService) *CategoryHandler {
	return &CategoryHandler{
		categoriesService: categoriesService,
		coursesService:    coursesService,
	}
}

// RenderCategories отображает страницу со списком категорий.
// Для каждой категории также загружается небольшое превью курсов.
func (h *CategoryHandler) RenderCategories(c *fiber.Ctx) error {
	const COURSE_LIMIT = 5

	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}
	ctx := c.UserContext()

	// Получаем только те категории, в которых есть курсы.
	categoriesDTOs, pagination, err := h.categoriesService.GetAllNotEmpty(ctx, query.Page, query.Limit)
	if err != nil {
		slog.Error("Failed to get categories for home page", "error", err)
		categoriesDTOs = []response.CategoryDTO{}
	}

	// Для каждой категории загружаем превью из нескольких курсов.
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
	vm := viewmodel.NewCategoriesPageViewMode(categories, pagination)

	return c.Render("pages/categories", fiber.Map{
		"Header":  viewmodel.NewHeader(),
		"User":    viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":    viewmodel.NewMain("Categories"),
		"Context": vm,
	}, "layouts/main")
}
