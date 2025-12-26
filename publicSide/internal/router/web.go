// Package router отвечает за настройку и регистрацию маршрутов приложения.
package router

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

// WebRouter инкапсулирует обработчики для всех маршрутов веб-интерфейса.
type WebRouter struct {
	Config              *config.AppConfig
	HomeHandler         *web.HomeHandler
	CategoryPageHandler *web.CategoryHandler
	CoursesHandler      *web.CoursesHandler
	WebLessonHandler    *web.LessonHandler
	AuthHandler         *web.AuthHandler
	AuthMiddleware      *web.AuthMiddleware
}

// Setup настраивает и регистрирует все маршруты для веб-интерфейса.
// Он также применяет глобальные middleware, такие как `NoCache` в режиме разработки,
// раздачу статики и проверку аутентификации пользователя.
func (r *WebRouter) Setup(app *fiber.App) {
	if r.Config.Dev {
		app.Use(middleware.NoCache())
	}

	app.Static("/static", "./static")

	// Middleware для извлечения информации о пользователе из cookie.
	app.Use(r.AuthMiddleware.WithUser)

	// Маршруты аутентификации
	app.Get("/login", r.AuthHandler.Login)
	app.Get("/logout", r.AuthHandler.Logout)
	app.Get("/auth/callback", r.AuthHandler.Callback)

	// Основные маршруты веб-приложения
	app.Get(routing.RouteHome, r.HomeHandler.RenderHome)
	app.Get(routing.RouteCategories, r.CategoryPageHandler.RenderCategories)
	app.Get(routing.RouteCourses, r.CoursesHandler.RenderCourses)
	app.Get(routing.RouteCourse, r.CoursesHandler.RenderCoursePage)
	app.Get(routing.RouteLesson, r.WebLessonHandler.RenderLesson)
}
