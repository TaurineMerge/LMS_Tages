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

// CategoryHandler обрабатывает HTTP-запросы, связанные с категориями.
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler создает новый экземпляр CategoryHandler.
func NewCategoryHandler(s service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

// GetAllCategories обрабатывает запрос на получение списка всех категорий с пагинацией.
// @Summary Получить список всех категорий
// @Description Получает страницы списка всех категорий.
// @Tags Categories
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20)
// @Success 200 {object} response.SuccessResponse{data=response.PaginatedCategoriesData} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверные параметры запроса"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	var query request.PaginationQuery
	if err := c.QueryParser(&query); err != nil {
		return apperrors.NewInvalidRequest("Wrong query parameters")
	}

	categories, pagination, err := h.service.GetAll(c.UserContext(), query.Page, query.Limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data: response.PaginatedCategoriesData{
			Items:      categories,
			Pagination: pagination,
		},
	})
}

// GetCategoryByID обрабатывает запрос на получение одной категории по ее ID.
// @Summary Получить категорию по ID
// @Description Получает детали одной категории по ее UUID.
// @Tags Categories
// @Accept json
// @Produce json
// @Param category_id path string true "Уникальный идентификатор категории"
// @Success 200 {object} response.SuccessResponse{data=response.CategoryDTO} "Успешный ответ"
// @Failure 400 {object} response.ErrorResponse "Неверный формат ID"
// @Failure 404 {object} response.ErrorResponse "Категория не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /categories/{category_id} [get]
func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}

	category, err := h.service.GetByID(c.UserContext(), categoryID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status: response.StatusSuccess,
		Data:   category,
	})
}
