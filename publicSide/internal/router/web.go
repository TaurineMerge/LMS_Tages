package router

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/web"
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
}

// Setup регистрирует все маршруты для веб-интерфейса.
func (r *WebRouter) Setup(app *fiber.App) {
	// Middleware для no-cache в режиме разработки.
	if r.Config.Dev {
		app.Use(func(c *fiber.Ctx) error {
			c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return c.Next()
		})
	}

	// Статические файлы
	app.Static("/doc", "./doc/swagger")
	app.Static("/static", "./static")

	// Регистрация маршрутов с использованием полных путей из пакета routing
	app.Get(routing.RouteHome, r.HomeHandler.RenderHome)
	app.Get(routing.RouteCategories, r.CategoryPageHandler.RenderCategories)
	app.Get(routing.RouteCourses, r.CoursesHandler.RenderCourses)
	app.Get(routing.RouteCourse, r.CoursesHandler.RenderCoursePage)
	app.Get(routing.RouteLesson, r.WebLessonHandler.RenderLesson)
}
