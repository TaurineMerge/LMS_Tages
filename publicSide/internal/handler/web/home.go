package web

import (
	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) RenderHome(c *fiber.Ctx) error {
	// In the future, we would fetch courses from a service here.
	// For now, we pass nil, and the template will show the "empty state".
	return c.Render("pages/home", fiber.Map{
		"title":   "Home",
		"courses": nil,
	}, "layouts/main")
}
