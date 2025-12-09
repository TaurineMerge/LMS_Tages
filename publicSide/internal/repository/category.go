package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	categoriesTable = "categories"
)

type CategoryRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	GetByID(ctx context.Context, id string) (domain.Category, error)
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", categoriesTable)
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	query := fmt.Sprintf("SELECT id, title, created_at, updated_at FROM %s ORDER BY title LIMIT $1 OFFSET $2", categoriesTable)
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, total, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (domain.Category, error) {
	var category domain.Category
	query := fmt.Sprintf("SELECT id, title, created_at, updated_at FROM %s WHERE id = $1", categoriesTable)

	err := r.db.QueryRow(ctx, query, id).Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Category{}, fmt.Errorf("category with id %s not found", id)
		}
		return domain.Category{}, fmt.Errorf("failed to get category by id: %w", err)
	}

	return category, nil
}
