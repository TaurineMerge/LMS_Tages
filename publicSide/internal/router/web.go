package router

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/middleware"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
)

// WebRouter инкапсулирует зависимости и логику для регистрации веб-маршрутов.
type WebRouter struct {
	Config              *config.AppConfig
	HomeHandler         *web.HomeHandler
	CategoryPageHandler *web.CategoryHandler
	CoursesHandler      *web.CoursesHandler
	WebLessonHandler    *web.LessonHandler
	AuthHandler         *web.AuthHandler
	AuthMiddleware      *web.AuthMiddleware
}

// Setup регистрирует все маршруты для веб-интерфейса.
func (r *WebRouter) Setup(webRouter *fiber.App) {
	// Middleware для no-cache в режиме разработки.
	if r.Config.Dev {
		webRouter.Use(middleware.NoCache())
	}

	webRouter.Static("/static", "./static")

	webRouter.Use(r.AuthMiddleware.WithUser)

	// Auth
	webRouter.Get("/login", r.AuthHandler.Login)
	webRouter.Get("/logout", r.AuthHandler.Logout)
	webRouter.Get("/auth/callback", r.AuthHandler.Callback)

	// Регистрация маршрутов с использованием полных путей из пакета routing
	webRouter.Get(routing.RouteHome, r.HomeHandler.RenderHome)
	webRouter.Get(routing.RouteCategories, r.CategoryPageHandler.RenderCategories)
	webRouter.Get(routing.RouteCourses, r.CoursesHandler.RenderCourses)
	webRouter.Get(routing.RouteCourse, r.CoursesHandler.RenderCoursePage)
	webRouter.Get(routing.RouteLesson, r.WebLessonHandler.RenderLesson)
}
