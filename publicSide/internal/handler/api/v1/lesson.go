// Package v1 содержит обработчики для API версии 1.
package v1

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/request"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// LessonHandler обрабатывает HTTP-запросы, связанные с уроками.
type LessonHandler struct {
	service service.LessonService
}

// NewLessonHandler создает новый экземпляр LessonHandler.
func NewLessonHandler(s service.LessonService) *LessonHandler {
	return &LessonHandler{service: s}
}

// GetLessonsByCourseID обрабатывает запрос на получение списка уроков для конкретного курса.
// @Summary Получить список уроков курса
// @Description Получает страницы списка уроков для указанного курса. Поддерживает пагинацию и сортировку.
// @Tags Lessons
// @Accept json
// @Produce json
// @Param category_id path string true "Уникальный идентификатор категории"
// @Param course_id path string true "Уникальный идентификатор курса"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Param sort query string false "Поле и порядок сортировки (например, -created_at)"
// @Success 200 {object} response.SuccessResponse{data=response.PaginatedLessonsData} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверные параметры запроса"
// @Failure 404 {object} response.ErrorResponse "Категория или курс не найдены"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{category_id}/courses/{course_id}/lessons [get]
func (h *LessonHandler) GetLessonsByCourseID(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}
	courseID := c.Params(routing.PathVariableCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCourseID)
	}

	var query request.ListQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	lessons, pagination, err := h.service.GetAllByCourseID(c.UserContext(), categoryID, courseID, query.Page, query.Limit, query.Sort)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedLessonsData{
			Items:      lessons,
			Pagination: pagination,
		},
	})
}

// GetLessonByID обрабатывает запрос на получение одного урока по его ID.
// @Summary Получить урок по ID
// @Description Получает детали одного урока по его UUID в рамках курса и категории.
// @Tags Lessons
// @Accept json
// @Produce json
// @Param category_id path string true "Уникальный идентификатор категории"
// @Param course_id path string true "Уникальный идентификатор курса"
// @Param lesson_id path string true "Уникальный идентификатор урока"
// @Success 200 {object} response.SuccessResponse{data=response.LessonDTODetailed} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} response.ErrorResponse "Категория, курс или урок не найдены"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{category_id}/courses/{course_id}/lessons/{lesson_id} [get]
func (h *LessonHandler) GetLessonByID(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}
	courseID := c.Params(routing.PathVariableCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCourseID)
	}

	lessonID := c.Params(routing.PathVariableLessonID)
	if _, err := uuid.Parse(lessonID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableLessonID)
	}

	lesson, err := h.service.GetByID(c.UserContext(), categoryID, courseID, lessonID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   lesson,
	})
}
