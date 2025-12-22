package web

import (
	"fmt"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/request"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// CategoryHandler - обработчик для страниц категорий.
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

// RenderCategories отображает страницу категорий.
func (h *CategoryHandler) RenderCategories(c *fiber.Ctx) error {
	const (
		PAGE         = 1
		LIMIT        = 10
		COURSE_LIMIT = 5
	)
	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}
	ctx := c.UserContext()

	categoriesDTOs, pagination, err := h.categoriesService.GetAllNotEmpty(ctx, PAGE, LIMIT)
	if err != nil {
		slog.Error("Failed to get categories for home page", "error", err)
		categoriesDTOs = []response.CategoryDTO{}
	}
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
	if len(vm.Categories) == 0 {
		return apperrors.NewNotFound(fmt.Sprintf("Page %d", query.Page))
	}
	return c.Render("pages/categories", fiber.Map{
		"Header":  viewmodel.NewHeader(),
		"User":    viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":    viewmodel.NewMain("Categories"),
		"Context": vm,
	}, "layouts/main")
}
