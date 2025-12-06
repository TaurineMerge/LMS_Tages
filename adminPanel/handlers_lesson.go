package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ============ LESSONS & COURSE LESSONS ============

// GET /courses/:course_id/lessons
func getCourseLessons(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(newBadRequestError("Invalid course ID format"))
	}

	ctx := context.Background()
	// Проверяем, что курс существует
	if _, err := GetCourseService(ctx, courseID); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Course not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to check course"))
	}

	lessons, err := GetCourseLessonsService(ctx, courseID)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to get lessons"))
	}

	return c.JSON(toLessonDTOs(lessons))
}

// GET /courses/:course_id/lessons (alias)
func getLessons(c *fiber.Ctx) error {
	return getCourseLessons(c)
}

// POST /courses/:course_id/lessons
func createLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(newBadRequestError("Invalid course ID format"))
	}

	ctx := context.Background()
	course, err := GetCourseService(ctx, courseID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Course not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to check course"))
	}

	var input LessonCreateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title == "" {
		return c.Status(400).JSON(newValidationError("Title is required"))
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(newValidationError("Title must be less than 255 characters"))
	}

	contentJSON := []byte("{}")
	if input.Content != nil {
		var err error
		contentJSON, err = json.Marshal(input.Content)
		if err != nil {
			log.Printf("❌ Ошибка маршалинга content: %v", err)
			return c.Status(400).JSON(newBadRequestError("Invalid content format"))
		}
	}

	lesson, err := CreateLessonService(ctx, courseID, course.Title, input.Title, contentJSON)
	if err != nil {
		return c.Status(500).JSON(newInternalError("Failed to create lesson"))
	}

	return c.Status(201).JSON(toLessonDTO(lesson))
}

// GET /courses/:course_id/lessons/:lesson_id
func getLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(newBadRequestError("Invalid ID format"))
	}

	ctx := context.Background()
	// Проверяем, что курс существует
	if _, err := GetCourseService(ctx, courseID); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Course not found"))
		}
		return c.Status(500).JSON(newInternalError("Failed to check course"))
	}

	lesson, err := GetLessonService(ctx, courseID, lessonID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Lesson not found in this course"))
		}
		return c.Status(500).JSON(newInternalError("Failed to get lesson"))
	}

	return c.JSON(toLessonDTO(lesson))
}

// PUT /courses/:course_id/lessons/:lesson_id
func updateLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(newBadRequestError("Invalid ID format"))
	}

	var input LessonUpdateDTO
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(newBadRequestError("Invalid request body"))
	}

	if input.Title != "" && len(input.Title) > 255 {
		return c.Status(400).JSON(newValidationError("Title must be less than 255 characters"))
	}

	var contentJSON []byte
	var updateContent bool
	if input.Content != nil {
		var err error
		contentJSON, err = json.Marshal(input.Content)
		if err != nil {
			log.Printf("❌ Ошибка маршалинга content: %v", err)
			return c.Status(400).JSON(newBadRequestError("Invalid content format"))
		}
		updateContent = true
	}

	ctx := context.Background()
	lesson, err := UpdateLessonService(ctx, courseID, lessonID, input, contentJSON, updateContent)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Lesson not found in this course"))
		}
		return c.Status(500).JSON(newInternalError("Failed to update lesson"))
	}

	return c.JSON(toLessonDTO(lesson))
}

// DELETE /courses/:course_id/lessons/:lesson_id
func deleteLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(newBadRequestError("Invalid ID format"))
	}

	ctx := context.Background()
	lessonTitle, rows, err := DeleteLessonService(ctx, courseID, lessonID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(newNotFoundError("Lesson not found in this course"))
		}
		return c.Status(500).JSON(newInternalError("Failed to delete lesson"))
	}

	if rows == 0 {
		return c.Status(404).JSON(newNotFoundError("Lesson not found"))
	}

	log.Printf("🗑️ Удален урок: %s", lessonTitle)
	return c.SendStatus(204)
}
