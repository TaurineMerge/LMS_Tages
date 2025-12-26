// Package repository provides the data persistence layer for the application.
// It abstracts the database interactions for domain models.
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

// CategoryRepository defines the interface for database operations on categories.
type CategoryRepository interface {
	// GetAll retrieves a paginated list of all categories.
	GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	GetAllNotEmpty(ctx context.Context, page, limit int) ([]domain.Category, int, error)
	// GetByID retrieves a single category by its ID.
	GetByID(ctx context.Context, categoryID string) (domain.Category, error)
}

type categoryRepository struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

// NewCategoryRepository creates a new instance of a category repository.
func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

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

func (r *categoryRepository) GetAll(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	// Count total categories
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

	// Get paginated categories
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

func (r *categoryRepository) GetAllNotEmpty(ctx context.Context, page, limit int) ([]domain.Category, int, error) {
	// Count total categories that have at least one public course
	countQuery := r.psql.Select("COUNT(DISTINCT c.id)").
		From(categoryTable + " AS c").
		Join(courseTable + " AS co ON c.id = co.category_id").
		Where(squirrel.Eq{
			"visibility":  "public",
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

	// Get paginated categories that have at least one public course
	queryBuilder := r.psql.Select("c.id", "c.title", "c.created_at", "c.updated_at").
		From(categoryTable + " AS c").
		Join(courseTable + " AS co ON c.id = co.category_id").
		Where(squirrel.Eq{
			"visibility":  "public",
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
