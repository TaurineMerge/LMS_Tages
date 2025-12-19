package web

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler/api/v1/dto/response"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apiconst"
	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonService service.LessonService
	courseService service.CourseService
}

func NewLessonHandler(ls service.LessonService, cs service.CourseService) *LessonHandler {
	return &LessonHandler{
		lessonService: ls,
		courseService: cs,
	}
}

func (h *LessonHandler) RenderLesson(c *fiber.Ctx) error {
	categoryID := c.Params(apiconst.PathVariableCategoryID)
	courseID := c.Params(apiconst.PathVariableCourseID)
	lessonID := c.Params(apiconst.PathVariableLessonID)

	// In a real app, you'd want to handle these errors properly
	lesson, err := h.lessonService.GetByID(c.UserContext(), categoryID, courseID, lessonID)
	if err != nil {
		return err // Let the global error handler deal with it
	}

	course, err := h.courseService.GetCourseByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return err
	}

	prevLesson, nextLesson, err := h.lessonService.GetNeighboringLessons(c.UserContext(), categoryID, courseID, lessonID)
	if err != nil {
		// Log the error but don't fail the page render if neighbors can't be fetched
		prevLesson = response.LessonDTO{}
		nextLesson = response.LessonDTO{}
	}

	// Fetch all lessons for the sidebar
	allLessons, _, err := h.lessonService.GetAllByCourseID(c.UserContext(), categoryID, courseID, 1, 999, "") // page 1, limit 999, default sort
	if err != nil {
		// Don't fail the render, just show an empty sidebar
		allLessons = []response.LessonDTO{}
	}

	return c.Render("pages/lesson", fiber.Map{
		"title":       lesson.Title,
		"lesson":      lesson,
		"course":      course,
		"prevLesson":  prevLesson,
		"nextLesson":  nextLesson,
		"allLessons":  allLessons,
		"categoryID": categoryID,
	}, "layouts/main")
}
