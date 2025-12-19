package web

import (
	"adminPanel/handlers/dto/request"
	"adminPanel/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CategoryView представляет категорию для отображения
type CategoryView struct {
	ID        string
	Title     string
	CreatedAt string
	UpdatedAt string
}

// CategoryWebHandler обрабатывает веб-страницы для управления категориями
type CategoryWebHandler struct {
	categoryService *services.CategoryService
}

// NewCategoryWebHandler создает новый обработчик веб-страниц категорий
func NewCategoryWebHandler(categoryService *services.CategoryService) *CategoryWebHandler {
	return &CategoryWebHandler{
		categoryService: categoryService,
	}
}

// RenderCategoriesEditor отображает страницу со списком категорий
func (h *CategoryWebHandler) RenderCategoriesEditor(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Получаем все категории
	categories, err := h.categoryService.GetCategories(ctx)
	if err != nil {
		return c.Status(500).Render("pages/categories-editor", fiber.Map{
			"title": "Редактор категорий",
			"error": "Ошибка загрузки категорий",
		}, "layouts/main")
	}

	// Преобразуем в CategoryView
	categoryViews := make([]CategoryView, 0, len(categories))
	for _, cat := range categories {
		categoryViews = append(categoryViews, CategoryView{
			ID:        cat.ID,
			Title:     cat.Title,
			CreatedAt: formatDateTime(cat.CreatedAt),
			UpdatedAt: formatDateTime(cat.UpdatedAt),
		})
	}

	return c.Render("pages/categories-editor", fiber.Map{
		"title":      "Редактор категорий",
		"categories": categoryViews,
	}, "layouts/main")
}

// RenderNewCategoryForm отображает форму создания новой категории
func (h *CategoryWebHandler) RenderNewCategoryForm(c *fiber.Ctx) error {
	return c.Render("pages/category-form", fiber.Map{
		"title": "Новая категория",
	}, "layouts/main")
}

// RenderEditCategoryForm отображает форму редактирования категории
func (h *CategoryWebHandler) RenderEditCategoryForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("id")

	// Получаем категорию
	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		return c.Status(404).Render("pages/category-form", fiber.Map{
			"title": "Категория не найдена",
			"error": "Категория с указанным ID не найдена",
		}, "layouts/main")
	}

	categoryView := CategoryView{
		ID:        category.ID,
		Title:     category.Title,
		CreatedAt: formatDateTime(category.CreatedAt),
		UpdatedAt: formatDateTime(category.UpdatedAt),
	}

	return c.Render("pages/category-form", fiber.Map{
		"title":    "Редактировать категорию",
		"category": categoryView,
	}, "layouts/main")
}

// CreateCategory обрабатывает создание новой категории
func (h *CategoryWebHandler) CreateCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()

	title := c.FormValue("title")
	if title == "" {
		return c.Status(400).Render("pages/category-form", fiber.Map{
			"title": "Новая категория",
			"error": "Название категории не может быть пустым",
		}, "layouts/main")
	}

	// Создаем категорию
	input := request.CategoryCreate{
		Title: title,
	}

	_, err := h.categoryService.CreateCategory(ctx, input)
	if err != nil {
		return c.Status(400).Render("pages/category-form", fiber.Map{
			"title": "Новая категория",
			"error": "Ошибка создания категории: " + err.Error(),
		}, "layouts/main")
	}

	// Перенаправляем на список категорий (nginx проксирует /admin -> backend root)
	return c.Redirect("/admin/categories")
}

// UpdateCategory обрабатывает обновление категории
func (h *CategoryWebHandler) UpdateCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("id")

	title := c.FormValue("title")
	if title == "" {
		// Получаем категорию для отображения в форме
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		var categoryView *CategoryView
		if category != nil {
			categoryView = &CategoryView{
				ID:        category.ID,
				Title:     category.Title,
				CreatedAt: formatDateTime(category.CreatedAt),
				UpdatedAt: formatDateTime(category.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/category-form", fiber.Map{
			"title":    "Редактировать категорию",
			"category": categoryView,
			"error":    "Название категории не может быть пустым",
		}, "layouts/main")
	}

	// Обновляем категорию
	input := request.CategoryUpdate{
		Title: title,
	}

	_, err := h.categoryService.UpdateCategory(ctx, categoryID, input)
	if err != nil {
		// Получаем категорию для отображения в форме
		category, _ := h.categoryService.GetCategory(ctx, categoryID)
		var categoryView *CategoryView
		if category != nil {
			categoryView = &CategoryView{
				ID:        category.ID,
				Title:     category.Title,
				CreatedAt: formatDateTime(category.CreatedAt),
				UpdatedAt: formatDateTime(category.UpdatedAt),
			}
		}

		return c.Status(400).Render("pages/category-form", fiber.Map{
			"title":    "Редактировать категорию",
			"category": categoryView,
			"error":    "Ошибка обновления категории: " + err.Error(),
		}, "layouts/main")
	}

	// Перенаправляем на список категорий (nginx проксирует /admin -> backend root)
	return c.Redirect("/admin/categories")
}

// DeleteCategory обрабатывает удаление категории
func (h *CategoryWebHandler) DeleteCategory(c *fiber.Ctx) error {
	ctx := c.UserContext()
	categoryID := c.Params("id")

	err := h.categoryService.DeleteCategory(ctx, categoryID)
	if err != nil {
		// Можно добавить flash-сообщение об ошибке
		return c.Redirect("/admin/categories")
	}

	// Перенаправляем на список категорий
	return c.Redirect("/admin/categories")
}

// formatDateTime форматирует время для отображения
func formatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}
