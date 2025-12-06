package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

// Сервисный слой для работы с категориями. Здесь только работа с БД и логирование,
// без HTTP-деталей.

func GetCategoriesService(ctx context.Context) ([]Category, error) {
	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d
		ORDER BY title
	`

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		log.Printf("❌ Ошибка запроса категорий: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Title, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			log.Printf("❌ Ошибка сканирования категории: %v", err)
			continue
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Ошибка итерации категорий: %v", err)
	}

	log.Printf("📋 Получено категорий: %d", len(categories))
	return categories, nil
}

func CreateCategoryService(ctx context.Context, title string) (Category, error) {
	query := `
		INSERT INTO knowledge_base.category_d (id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, created_at, updated_at
	`

	categoryID := uuid.NewString()
	now := time.Now()

	var category Category
	err := dbPool.QueryRow(ctx, query,
		categoryID,
		title,
		now,
		now,
	).Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		log.Printf("❌ Ошибка создания категории: %v", err)
		return Category{}, err
	}

	log.Printf("✅ Создана категория: %s (ID: %s)", category.Title, category.ID)
	return category, nil
}

func GetCategoryService(ctx context.Context, id string) (Category, error) {
	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d 
		WHERE id = $1
	`

	var category Category
	err := dbPool.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Title,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка получения категории: %v", err)
		return Category{}, err
	}

	log.Printf("📖 Получена категория: %s", category.Title)
	return category, nil
}

func UpdateCategoryService(ctx context.Context, id, title string) (Category, error) {
	updateQuery := `
		UPDATE knowledge_base.category_d 
		SET title = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, title, created_at, updated_at
	`

	var category Category
	err := dbPool.QueryRow(ctx, updateQuery,
		title,
		time.Now(),
		id,
	).Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		log.Printf("❌ Ошибка обновления категории: %v", err)
		return Category{}, err
	}

	log.Printf("✅ Обновлена категория: %s", category.Title)
	return category, nil
}

func CountCoursesForCategory(ctx context.Context, categoryID string) (int, error) {
	checkQuery := `
		SELECT COUNT(*) 
		FROM knowledge_base.course_b 
		WHERE category_id = $1
	`

	var courseCount int
	err := dbPool.QueryRow(ctx, checkQuery, categoryID).Scan(&courseCount)
	if err != nil {
		log.Printf("❌ Ошибка проверки курсов: %v", err)
		return 0, err
	}

	return courseCount, nil
}

func DeleteCategoryService(ctx context.Context, categoryID string) (int64, error) {
	deleteQuery := `DELETE FROM knowledge_base.category_d WHERE id = $1`
	result, err := dbPool.Exec(ctx, deleteQuery, categoryID)
	if err != nil {
		log.Printf("❌ Ошибка удаления категории: %v", err)
		return 0, err
	}

	log.Printf("🗑️ Удалена категория с ID: %s", categoryID)
	return result.RowsAffected(), nil
}
