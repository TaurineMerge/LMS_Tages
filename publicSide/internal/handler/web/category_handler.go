package web

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/gofiber/fiber/v2"
)

// categoryCourseView описывает краткую информацию о курсе для страницы категорий.
type categoryCourseView struct {
	ID    string
	Title string
}

// categoryView описывает категорию с ограниченным списком курсов для отображения на странице.
type categoryView struct {
	ID             string
	Title          string
	TotalCourses   int
	Courses        []categoryCourseView
	HasMoreCourses bool
}

// CategoryHandler - обработчик для страниц категорий.
// Он использует сервисы категорий и курсов, чтобы получать данные из БД knowledge_base.
type CategoryHandler struct {
	categoryService service.CategoryService
	courseService   service.CourseService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
// Сервисы передаются из main.go и внутри уже используют репозитории и подключение к БД.
func NewCategoryHandler(categoryService service.CategoryService, courseService service.CourseService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		courseService:   courseService,
	}
}

// RenderCategories отображает страницу категорий.
// Данные берутся из БД knowledge_base через сервисы категорий и курсов.
func (h *CategoryHandler) RenderCategories(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// Получаем список всех категорий (ограничим 100, чтобы не перегружать страницу).
	const categoriesPage = 1
	const categoriesLimit = 100

	categoryDTOs, _, err := h.categoryService.GetAll(ctx, categoriesPage, categoriesLimit)
	if err != nil {
		return err
	}

	categories := make([]categoryView, 0, len(categoryDTOs))

	// Для каждой категории подтягиваем до 10 публичных курсов.
	const coursesPage = 1
	const coursesLimit = 10

	for _, cat := range categoryDTOs {
		catView := categoryView{
			ID:    cat.ID,
			Title: cat.Title,
		}

		courses, pagination, err := h.courseService.GetCoursesByCategoryID(ctx, cat.ID, coursesPage, coursesLimit)
		if err != nil {
			// Если для конкретной категории произошла ошибка, прерываем рендер с ней —
			// это позволит увидеть проблему в логах и в Swagger.
			return err
		}

		catView.TotalCourses = pagination.Total

		courseViews := make([]categoryCourseView, 0, len(courses))
		for _, course := range courses {
			courseViews = append(courseViews, categoryCourseView{
				ID:    course.ID,
				Title: course.Title,
			})
		}
		catView.Courses = courseViews

		// Показываем кнопку "Показать ещё", если всего курсов больше, чем выведено на странице.
		catView.HasMoreCourses = pagination.Total > len(courseViews)

		categories = append(categories, catView)
	}

	// Если курсов ещё нет (текущая версия проекта), категории всё равно отобразятся,
	// просто без списка курсов и без кнопки "Показать ещё".
	return c.Render("pages/categories", fiber.Map{
		"title":      "Категории курсов - LMS Tages",
		"categories": categories,
	}, "layouts/main")
}

// RegisterHelpers зарезервирован для регистрации вспомогательных функций шаблонов, если они понадобятся.
func (h *CategoryHandler) RegisterHelpers(engine fiber.Views) {
	// Пока дополнительных helper-ов для шаблонов не требуется.
}
