// Package router отвечает за настройку и регистрацию маршрутов приложения.
package router

import (
	v1 "github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// APIRouter инкапсулирует обработчики для всех маршрутов API.
type APIRouter struct {
	APICategoryHandler *v1.CategoryHandler
	APICourseHandler   *v1.CourseHandler
	APILessonHandler   *v1.LessonHandler
}

// Setup настраивает и регистрирует все маршруты API v1.
// Он также настраивает маршрут для отображения документации Swagger.
func (r *APIRouter) Setup(app *fiber.App) {
	// Раздача статического файла swagger.json
	app.Static("/doc", "./doc/swagger")

	apiV1 := app.Group(routing.RouteAPIV1)

	// Настройка Swagger UI
	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	// Маршруты для категорий
	apiV1.Get(routing.RouteCategories, r.APICategoryHandler.GetAllCategories)
	apiV1.Get(routing.RouteCategory, r.APICategoryHandler.GetCategoryByID)

	// Маршруты для курсов
	apiV1.Get(routing.RouteCourses, r.APICourseHandler.GetCoursesByCategoryID)
	apiV1.Get(routing.RouteCourse, r.APICourseHandler.GetCourseByID)

	// Маршруты для уроков
	apiV1.Get(routing.RouteLessons, r.APILessonHandler.GetLessonsByCourseID)
	apiV1.Get(routing.RouteLesson, r.APILessonHandler.GetLessonByID)
}
