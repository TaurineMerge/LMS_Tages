// Package repository предоставляет слой для взаимодействия с базой данных.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// DirectionNext указывает на выборку следующего элемента.
	DirectionNext = "next"
	// DirectionPrevious указывает на выборку предыдущего элемента.
	DirectionPrevious = "previous"
)

// LessonChunkOptions определяет параметры для выборки "чанка" (порции) уроков.
// Используется для получения соседних уроков.
type LessonChunkOptions struct {
	PivotValue interface{} // Значение поля, от которого идет выборка (например, `created_at` текущего урока).
	OrderBy    string      // Поле для сортировки.
	Direction  string      // Направление выборки (`next` или `previous`).
	Limit      int         // Количество записей для выборки.
}

// LessonRepository определяет интерфейс для работы с уроками в базе данных.
type LessonRepository interface {
	// GetAllByCourseID получает все уроки для данного курса с пагинацией и сортировкой.
	GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]domain.Lesson, int, error)
	// GetByID получает один урок по его ID, ID курса и ID категории.
	GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error)
	// GetLessonsChunk получает порцию уроков на основе заданных опций.
	GetLessonsChunk(ctx context.Context, courseID string, options LessonChunkOptions) ([]domain.Lesson, error)
}

// lessonRepository является реализацией LessonRepository.
type lessonRepository struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

// NewLessonRepository создает новый экземпляр lessonRepository.
func NewLessonRepository(db *pgxpool.Pool) LessonRepository {
	return &lessonRepository{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// scanner - это общий интерфейс для pgx.Row и pgx.Rows, чтобы избежать дублирования.
type scanner interface {
	Scan(dest ...any) error
}

// scanLesson сканирует одну строку из результата запроса в структуру domain.Lesson.
func (r *lessonRepository) scanLesson(row scanner) (domain.Lesson, error) {
	var lesson domain.Lesson
	err := row.Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&lesson.Content,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)
	return lesson, err
}

// scanLessons итерирует по pgx.Rows и сканирует каждую строку в срез []domain.Lesson.
func (r *lessonRepository) scanLessons(rows pgx.Rows) ([]domain.Lesson, error) {
	var lessons []domain.Lesson
	defer rows.Close()
	for rows.Next() {
		lesson, err := r.scanLesson(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}
		lessons = append(lessons, lesson)
	}
	return lessons, nil
}

// GetAllByCourseID извлекает срез уроков для указанного курса с пагинацией и сортировкой.
// Возвращает срез уроков, общее количество уроков в курсе и ошибку.
func (r *lessonRepository) GetAllByCourseID(ctx context.Context, categoryID, courseID string, page, limit int, sort string) ([]domain.Lesson, int, error) {
	// Сначала получаем общее количество уроков для пагинации.
	countBuilder := r.psql.Select("COUNT(l.id)").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"c.visibility":  "public",
		})

	countQuery, args, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query for lessons: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count lessons by course: %w", err)
	}

	if total == 0 {
		return []domain.Lesson{}, 0, nil
	}

	// Затем получаем срез уроков для текущей страницы.
	queryBuilder := r.psql.Select("l.id", "l.title", "l.course_id", "l.content", "l.created_at", "l.updated_at").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"c.visibility":  "public",
		}).
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit))

	queryBuilder = r.applySorting(queryBuilder, sort)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build get all lessons query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get lessons by course: %w", err)
	}

	lessons, err := r.scanLessons(rows)
	if err != nil {
		return nil, 0, err
	}
	return lessons, total, nil
}

// GetByID находит и возвращает один видимый урок по его ID, ID курса и ID категории.
// Если урок не найден, возвращает ошибку.
func (r *lessonRepository) GetByID(ctx context.Context, categoryID, courseID, lessonID string) (domain.Lesson, error) {
	queryBuilder := r.psql.Select("l.id", "l.title", "l.course_id", "l.content", "l.created_at", "l.updated_at").
		From(lessonsTable + " AS l").
		Join(courseTable + " AS c ON l.course_id = c.id").
		Where(squirrel.Eq{
			"c.category_id": categoryID,
			"l.course_id":   courseID,
			"l.id":          lessonID,
			"c.visibility":  "public",
		})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return domain.Lesson{}, fmt.Errorf("failed to build get lesson by id query: %w", err)
	}

	row := r.db.QueryRow(ctx, query, args...)
	lesson, err := r.scanLesson(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Lesson{}, fmt.Errorf("lesson with id %s not found in course %s", lessonID, courseID)
		}
		return domain.Lesson{}, fmt.Errorf("failed to get lesson by id: %w", err)
	}

	return lesson, nil
}

// GetLessonsChunk получает "порцию" уроков (следующий или предыдущий) относительно опорного урока.
// Это используется для навигации "следующий/предыдущий урок".
func (r *lessonRepository) GetLessonsChunk(ctx context.Context, courseID string, options LessonChunkOptions) ([]domain.Lesson, error) {
	if !r.isValidOrderBy(options.OrderBy) {
		return nil, fmt.Errorf("invalid order by field: %s", options.OrderBy)
	}

	queryBuilder := r.psql.Select("l.id", "l.title", "l.course_id", "l.content", "l.created_at", "l.updated_at").
		From(lessonsTable + " AS l").
		Where(squirrel.Eq{"l.course_id": courseID})

	// Устанавливаем условие для выборки относительно опорного значения.
	if options.PivotValue != nil {
		column := fmt.Sprintf("l.%s", options.OrderBy)
		if options.Direction == DirectionNext {
			queryBuilder = queryBuilder.Where(squirrel.Gt{column: options.PivotValue})
		} else {
			queryBuilder = queryBuilder.Where(squirrel.Lt{column: options.PivotValue})
		}
	}

	// Устанавливаем порядок сортировки для корректной выборки "следующего" или "предыдущего".
	if options.Direction == DirectionNext {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("l.%s ASC", options.OrderBy))
	} else {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("l.%s DESC", options.OrderBy))
	}

	queryBuilder = queryBuilder.Limit(uint64(options.Limit))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build lessons chunk query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons chunk: %w", err)
	}

	return r.scanLessons(rows)
}

// isValidOrderBy проверяет, является ли поле сортировки допустимым.
func (r *lessonRepository) isValidOrderBy(field string) bool {
	switch field {
	case "created_at", "title", "updated_at":
		return true
	default:
		return false
	}
}

// applySorting применяет к запросу сортировку на основе строки `sort`.
func (r *lessonRepository) applySorting(builder squirrel.SelectBuilder, sort string) squirrel.SelectBuilder {
	if sort == "" {
		return builder.OrderBy("l.created_at ASC")
	}

	allowedFields := map[string]string{
		"title":      "l.title",
		"created_at": "l.created_at",
		"updated_at": "l.updated_at",
	}

	direction := "ASC"
	if strings.HasPrefix(sort, "-") {
		direction = "DESC"
		sort = strings.TrimPrefix(sort, "-")
	}

	dbColumn, ok := allowedFields[sort]
	if !ok {
		return builder.OrderBy("l.created_at ASC") // Сортировка по умолчанию, если поле не разрешено
	}

	return builder.OrderBy(fmt.Sprintf("%s %s", dbColumn, direction))
}
