package web

import (
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/clients/testing"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CoursesHandler handles web requests for the courses page.
type CoursesHandler struct {
	courseService   service.CourseService
	categoryService service.CategoryService
	lessonService   service.LessonService
	testService     service.TestService
	testingConfig   config.TestingServiceConfig
}

// NewCoursesHandler creates a new instance of CoursesHandler.
func NewCoursesHandler(
	courseService service.CourseService,
	categoryService service.CategoryService,
	lessonService service.LessonService,
	testService service.TestService,
	testingConfig config.TestingServiceConfig,
) *CoursesHandler {
	return &CoursesHandler{
		courseService:   courseService,
		categoryService: categoryService,
		lessonService:   lessonService,
		testService:     testService,
		testingConfig:   testingConfig,
	}
}

// RenderCourses renders the courses page with filters and sorting.
func (h *CoursesHandler) RenderCourses(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return apperrors.NewInvalidUUID(routing.PathVariableCategoryID)
	}
	page := c.QueryInt("page", 1)
	level := c.Query("level", "all")
	sortBy := c.Query("sort_by", "updated_at")
	limit := c.QueryInt("limit", 28)

	categoryDTO, err := h.categoryService.GetByID(c.UserContext(), categoryID)
	if err != nil {
		return err
	}

	coursesDTOs, coursesPagination, err := h.courseService.GetCoursesByCategoryID(
		c.UserContext(), categoryID, page, limit, level, sortBy,
	)
	if err != nil {
		return err
	}

	lessonAmounts := make([]int, 0, len(coursesDTOs))
	for _, course := range coursesDTOs {
		_, pag, err := h.lessonService.GetAllByCourseID(c.UserContext(), categoryID, course.ID, 1, 1, "")
		if err != nil {
			return err
		}
		lessonAmounts = append(lessonAmounts, pag.Total)
	}

	vm := viewmodel.NewCoursesPageViewModel(
		categoryDTO,
		coursesDTOs,
		coursesPagination,
		lessonAmounts,
		level,
		sortBy,
	)

	vm.Courses = russifyCoursesLevel(vm.Courses)

	return c.Render("pages/courses", fiber.Map{
		"Header":  viewmodel.NewHeader(),
		"User":    viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":    viewmodel.NewMain("Courses"),
		"Context": vm,
	}, "layouts/main")
}

// RenderCoursePage renders the individual course page with course details.
func (h *CoursesHandler) RenderCoursePage(c *fiber.Ctx) error {
	categoryID := c.Params(routing.PathVariableCategoryID)
	if _, err := uuid.Parse(categoryID); err != nil {
		return err
	}
	courseID := c.Params(routing.PathVariableCourseID)
	if _, err := uuid.Parse(courseID); err != nil {
		return err
	}

	categoryDTO, err := h.categoryService.GetByID(c.UserContext(), categoryID)
	if err != nil {
		return err
	}
	courseDTO, err := h.courseService.GetCourseByID(c.UserContext(), categoryID, courseID)
	if err != nil {
		return err
	}
	lessonsDTOs, _, err := h.lessonService.GetAllByCourseID(c.UserContext(), categoryID, courseID, 1, 10, "created_at")
	if err != nil {
		slog.Error("Failed to get first 10 lessons for course page", "courseId", courseID, "error", err)
		return err
	}

	// Fetch test details
	var testVM *viewmodel.TestViewModel
	var testIsNotFound, testServiceIsUnavailable bool

	testData, err := h.testService.GetTest(c.UserContext(), categoryID, courseID)
	if err != nil {
		var appErr *apperrors.AppError
		var unavailableErr *apperrors.ServiceUnavailableError

		if errors.As(err, &appErr) && appErr.HTTPStatus == 404 {
			testIsNotFound = true
		} else if errors.As(err, &unavailableErr) {
			testServiceIsUnavailable = true
		} else {
			slog.Error("Unexpected error fetching test details", "error", err, "courseID", courseID)
			return err
		}
	} else {
		testVM = viewmodel.NewTestViewModel(testData, testing.GetUITestURL(h.testingConfig.BaseURL, categoryID, courseID))
	}

	vm := viewmodel.NewCoursePageViewModel(
		categoryDTO,
		courseDTO,
		lessonsDTOs,
		testVM,
		testIsNotFound,
		testServiceIsUnavailable,
	)

	russifyCourseDetailLevel(vm.Course)

	return c.Render("pages/course", fiber.Map{
		"Header":  viewmodel.NewHeader(),
		"User":    viewmodel.NewUserViewModel(c.Locals(domain.UserContextKey).(domain.UserClaims)),
		"Main":    viewmodel.NewMain("Course"),
		"Context": vm,
	}, "layouts/main")
}

func getLevelRussification(level string) string {
	enLvlToRu := map[string]string{
		"all":    "Все уровни",
		"easy":   "Легкий",
		"medium": "Средний",
		"hard":   "Сложный",
	}

	if ru, ok := enLvlToRu[level]; ok {
		return ru
	}
	return level
}

func russifyCoursesLevel(courses []viewmodel.CourseViewModel) []viewmodel.CourseViewModel {
	for i := range courses {
		courses[i].LevelRu = getLevelRussification(courses[i].Level)
	}
	return courses
}

func russifyCourseDetailLevel(course *viewmodel.CourseDetailViewModel) {
	course.LevelRu = getLevelRussification(course.Level)
}