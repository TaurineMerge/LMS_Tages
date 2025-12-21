package web

import (
	"context"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/breadcrumbs"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonService service.LessonService
	courseService service.CourseService
}

func NewLessonHandler(ls service.LessonService, cs service.CourseService) *LessonHandler {
	return &LessonHandler{
		lessonService: ls,
		courseService: cs,
	}
}

// RenderLesson отображает страницу урока.
// Вместо fiber.Map теперь используется строго типизированная ViewModel.
func (h *LessonHandler) RenderLesson(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	courseID := c.Params(routing.PathVariableCourseID)
	lessonID := c.Params(routing.PathVariableLessonID)

	// Теперь вся логика сбора данных инкапсулирована в функции buildLessonPageViewModel
	vm, err := h.buildLessonPageViewModel(c.UserContext(), categoryID, courseID, lessonID)
	if err != nil {
		// Обработка ошибок, которые вернулись из buildLessonPageViewModel
		// Например, NotFound ошибки будут обработаны глобальным ErrorHandler
		slog.Error("Failed to build lesson page view model", "lessonID", lessonID, "error", err)
		return err
	}

	// Передаем ViewModel напрямую в шаблон
	return c.Render("pages/lesson", vm, "layouts/main")
}

// buildLessonPageViewModel собирает все необходимые данные для LessonPageViewModel.
// Она вызывает соответствующие сервисы и преобразует их результаты в единую ViewModel.
func (h *LessonHandler) buildLessonPageViewModel(ctx context.Context, categoryID, courseID, lessonID string) (viewmodel.LessonPageViewModel, error) {
	// Инициализация пустой ViewModel
	vm := viewmodel.LessonPageViewModel{}

	// 1. Получаем детальную информацию об уроке
	lesson, err := h.lessonService.GetByID(ctx, categoryID, courseID, lessonID)
	if err != nil {
		return viewmodel.LessonPageViewModel{}, apperrors.NewNotFound("Lesson")
	}
	vm.Lesson = lesson

	// 2. Получаем информацию о курсе
	course, err := h.courseService.GetCourseByID(ctx, categoryID, courseID)
	if err != nil {
		return viewmodel.LessonPageViewModel{}, apperrors.NewNotFound("Course")
	}
	vm.Course = course

	// 3. Получаем информацию о категории
	category, err := h.courseService.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return viewmodel.LessonPageViewModel{}, apperrors.NewNotFound("Category")
	}
	vm.Category = *category // Разыменовываем указатель

	// 4. Собираем хедер страницы с хлебными крошками
	vm.PageHeader = viewmodel.PageHeaderViewModel{
		Title:       lesson.Title,
		Breadcrumbs: breadcrumbs.ForLessonPage(*category, course, lesson),
	}

	// 5. Получаем соседние уроки
	prevLesson, nextLesson, err := h.lessonService.GetNeighboringLessons(ctx, categoryID, courseID, lessonID)
	if err != nil {
		// Логируем ошибку, но не прерываем, так как соседних уроков может и не быть
		slog.Warn("Failed to get neighboring lessons", "lessonID", lessonID, "error", err)
		vm.PrevLesson = response.LessonDTO{}
		vm.NextLesson = response.LessonDTO{}
	} else {
		vm.PrevLesson = prevLesson
		vm.NextLesson = nextLesson
	}

	// 6. Получаем все уроки для сайдбара
	// Используем большое значение limit, чтобы получить все уроки одним запросом.
	allLessons, _, err := h.lessonService.GetAllByCourseID(ctx, categoryID, courseID, 1, 999, "created_at")
	if err != nil {
		slog.Warn("Failed to get all lessons for sidebar", "courseID", courseID, "error", err)
		vm.AllLessons = []response.LessonDTO{} // Пустой список, если ошибка
	} else {
		vm.AllLessons = allLessons
	}

	return vm, nil
}
