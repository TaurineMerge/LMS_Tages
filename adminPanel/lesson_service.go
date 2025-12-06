package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"fmt"

	"github.com/google/uuid"
)

// Сервисный слой для работы с уроками.

func GetCourseLessonsService(ctx context.Context, courseID string) ([]Lesson, error) {
	query := `
		SELECT id, title, course_id, content, created_at, updated_at
		FROM knowledge_base.lesson_d
		WHERE course_id = $1
		ORDER BY created_at
	`

	rows, err := dbPool.Query(ctx, query, courseID)
	if err != nil {
		log.Printf("❌ Ошибка запроса уроков: %v", err)
		return nil, err
	}
	defer rows.Close()

	var lessonsList []Lesson
	for rows.Next() {
		var lesson Lesson
		var contentJSON []byte

		if err := rows.Scan(
			&lesson.ID,
			&lesson.Title,
			&lesson.CourseID,
			&contentJSON,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			log.Printf("❌ Ошибка сканирования урока: %v", err)
			continue
		}

		if len(contentJSON) > 0 {
			if err := json.Unmarshal(contentJSON, &lesson.Content); err != nil {
				log.Printf("❌ Ошибка парсинга content урока: %v", err)
				lesson.Content = make(map[string]interface{})
			}
		} else {
			lesson.Content = make(map[string]interface{})
		}

		lessonsList = append(lessonsList, lesson)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Ошибка итерации уроков: %v", err)
	}

	log.Printf("📖 Получено уроков для курса %s: %d", courseID, len(lessonsList))
	return lessonsList, nil
}

func CreateLessonService(ctx context.Context, courseID, courseTitle, title string, contentJSON []byte) (Lesson, error) {
	lessonID := uuid.NewString()
	now := time.Now()

	query := `
		INSERT INTO knowledge_base.lesson_d 
		(id, title, course_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, course_id, content, created_at, updated_at
	`

	var lesson Lesson
	var contentBytes []byte

	err := dbPool.QueryRow(ctx, query,
		lessonID,
		title,
		courseID,
		contentJSON,
		now,
		now,
	).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&contentBytes,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка создания урока: %v", err)
		return Lesson{}, err
	}

	if len(contentBytes) > 0 {
		if err := json.Unmarshal(contentBytes, &lesson.Content); err != nil {
			log.Printf("❌ Ошибка парсинга content: %v", err)
			lesson.Content = make(map[string]interface{})
		}
	} else {
		lesson.Content = make(map[string]interface{})
	}

	log.Printf("✅ Создан урок: %s (ID: %s) для курса: %s",
		lesson.Title, lesson.ID, courseTitle)

	return lesson, nil
}

func GetLessonService(ctx context.Context, courseID, lessonID string) (Lesson, error) {
	query := `
		SELECT id, title, course_id, content, created_at, updated_at
		FROM knowledge_base.lesson_d
		WHERE id = $1 AND course_id = $2
	`

	var lesson Lesson
	var contentJSON []byte

	err := dbPool.QueryRow(ctx, query, lessonID, courseID).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&contentJSON,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка получения урока: %v", err)
		return Lesson{}, err
	}

	if len(contentJSON) > 0 {
		if err := json.Unmarshal(contentJSON, &lesson.Content); err != nil {
			log.Printf("❌ Ошибка парсинга content урока: %v", err)
			lesson.Content = make(map[string]interface{})
		}
	} else {
		lesson.Content = make(map[string]interface{})
	}

	log.Printf("📖 Получен урок: %s", lesson.Title)
	return lesson, nil
}

func UpdateLessonService(ctx context.Context, courseID, lessonID string, input LessonUpdateDTO, contentJSON []byte, updateContent bool) (Lesson, error) {
	updateQuery := `UPDATE knowledge_base.lesson_d SET `
	var params []interface{}
	paramCounter := 1

	if input.Title != "" {
		updateQuery += fmt.Sprintf("title = $%d, ", paramCounter)
		params = append(params, input.Title)
		paramCounter++
	}

	if updateContent {
		updateQuery += fmt.Sprintf("content = $%d, ", paramCounter)
		params = append(params, contentJSON)
		paramCounter++
	}

	updateQuery += fmt.Sprintf("updated_at = $%d ", paramCounter)
	params = append(params, time.Now())
	paramCounter++

	updateQuery += fmt.Sprintf("WHERE id = $%d AND course_id = $%d ", paramCounter, paramCounter+1)
	params = append(params, lessonID, courseID)
	paramCounter += 2

	updateQuery += "RETURNING id, title, course_id, content, created_at, updated_at"

	var lesson Lesson
	var contentBytes []byte

	err := dbPool.QueryRow(ctx, updateQuery, params...).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&contentBytes,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка обновления урока: %v", err)
		return Lesson{}, err
	}

	if len(contentBytes) > 0 {
		if err := json.Unmarshal(contentBytes, &lesson.Content); err != nil {
			log.Printf("❌ Ошибка парсинга content: %v", err)
			lesson.Content = make(map[string]interface{})
		}
	} else {
		lesson.Content = make(map[string]interface{})
	}

	log.Printf("✅ Обновлен урок: %s", lesson.Title)
	return lesson, nil
}

func DeleteLessonService(ctx context.Context, courseID, lessonID string) (string, int64, error) {
	checkQuery := `SELECT id, title FROM knowledge_base.lesson_d WHERE id = $1 AND course_id = $2`
	var lessonTitle string
	if err := dbPool.QueryRow(ctx, checkQuery, lessonID, courseID).Scan(&lessonID, &lessonTitle); err != nil {
		log.Printf("❌ Ошибка проверки урока: %v", err)
		return "", 0, err
	}

	deleteQuery := `DELETE FROM knowledge_base.lesson_d WHERE id = $1 AND course_id = $2`
	result, err := dbPool.Exec(ctx, deleteQuery, lessonID, courseID)
	if err != nil {
		log.Printf("❌ Ошибка удаления урока: %v", err)
		return "", 0, err
	}

	rows := result.RowsAffected()
	if rows > 0 {
		log.Printf("🗑️ Удален урок: %s", lessonTitle)
	}

	return lessonTitle, rows, nil
}
