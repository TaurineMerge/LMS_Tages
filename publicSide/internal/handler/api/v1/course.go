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

// CourseHandler обрабатывает HTTP-запросы, связанные с курсами.
type CourseHandler struct {
	courseService service.CourseService
}

// NewCourseHandler создает новый экземпляр CourseHandler.
func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// GetCoursesByCategoryID обрабатывает запрос на получение списка курсов для конкретной категории.
// @Summary Получить список курсов категории
// @Description Получает страницы списка курсов для указанной категории.
// @Tags Courses
// @Accept json
// @Produce json
// @Param category_id path string true "Уникальный идентификатор категории"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} response.SuccessResponse{data=response.PaginatedCoursesData} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверные параметры запроса"
// @Failure 404 {object} response.ErrorResponse "Категория не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{category_id}/courses [get]
func (h *CourseHandler) GetCoursesByCategoryID(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}

	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	// В API не используются фильтры по уровню и сортировка, передаем пустые строки.
	courses, pagination, err := h.courseService.GetCoursesByCategoryID(c.UserContext(), categoryID, query.Page, query.Limit, "", "")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCoursesData{
			Items:      courses,
			Pagination: pagination,
		},
	})
}

// GetCourseByID обрабатывает запрос на получение одного курса по его ID.
// @Summary Получить курс по ID
// @Description Получает детали одного курса по его UUID в рамках категории.
// @Tags Courses
// @Accept json
// @Produce json
// @Param category_id path string true "Уникальный идентификатор категории"
// @Param course_id path string true "Уникальный идентификатор курса"
// @Success 200 {object} response.SuccessResponse{data=response.CourseDTO} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} response.ErrorResponse "Категория или курс не найдены"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{category_id}/courses/{course_id} [get]
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}

	courseID := c.Params(routing.PathVariableCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCourseID)
	}

	course, err := h.courseService.GetCourseByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   course,
	})
}
