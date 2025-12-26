// Package middleware предоставляет промежуточные обработчики для Fiber.
package middleware

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/gofiber/fiber/v2"
)

// WebErrorHandler является обработчиком ошибок для веб-страниц (не API).
// Он перехватывает ошибки и рендерит HTML-страницу с информацией об ошибке.
func WebErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		slog.Info("Handler error", "error", err)
		appErr = err.(*apperrors.AppError)
	} else if strings.Contains(err.Error(), "Cannot GET") {
		// Обработка стандартной ошибки Fiber для несуществующих маршрутов.
		appErr = apperrors.NewNotFound("Page").(*apperrors.AppError)
	} else {
		// Все остальные ошибки считаются внутренними.
		slog.Error("Unhandled web error", "error", err)
		appErr = apperrors.NewInternal().(*apperrors.AppError)
	}

	return c.Status(appErr.HTTPStatus).Render("pages/error", fiber.Map{
		"Header":     viewmodel.NewHeader(),
		"User":       viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":       viewmodel.NewMain("Home"),
		"Title":      "Error",
		"HTTPStatus": appErr.HTTPStatus,
		"Message":    appErr.Message,
	}, "layouts/main")
}

// NoCache устанавливает заголовки, запрещающие кэширование на стороне клиента.
// Используется в режиме разработки для гарантии загрузки свежих версий ассетов.
func NoCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		return c.Next()
	}
}
