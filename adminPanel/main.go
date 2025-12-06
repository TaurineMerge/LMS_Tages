package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"context"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	
	// Swagger UI
	"github.com/gofiber/swagger"
)

// ============ КОНФИГУРАЦИЯ ============

type Config struct {
	DatabaseURL string
}

func getConfig() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://appuser:password@app-db:5432/appdb?sslmode=disable"
	}

	return Config{
		DatabaseURL: dbURL,
	}
}

// ============ ПУЛ ПОДКЛЮЧЕНИЙ ============

var dbPool *pgxpool.Pool

func initDB() error {
	config := getConfig()

	poolConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("ошибка парсинга конфигурации БД: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 5

	ctx := context.Background()
	dbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("ошибка создания пула соединений: %w", err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	log.Println("✅ Успешно подключено к PostgreSQL")
	return nil
}

// ============ МОДЕЛИ ============

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type Category struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryCreate struct {
	Title string `json:"title"`
}

type CategoryUpdate struct {
	Title string `json:"title"`
}

type Course struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	CategoryID  string    `json:"category_id"`
	Visibility  string    `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseCreate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

type CourseUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

type PaginatedCourses struct {
	Data  []Course `json:"data"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Pages int      `json:"pages"`
}

type Lesson struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CourseID  string                 `json:"course_id"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type LessonCreate struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

type LessonUpdate struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

// ============ MAIN ============

func main() {
	// Инициализация базы данных
	if err := initDB(); err != nil {
		log.Fatalf("❌ Ошибка инициализации БД: %v", err)
	}
	defer dbPool.Close()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Middleware
	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Next()
	})

	// Swagger UI (использует ваш существующий swagger.json)
	setupSwagger(app)

	// API routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", healthCheck)

	// Categories CRUD
	categories := api.Group("/categories")
	categories.Get("/", getCategories)
	categories.Post("/", createCategory)
	categories.Get("/:category_id", getCategory)
	categories.Put("/:category_id", updateCategory)
	categories.Delete("/:category_id", deleteCategory)
	categories.Get("/:category_id/courses", getCategoryCourses)

	// Courses CRUD
	courses := api.Group("/courses")
	courses.Get("/", getCourses)
	courses.Post("/", createCourse)
	courses.Get("/:course_id", getCourse)
	courses.Put("/:course_id", updateCourse)
	courses.Delete("/:course_id", deleteCourse)
	courses.Get("/:course_id/lessons", getCourseLessons)

	// Lessons CRUD
	lessons := courses.Group("/:course_id/lessons")
	lessons.Get("/", getLessons)
	lessons.Post("/", createLesson)
	lessons.Get("/:lesson_id", getLesson)
	lessons.Put("/:lesson_id", updateLesson)
	lessons.Delete("/:lesson_id", deleteLesson)

	// Start server
	log.Println("🚀 Сервер запущен на :4000")
	log.Println("📚 API доступен по адресу: http://localhost:4000/api")
	log.Println("📄 Swagger UI: http://localhost:4000/swagger/index.html")
	log.Fatal(app.Listen(":4000"))
}

// ============ SWAGGER SETUP ============

func setupSwagger(app *fiber.App) {
	// 1. Endpoint для вашего swagger.json
	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			log.Printf("❌ Ошибка чтения swagger.json: %v", err)
			return c.Status(500).JSON(ErrorResponse{
				Error: "Failed to load Swagger documentation",
				Code:  "INTERNAL_ERROR",
			})
		}

		var swaggerSpec map[string]interface{}
		if err := json.Unmarshal(data, &swaggerSpec); err != nil {
			log.Printf("❌ Ошибка парсинга swagger.json: %v", err)
			return c.Status(500).JSON(ErrorResponse{
				Error: "Invalid Swagger JSON",
				Code:  "INTERNAL_ERROR",
			})
		}

		// Обновляем basePath и host для вашего API
		swaggerSpec["basePath"] = "/api/v1"
		swaggerSpec["host"] = "localhost:4000"

		return c.JSON(swaggerSpec)
	})

	// 2. Swagger UI от библиотеки (будет использовать наш /swagger.json)
	app.Get("/swagger/*", swagger.New(swagger.Config{
	URL:          "/swagger.json",
	DeepLinking:  true,
	DocExpansion: "list",
}))
}

// ============ HANDLERS ============

// Health check
func healthCheck(c *fiber.Ctx) error {
	ctx := context.Background()
	err := dbPool.Ping(ctx)
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
		log.Printf("❌ Ошибка проверки БД: %v", err)
	}

	return c.JSON(HealthResponse{
		Status:   "healthy",
		Database: dbStatus,
		Version:  "1.0.0",
	})
}

// ============ CATEGORIES ============

// GET /categories
func getCategories(c *fiber.Ctx) error {
	ctx := context.Background()

	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d
		ORDER BY title
	`

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		log.Printf("❌ Ошибка запроса категорий: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Internal server error",
			Code:  "INTERNAL_ERROR",
		})
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
	return c.JSON(categories)
}

// POST /categories
func createCategory(c *fiber.Ctx) error {
	var input CategoryCreate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title is required",
			Code:  "VALIDATION_ERROR",
		})
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title must be less than 255 characters",
			Code:  "VALIDATION_ERROR",
		})
	}

	ctx := context.Background()
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
		input.Title,
		now,
		now,
	).Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		log.Printf("❌ Ошибка создания категории: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(409).JSON(ErrorResponse{
				Error: "Category with this title already exists",
				Code:  "CONFLICT",
			})
		}
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create category",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("✅ Создана категория: %s (ID: %s)", category.Title, category.ID)
	return c.Status(201).JSON(category)
}

// GET /categories/:category_id
func getCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	query := `
		SELECT id, title, created_at, updated_at 
		FROM knowledge_base.category_d 
		WHERE id = $1
	`

	var category Category
	err := dbPool.QueryRow(ctx, query, categoryID).Scan(
		&category.ID,
		&category.Title,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Category not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка получения категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get category",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("📖 Получена категория: %s", category.Title)
	return c.JSON(category)
}

// PUT /categories/:category_id
func updateCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.category_d WHERE id = $1`
	var existingID string
	err := dbPool.QueryRow(ctx, checkQuery, categoryID).Scan(&existingID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Category not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check category",
			Code:  "INTERNAL_ERROR",
		})
	}

	var input CategoryUpdate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title is required for update",
			Code:  "VALIDATION_ERROR",
		})
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title must be less than 255 characters",
			Code:  "VALIDATION_ERROR",
		})
	}

	updateQuery := `
		UPDATE knowledge_base.category_d 
		SET title = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, title, created_at, updated_at
	`

	var category Category
	err = dbPool.QueryRow(ctx, updateQuery,
		input.Title,
		time.Now(),
		categoryID,
	).Scan(&category.ID, &category.Title, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		log.Printf("❌ Ошибка обновления категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update category",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("✅ Обновлена категория: %s", category.Title)
	return c.JSON(category)
}

// DELETE /categories/:category_id
func deleteCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `
		SELECT COUNT(*) 
		FROM knowledge_base.course_b 
		WHERE category_id = $1
	`

	var courseCount int
	err := dbPool.QueryRow(ctx, checkQuery, categoryID).Scan(&courseCount)
	if err != nil {
		log.Printf("❌ Ошибка проверки курсов: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check associated courses",
			Code:  "INTERNAL_ERROR",
		})
	}

	if courseCount > 0 {
		return c.Status(409).JSON(ErrorResponse{
			Error: "Cannot delete category with associated courses",
			Code:  "CONFLICT",
		})
	}

	deleteQuery := `DELETE FROM knowledge_base.category_d WHERE id = $1`
	result, err := dbPool.Exec(ctx, deleteQuery, categoryID)
	if err != nil {
		log.Printf("❌ Ошибка удаления категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to delete category",
			Code:  "INTERNAL_ERROR",
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(404).JSON(ErrorResponse{
			Error: "Category not found",
			Code:  "NOT_FOUND",
		})
	}

	log.Printf("🗑️ Удалена категория с ID: %s", categoryID)
	return c.SendStatus(204)
}

// GET /categories/:category_id/courses
func getCategoryCourses(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")

	if !isValidUUID(categoryID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.category_d WHERE id = $1`
	var existingID string
	err := dbPool.QueryRow(ctx, checkQuery, categoryID).Scan(&existingID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Category not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check category",
			Code:  "INTERNAL_ERROR",
		})
	}

	query := `
		SELECT id, title, description, level, category_id, visibility, created_at, updated_at
		FROM knowledge_base.course_b
		WHERE category_id = $1
		ORDER BY created_at DESC
	`

	rows, err := dbPool.Query(ctx, query, categoryID)
	if err != nil {
		log.Printf("❌ Ошибка запроса курсов категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get courses",
			Code:  "INTERNAL_ERROR",
		})
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

	log.Printf("📚 Найдено курсов в категории: %d", len(coursesList))
	return c.JSON(coursesList)
}

// ============ COURSES ============

// GET /courses
func getCourses(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	level := c.Query("level")
	visibility := c.Query("visibility")
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	ctx := context.Background()

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
		if !isValidUUID(categoryID) {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Invalid category ID format",
				Code:  "BAD_REQUEST",
			})
		}
		baseQuery += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		countQuery += fmt.Sprintf(" AND category_id = $%d", paramCounter)
		queryParams = append(queryParams, categoryID)
		countParams = append(countParams, categoryID)
		paramCounter++
	}

	baseQuery += " ORDER BY created_at DESC"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	queryParams = append(queryParams, limit, offset)

	log.Printf("🔍 Поиск курсов: page=%d, limit=%d, level=%s, visibility=%s, category=%s",
		page, limit, level, visibility, categoryID)

	var total int
	err := dbPool.QueryRow(ctx, countQuery, countParams...).Scan(&total)
	if err != nil {
		log.Printf("❌ Ошибка подсчета курсов: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to count courses",
			Code:  "INTERNAL_ERROR",
		})
	}

	rows, err := dbPool.Query(ctx, baseQuery, queryParams...)
	if err != nil {
		log.Printf("❌ Ошибка запроса курсов: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get courses",
			Code:  "INTERNAL_ERROR",
		})
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
	return c.JSON(response)
}

// POST /courses
func createCourse(c *fiber.Ctx) error {
	var input CourseCreate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title is required",
			Code:  "VALIDATION_ERROR",
		})
	}

	if input.CategoryID == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Category ID is required",
			Code:  "VALIDATION_ERROR",
		})
	}

	if !isValidUUID(input.CategoryID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid category ID format",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Level must be one of: hard, medium, easy",
			Code:  "VALIDATION_ERROR",
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Visibility must be one of: draft, public, private",
			Code:  "VALIDATION_ERROR",
		})
	}

	if input.Level == "" {
		input.Level = "medium"
	}
	if input.Visibility == "" {
		input.Visibility = "draft"
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.category_d WHERE id = $1`
	var categoryExists string
	err := dbPool.QueryRow(ctx, checkQuery, input.CategoryID).Scan(&categoryExists)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Category not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки категории: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check category",
			Code:  "INTERNAL_ERROR",
		})
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
	err = dbPool.QueryRow(ctx, query,
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
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create course",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("✅ Создан курс: %s (ID: %s)", course.Title, course.ID)
	return c.Status(201).JSON(course)
}

// GET /courses/:course_id
func getCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
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
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка получения курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get course",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("📖 Получен курс: %s", course.Title)
	return c.JSON(course)
}

// PUT /courses/:course_id
func updateCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.course_b WHERE id = $1`
	var existingID string
	err := dbPool.QueryRow(ctx, checkQuery, courseID).Scan(&existingID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check course",
			Code:  "INTERNAL_ERROR",
		})
	}

	var input CourseUpdate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title != "" {
		if len(input.Title) > 255 {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Title must be less than 255 characters",
				Code:  "VALIDATION_ERROR",
			})
		}
	}

	if input.Level != "" && !isValidLevel(input.Level) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Level must be one of: hard, medium, easy",
			Code:  "VALIDATION_ERROR",
		})
	}

	if input.Visibility != "" && !isValidVisibility(input.Visibility) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Visibility must be one of: draft, public, private",
			Code:  "VALIDATION_ERROR",
		})
	}

	if input.CategoryID != "" {
		if !isValidUUID(input.CategoryID) {
			return c.Status(400).JSON(ErrorResponse{
				Error: "Invalid category ID format",
				Code:  "BAD_REQUEST",
			})
		}
		var categoryExists string
		err := dbPool.QueryRow(ctx, "SELECT id FROM knowledge_base.category_d WHERE id = $1", input.CategoryID).Scan(&categoryExists)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return c.Status(404).JSON(ErrorResponse{
					Error: "Category not found",
					Code:  "NOT_FOUND",
				})
			}
			log.Printf("❌ Ошибка проверки категории: %v", err)
			return c.Status(500).JSON(ErrorResponse{
				Error: "Failed to check category",
				Code:  "INTERNAL_ERROR",
			})
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
	err = dbPool.QueryRow(ctx, updateQuery, params...).Scan(
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
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update course",
			Code:  "INTERNAL_ERROR",
		})
	}

	log.Printf("✅ Обновлен курс: %s", course.Title)
	return c.JSON(course)
}

// DELETE /courses/:course_id
func deleteCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.course_b WHERE id = $1`
	var existingID string
	err := dbPool.QueryRow(ctx, checkQuery, courseID).Scan(&existingID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check course",
			Code:  "INTERNAL_ERROR",
		})
	}

	deleteQuery := `DELETE FROM knowledge_base.course_b WHERE id = $1`
	result, err := dbPool.Exec(ctx, deleteQuery, courseID)
	if err != nil {
		log.Printf("❌ Ошибка удаления курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to delete course",
			Code:  "INTERNAL_ERROR",
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(404).JSON(ErrorResponse{
			Error: "Course not found",
			Code:  "NOT_FOUND",
		})
	}

	log.Printf("🗑️ Удален курс с ID: %s", courseID)
	return c.SendStatus(204)
}

// GET /courses/:course_id/lessons
func getCourseLessons(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.course_b WHERE id = $1`
	var courseExists string
	err := dbPool.QueryRow(ctx, checkQuery, courseID).Scan(&courseExists)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check course",
			Code:  "INTERNAL_ERROR",
		})
	}

	query := `
		SELECT id, title, course_id, content, created_at, updated_at
		FROM knowledge_base.lesson_d
		WHERE course_id = $1
		ORDER BY created_at
	`

	rows, err := dbPool.Query(ctx, query, courseID)
	if err != nil {
		log.Printf("❌ Ошибка запроса уроков: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get lessons",
			Code:  "INTERNAL_ERROR",
		})
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
	return c.JSON(lessonsList)
}

// ============ LESSONS ============

// GET /courses/:course_id/lessons
func getLessons(c *fiber.Ctx) error {
	return getCourseLessons(c)
}

// POST /courses/:course_id/lessons
func createLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")

	if !isValidUUID(courseID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid course ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id, title FROM knowledge_base.course_b WHERE id = $1`
	var courseTitle string
	err := dbPool.QueryRow(ctx, checkQuery, courseID).Scan(&courseID, &courseTitle)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check course",
			Code:  "INTERNAL_ERROR",
		})
	}

	var input LessonCreate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title == "" {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title is required",
			Code:  "VALIDATION_ERROR",
		})
	}

	if len(input.Title) > 255 {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title must be less than 255 characters",
			Code:  "VALIDATION_ERROR",
		})
	}

	contentJSON := []byte("{}")
	if input.Content != nil {
		var err error
		contentJSON, err = json.Marshal(input.Content)
		if err != nil {
			log.Printf("❌ Ошибка маршалинга content: %v", err)
			return c.Status(400).JSON(ErrorResponse{
				Error: "Invalid content format",
				Code:  "BAD_REQUEST",
			})
		}
	}

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

	err = dbPool.QueryRow(ctx, query,
		lessonID,
		input.Title,
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
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to create lesson",
			Code:  "INTERNAL_ERROR",
		})
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

	return c.Status(201).JSON(lesson)
}

// GET /courses/:course_id/lessons/:lesson_id
func getLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkCourseQuery := `SELECT id FROM knowledge_base.course_b WHERE id = $1`
	var courseExists string
	err := dbPool.QueryRow(ctx, checkCourseQuery, courseID).Scan(&courseExists)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Course not found",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки курса: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check course",
			Code:  "INTERNAL_ERROR",
		})
	}

	query := `
		SELECT id, title, course_id, content, created_at, updated_at
		FROM knowledge_base.lesson_d
		WHERE id = $1 AND course_id = $2
	`

	var lesson Lesson
	var contentJSON []byte

	err = dbPool.QueryRow(ctx, query, lessonID, courseID).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&contentJSON,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Lesson not found in this course",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка получения урока: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to get lesson",
			Code:  "INTERNAL_ERROR",
		})
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
	return c.JSON(lesson)
}

// PUT /courses/:course_id/lessons/:lesson_id
func updateLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id FROM knowledge_base.lesson_d WHERE id = $1 AND course_id = $2`
	var existingID string
	err := dbPool.QueryRow(ctx, checkQuery, lessonID, courseID).Scan(&existingID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Lesson not found in this course",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки урока: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check lesson",
			Code:  "INTERNAL_ERROR",
		})
	}

	var input LessonUpdate
	if err := c.BodyParser(&input); err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid request body",
			Code:  "BAD_REQUEST",
		})
	}

	if input.Title != "" && len(input.Title) > 255 {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Title must be less than 255 characters",
			Code:  "VALIDATION_ERROR",
		})
	}

	var contentJSON []byte
	var updateContent bool
	if input.Content != nil {
		var err error
		contentJSON, err = json.Marshal(input.Content)
		if err != nil {
			log.Printf("❌ Ошибка маршалинга content: %v", err)
			return c.Status(400).JSON(ErrorResponse{
				Error: "Invalid content format",
				Code:  "BAD_REQUEST",
			})
		}
		updateContent = true
	}

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

	err = dbPool.QueryRow(ctx, updateQuery, params...).Scan(
		&lesson.ID,
		&lesson.Title,
		&lesson.CourseID,
		&contentBytes,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		log.Printf("❌ Ошибка обновления урока: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to update lesson",
			Code:  "INTERNAL_ERROR",
		})
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
	return c.JSON(lesson)
}

// DELETE /courses/:course_id/lessons/:lesson_id
func deleteLesson(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	lessonID := c.Params("lesson_id")

	if !isValidUUID(courseID) || !isValidUUID(lessonID) {
		return c.Status(400).JSON(ErrorResponse{
			Error: "Invalid ID format",
			Code:  "BAD_REQUEST",
		})
	}

	ctx := context.Background()
	checkQuery := `SELECT id, title FROM knowledge_base.lesson_d WHERE id = $1 AND course_id = $2`
	var lessonTitle string
	err := dbPool.QueryRow(ctx, checkQuery, lessonID, courseID).Scan(&lessonID, &lessonTitle)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return c.Status(404).JSON(ErrorResponse{
				Error: "Lesson not found in this course",
				Code:  "NOT_FOUND",
			})
		}
		log.Printf("❌ Ошибка проверки урока: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to check lesson",
			Code:  "INTERNAL_ERROR",
		})
	}

	deleteQuery := `DELETE FROM knowledge_base.lesson_d WHERE id = $1 AND course_id = $2`
	result, err := dbPool.Exec(ctx, deleteQuery, lessonID, courseID)
	if err != nil {
		log.Printf("❌ Ошибка удаления урока: %v", err)
		return c.Status(500).JSON(ErrorResponse{
			Error: "Failed to delete lesson",
			Code:  "INTERNAL_ERROR",
		})
	}

	if result.RowsAffected() == 0 {
		return c.Status(404).JSON(ErrorResponse{
			Error: "Lesson not found",
			Code:  "NOT_FOUND",
		})
	}

	log.Printf("🗑️ Удален урок: %s", lessonTitle)
	return c.SendStatus(204)
}

// ============ ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ============

func isValidLevel(level string) bool {
	switch strings.ToLower(level) {
	case "hard", "medium", "easy":
		return true
	default:
		return false
	}
}

func isValidVisibility(visibility string) bool {
	switch strings.ToLower(visibility) {
	case "draft", "public", "private":
		return true
	default:
		return false
	}
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}