package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Сервисный слой для работы с курсами.

func GetCoursesService(ctx context.Context, page, limit int, level, visibility, categoryID string) (PaginatedCourses, error) {
	baseQuery := `
		SELECT id, title, description, level, category_id, visibility, created_at, updated_at
		FROM knowledge_base.course_b
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM knowledge_base.course_b WHERE 1=1`

	var queryParams []interface{}
	var countParams []interface{}
	paramCounter := 1

	if level != "" {
		baseQuery += fmt.Sprintf(" AND level = $%d", paramCounter)
		countQuery += fmt.Sprintf(" AND level = $%d", paramCounter)
		queryParams = append(queryParams, level)
		countParams = append(countParams, level)
		paramCounter++
	}

	if visibility != "" {
		baseQuery += fmt.Sprintf(" AND visibility = $%d", paramCounter)
		countQuery += fmt.Sprintf(" AND visibility = $%d", paramCounter)
		queryParams = append(queryParams, visibility)
		countParams = append(countParams, visibility)
		paramCounter++
	}

	if categoryID != "" {
		baseQuery += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		countQuery += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		queryParams = append(queryParams, categoryID)
		countParams = append(countParams, categoryID)
		paramCounter++
	}

	baseQuery += " ORDER BY created_at DESC"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	queryParams = append(queryParams, limit, (page-1)*limit)

	log.Printf("🔍 Поиск курсов: page=%d, limit=%d, level=%s, visibility=%s, category=%s",
		page, limit, level, visibility, categoryID)

	var total int
	if err := dbPool.QueryRow(ctx, countQuery, countParams...).Scan(&total); err != nil {
		log.Printf("❌ Ошибка подсчета курсов: %v", err)
		return PaginatedCourses{}, err
	}

	rows, err := dbPool.Query(ctx, baseQuery, queryParams...)
	if err != nil {
		log.Printf("❌ Ошибка запроса курсов: %v", err)
		return PaginatedCourses{}, err
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Level,
			&course.CategoryID,
			&course.Visibility,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			log.Printf("❌ Ошибка сканирования курса: %v", err)
			continue
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Ошибка итерации курсов: %v", err)
	}

	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	response := PaginatedCourses{
		Data:  courses,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}

	log.Printf("✅ Найдено курсов: %d (показано: %d)", total, len(courses))
	return response, nil
}

func CreateCourseService(ctx context.Context, input CourseCreateDTO) (Course, error) {
	checkQuery := `SELECT id FROM knowledge_base.category_d WHERE id = $1`
	var categoryExists string
	if err := dbPool.QueryRow(ctx, checkQuery, input.CategoryID).Scan(&categoryExists); err != nil {
		log.Printf("❌ Ошибка проверки категории: %v", err)
		return Course{}, err
	}

	courseID := uuid.NewString()
	now := time.Now()

	query := `
		INSERT INTO knowledge_base.course_b 
		(id, title, description, level, category_id, visibility, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, level, category_id, visibility, created_at, updated_at
	`

	var course Course
	err := dbPool.QueryRow(ctx, query,
		courseID,
		input.Title,
		input.Description,
		input.Level,
		input.CategoryID,
		input.Visibility,
		now,
		now,
	).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Level,
		&course.CategoryID,
		&course.Visibility,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка создания курса: %v", err)
		return Course{}, err
	}

	log.Printf("✅ Создан курс: %s (ID: %s)", course.Title, course.ID)
	return course, nil
}

func GetCourseService(ctx context.Context, courseID string) (Course, error) {
	query := `
		SELECT id, title, description, level, category_id, visibility, created_at, updated_at
		FROM knowledge_base.course_b
		WHERE id = $1
	`

	var course Course
	err := dbPool.QueryRow(ctx, query, courseID).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Level,
		&course.CategoryID,
		&course.Visibility,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка получения курса: %v", err)
		return Course{}, err
	}

	log.Printf("📖 Получен курс: %s", course.Title)
	return course, nil
}

func UpdateCourseService(ctx context.Context, courseID string, input CourseUpdateDTO) (Course, error) {
	// Проверяем существование категории, если передан category_id
	if input.CategoryID != "" {
		var categoryExists string
		if err := dbPool.QueryRow(ctx, "SELECT id FROM knowledge_base.category_d WHERE id = $1", input.CategoryID).Scan(&categoryExists); err != nil {
			log.Printf("❌ Ошибка проверки категории: %v", err)
			return Course{}, err
		}
	}

	updateQuery := `UPDATE knowledge_base.course_b SET `
	var params []interface{}
	paramCounter := 1

	if input.Title != "" {
		updateQuery += fmt.Sprintf("title = $%d, ", paramCounter)
		params = append(params, input.Title)
		paramCounter++
	}

	if input.Description != "" {
		updateQuery += fmt.Sprintf("description = $%d, ", paramCounter)
		params = append(params, input.Description)
		paramCounter++
	}

	if input.Level != "" {
		updateQuery += fmt.Sprintf("level = $%d, ", paramCounter)
		params = append(params, input.Level)
		paramCounter++
	}

	if input.CategoryID != "" {
		updateQuery += fmt.Sprintf("category_id = $%d, ", paramCounter)
		params = append(params, input.CategoryID)
		paramCounter++
	}

	if input.Visibility != "" {
		updateQuery += fmt.Sprintf("visibility = $%d, ", paramCounter)
		params = append(params, input.Visibility)
		paramCounter++
	}

	updateQuery += fmt.Sprintf("updated_at = $%d ", paramCounter)
	params = append(params, time.Now())
	paramCounter++

	updateQuery += fmt.Sprintf("WHERE id = $%d ", paramCounter)
	params = append(params, courseID)
	paramCounter++

	updateQuery += "RETURNING id, title, description, level, category_id, visibility, created_at, updated_at"

	var course Course
	err := dbPool.QueryRow(ctx, updateQuery, params...).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.Level,
		&course.CategoryID,
		&course.Visibility,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка обновления курса: %v", err)
		return Course{}, err
	}

	log.Printf("✅ Обновлен курс: %s", course.Title)
	return course, nil
}

func DeleteCourseService(ctx context.Context, courseID string) (int64, error) {
	deleteQuery := `DELETE FROM knowledge_base.course_b WHERE id = $1`
	result, err := dbPool.Exec(ctx, deleteQuery, courseID)
	if err != nil {
		log.Printf("❌ Ошибка удаления курса: %v", err)
		return 0, err
	}

	rows := result.RowsAffected()
	if rows > 0 {
		log.Printf("🗑️ Удален курс с ID: %s", courseID)
	}

	return rows, nil
}

func GetCategoryCoursesService(ctx context.Context, categoryID string) ([]Course, error) {
	query := `
		SELECT id, title, description, level, category_id, visibility, created_at, updated_at
		FROM knowledge_base.course_b
		WHERE category_id = $1
		ORDER BY created_at DESC
	`

	rows, err := dbPool.Query(ctx, query, categoryID)
	if err != nil {
		log.Printf("❌ Ошибка запроса курсов категории: %v", err)
		return nil, err
	}
	defer rows.Close()

	var coursesList []Course
	for rows.Next() {
		var course Course
		if err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.Level,
			&course.CategoryID,
			&course.Visibility,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			log.Printf("❌ Ошибка сканирования курса: %v", err)
			continue
		}
		coursesList = append(coursesList, course)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Ошибка итерации курсов: %v", err)
	}

	log.Printf("📚 Найдено курсов в категории: %d", len(coursesList))
	return coursesList, nil
}
