package web

import (
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LessonHandler struct {
	lessonsService   service.LessonService
	coursesService   service.CourseService
	categoriesService service.CategoryService
}

func NewLessonHandler(ls service.LessonService, cs service.CourseService, cats service.CategoryService) *LessonHandler {
	return &LessonHandler{
		lessonsService:   ls,
		coursesService:   cs,
		categoriesService: cats,
	}
}

// RenderLesson отображает страницу урока.
func (h *LessonHandler) RenderLesson(c *fiber.Ctx) error {
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
	ctx := c.UserContext()

	lessonDTODetailed, err := h.lessonsService.GetByID(ctx, categoryID, courseID, lessonID)
	if err != nil {
		slog.Error("Failed to get lesson by ID", "lessonID", lessonID, "error", err)
		return err
	}

	prevLessonDTO, nextLessonDTO, err := h.lessonsService.GetNeighboringLessons(ctx, categoryID, courseID, lessonID)
	if err != nil {
		slog.Error("Failed to get neighboring lessons", "lessonID", lessonID, "error", err)
		prevLessonDTO = response.LessonDTO{}
		nextLessonDTO = response.LessonDTO{}
	}

	lessonsDTOs, _, err := h.lessonsService.GetAllByCourseID(ctx, categoryID, courseID, 1, 100, "")
	if err != nil {
		slog.Error("Failed to get all lessons", "lessonID", lessonID, "error", err)
		return err
	}

	courseDTO, err := h.coursesService.GetCourseByID(ctx, categoryID, courseID)
	if err != nil {
		slog.Error("Failed to get course by ID", "courseID", courseID, "error", err)
		return err
	}

	categoryDTO, err := h.categoriesService.GetByID(ctx, categoryID)
	if err != nil {
		slog.Error("Failed to get category by ID", "categoryID", categoryID, "error", err)
		return err
	}

	return c.Render("pages/lesson", fiber.Map{
		"Header": viewmodel.NewHeader(),
		"User":   viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":   viewmodel.NewMain("Lesson"),
		"Context": viewmodel.NewLessonPageViewModel(
			lessonDTODetailed,
			courseDTO,
			categoryDTO,
			nextLessonDTO,
			prevLessonDTO,
			lessonsDTOs,
		),
	}, "layouts/main")
}
