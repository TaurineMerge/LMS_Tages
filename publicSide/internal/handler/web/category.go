package web

import (
	"context"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/breadcrumbs"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/gofiber/fiber/v2"
)

// CategoryHandler - обработчик для страниц категорий.
type CategoryHandler struct {
	categoryService service.CategoryService
	courseService   service.CourseService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
func NewCategoryHandler(categoryService service.CategoryService, courseService service.CourseService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		courseService:   courseService,
	}
}



// RenderCategories отображает страницу категорий.
func (h *CategoryHandler) RenderCategories(c *fiber.Ctx) error {
	vm, err := h.buildCategoriesPageViewModel(c.UserContext())
	if err != nil {
		return err
	}

	return c.Render("pages/categories", vm, "layouts/main")
}

func (h *CategoryHandler) buildCategoriesPageViewModel(ctx context.Context) (viewmodel.CategoriesPageViewModel, error) {
	vm := viewmodel.CategoriesPageViewModel{
		PageHeader: viewmodel.PageHeaderViewModel{
			Title:       "Категории курсов",
			Breadcrumbs: breadcrumbs.ForCategoriesPage(),
		},
	}

	// Получаем список всех категорий (ограничим 100, чтобы не перегружать страницу).
	const categoriesPage = 1
	const categoriesLimit = 100
	categoryDTOs, _, err := h.categoryService.GetAll(ctx, categoriesPage, categoriesLimit)
	if err != nil {
		return viewmodel.CategoriesPageViewModel{}, err
	}

	categories := make([]viewmodel.CategoryView, 0, len(categoryDTOs))

	// Для каждой категории подтягиваем до 10 публичных курсов.
	const coursesPage = 1
	const coursesLimit = 10
	for _, cat := range categoryDTOs {
		catView := viewmodel.CategoryView{
			ID:    cat.ID,
			Title: cat.Title,
		}

		courses, pagination, err := h.courseService.GetCoursesByCategoryID(ctx, cat.ID, coursesPage, coursesLimit, "", "")
		if err != nil {
			return viewmodel.CategoriesPageViewModel{}, err
		}

		catView.TotalCourses = pagination.Total

		courseViews := make([]viewmodel.CategoryCourseView, 0, len(courses))
		for _, course := range courses {
			courseViews = append(courseViews, viewmodel.CategoryCourseView{
				ID:    course.ID,
				Title: course.Title,
			})
		}
		catView.Courses = courseViews
		catView.HasMoreCourses = pagination.Total > len(courseViews)

		categories = append(categories, catView)
	}
	vm.Categories = categories

	return vm, nil
}
