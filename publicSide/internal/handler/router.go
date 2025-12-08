package handler

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(
	api fiber.Router,
	categoryHandler *CategoryHandler,
	courseHandler *CourseHandler,
	lessonHandler *LessonHandler,
) {
	// Categories
	categories := api.Group("/categories")
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/:category_id", categoryHandler.GetCategoryByID)
	categories.Get("/:category_id/courses", courseHandler.GetCoursesByCategoryID)

	// Courses
	courses := api.Group("/courses")
	courses.Get("/", courseHandler.GetAllCourses)
	courses.Get("/:course_id", courseHandler.GetCourseByID)

	// Lessons
	courses.Get("/:course_id/lessons", lessonHandler.GetLessonsByCourseID)
	courses.Get("/:course_id/lessons/:lesson_id", lessonHandler.GetLessonByID)

}
