// Package repository предоставляет слой для взаимодействия с базой данных.
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CategoryRepository определяет интерфейс для работы с категориями в базе данных.
type CategoryRepository interface {
	// GetAll получает все категории с пагинацией.
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	// GetAllNotEmpty получает все категории, в которых есть хотя бы один публичный курс, с пагинацией.
	GetAllNotEmpty(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	// GetByID получает категорию по ее уникальному идентификатору.
	GetByID(ctx context.Context, categoryID string) (domain.Category, error)
}

// categoryRepository является реализацией CategoryRepository.
type categoryRepository struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

// NewCategoryRepository создает новый экземпляр categoryRepository.
func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// scanCategory сканирует одну строку из результата запроса в структуру domain.Category.
func (r *categoryRepository) scanCategory(row scanner) (domain.Category, error) {
	var category domain.Category
	err := row.Scan(
		&category.ID,
		&category.Title,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	return category, err
}

// GetAll извлекает из базы данных срез категорий с учетом пагинации.
// Возвращает срез категорий, общее количество категорий и ошибку.
func (r *categoryRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	countQuery := r.psql.Select("COUNT(*)").
		From(categoryTable)

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query for categories: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, countSql, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	if total == 0 {
		return []domain.Category{}, 0, nil
	}

	queryBuilder := r.psql.Select("id", "title", "created_at", "updated_at").
		From(categoryTable).
		OrderBy("created_at ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build get all categories query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category, err := r.scanCategory(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// GetAllNotEmpty извлекает категории, содержащие хотя бы один видимый курс.
// Используется для отображения только тех категорий, которые имеют контент.
func (r *categoryRepository) GetAllNotEmpty(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	countQuery := r.psql.Select("COUNT(DISTINCT c.id)").
		From(categoryTable + " AS c").
		Join(courseTable + " AS co ON c.id = co.category_id").
		Where(squirrel.Eq{
			"visibility": "public",
		})

	countSql, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query for not empty categories: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, countSql, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count not empty categories: %w", err)
	}

	if total == 0 {
		return []domain.Category{}, 0, nil
	}

	queryBuilder := r.psql.Select("c.id", "c.title", "c.created_at", "c.updated_at").
		From(categoryTable+" AS c").
		Join(courseTable+" AS co ON c.id = co.category_id").
		Where(squirrel.Eq{
			"visibility": "public",
		}).
		GroupBy("c.id", "c.title", "c.created_at", "c.updated_at").
		OrderBy("c.created_at ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build get all not empty categories query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get not empty categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		category, err := r.scanCategory(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

// GetByID находит и возвращает одну категорию по её ID.
// Если категория не найдена, возвращает ошибку, содержащую pgx.ErrNoRows.
func (r *categoryRepository) GetByID(ctx context.Context, categoryID string) (domain.Category, error) {
	queryBuilder := r.psql.Select("id", "title", "created_at", "updated_at").
		From(categoryTable).
		Where(squirrel.Eq{
			"id": categoryID,
		})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return domain.Category{}, fmt.Errorf("failed to build get category by id query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	category, err := r.scanCategory(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, fmt.Errorf("category with id %s not found", categoryID)
		}
		return domain.Category{}, fmt.Errorf("failed to get category by id: %w", err)
	}

	return category, nil
}
