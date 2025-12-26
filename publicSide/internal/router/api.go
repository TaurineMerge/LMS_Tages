package router

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// APIRouter инкапсулирует зависимости и логику для регистрации API-маршрутов.
type APIRouter struct {
	APICategoryHandler *v1.CategoryHandler
	APICourseHandler   *v1.CourseHandler
	APILessonHandler   *v1.LessonHandler
}

// Setup регистрирует все маршруты для API в явном, централизованном виде.
func (r *APIRouter) Setup(app *fiber.App) {
	app.Static("/doc", "./doc/swagger")

	// /api/v1
	apiV1 := app.Group(routing.RouteAPIV1)

	// /api/v1/swagger/*
	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	apiV1.Get(routing.RouteCategories, r.APICategoryHandler.GetAllCategories)
	apiV1.Get(routing.RouteCategory, r.APICategoryHandler.GetCategoryByID)

	apiV1.Get(routing.RouteCourses, r.APICourseHandler.GetCoursesByCategoryID)
	apiV1.Get(routing.RouteCourse, r.APICourseHandler.GetCourseByID)
	
	apiV1.Get(routing.RouteLessons, r.APILessonHandler.GetLessonsByCourseID)
	apiV1.Get(routing.RouteLesson, r.APILessonHandler.GetLessonByID)
}
