package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/handler"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/repository"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/service"
	"github.com/aymerick/raymond"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// application holds the dependencies for our app.
type application struct {
	templateCache map[string]*raymond.Template
}

// Course represents a single course.
type Course struct {
	Title       string
	Description string
}

// CategoryWithCourses represents a category containing a list of courses.
type CategoryWithCourses struct {
	Name    string
	Courses []Course
}

func main() {
	// Initialize the template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Failed to create template cache: %v", err)
	}

	app := fiber.New()
	
	// Create an app instance
	appl := &application{
		templateCache: templateCache,
	}

	// Setup static file server
	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "./ui/static" // Default for local development
	}
	app.Static("/static", staticPath)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:9090",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// --- HTML Rendering Routes ---
	app.Get("/", appl.home)
	app.Get("/courses", appl.courses)
	
	// --- Existing API Routes ---
	// Serve the `doc` folder so swagger.json is reachable at `/doc/swagger.json`
	app.Static("/doc", "./doc/swagger")

	apiV1 := app.Group("/api/v1")

	// Setup swagger (serve UI at /api/v1/swagger/* and point it to the served JSON)
	apiV1.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/doc/swagger.json",
	}))

	// Initialize repositories
	categoryRepo := repository.NewCategoryMemoryRepository()
	courseRepo := repository.NewCourseMemoryRepository()
	lessonRepo := repository.NewLessonMemoryRepository()

	// Initialize services
	categoryService := service.NewCategoryService(categoryRepo)
	courseService := service.NewCourseService(courseRepo)
	lessonService := service.NewLessonService(lessonRepo)

	// Initialize handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	courseHandler := handler.NewCourseHandler(courseService)
	lessonHandler := handler.NewLessonHandler(lessonService)

	// Register routes
	handler.RegisterRoutes(apiV1, categoryHandler, courseHandler, lessonHandler)

	// Start the server
	log.Println("Starting server on :3000")
	log.Fatal(app.Listen(":3000"))
}

// home is the handler for the home page.
func (app *application) home(c *fiber.Ctx) error {
	data := map[string]interface{}{
		"title":     "Home",
		"pageTitle": "Welcome to Our LMS!",
	}
	c.Type("html")
	return app.render(c, "home.page.hbs", data)
}

// courses is the handler for the courses page.
func (app *application) courses(c *fiber.Ctx) error {
	// Mock data
	mockData := []CategoryWithCourses{
		{
			Name: "Programming",
			Courses: []Course{
				{Title: "Introduction to Go", Description: "Learn the fundamentals of the Go programming language."},
				{Title: "Advanced Go Concurrency", Description: "Master goroutines, channels, and complex concurrency patterns."},
				{Title: "Web Development with Fiber", Description: "Build high-performance web applications in Go."},
				{Title: "Data Structures in Python", Description: "Understand fundamental data structures with Python examples."},
				{Title: "Introduction to Rust", Description: "Learn about Rust's safety and performance features."},
			},
		},
		{
			Name: "Design",
			Courses: []Course{
				{Title: "UI/UX Design Fundamentals", Description: "An introduction to user interface and user experience design."},
				{Title: "Introduction to Figma", Description: "Learn to design and prototype using Figma."},
				{Title: "Color Theory for Designers", Description: "Understand the principles of color in design."},
				{Title: "Typography and Layout", Description: "Master the art of arranging text and elements."},
				{Title: "Design Systems 101", Description: "Learn how to create and manage a design system."},
				{Title: "Motion Design with After Effects", Description: "Bring your designs to life with animation."},
			},
		},
		{
			Name: "Data Science",
			Courses: []Course{
				{Title: "Introduction to Machine Learning", Description: "Learn the basic concepts of machine learning."},
				{Title: "Data Analysis with Pandas", Description: "Master data manipulation and analysis using the Pandas library."},
				{Title: "Deep Learning with TensorFlow", Description: "Build and train neural networks using TensorFlow."},
				{Title: "Natural Language Processing", Description: "Explore techniques for analyzing and understanding text."},
				{Title: "SQL for Data Scientists", Description: "Master SQL for data extraction and analysis."},
			},
		},
		{
			Name: "Business & Marketing",
			Courses: []Course{
				{Title: "Digital Marketing Fundamentals", Description: "Learn the essentials of online marketing."},
				{Title: "SEO Strategy for 2025", Description: "Optimize your web presence for search engines."},
				{Title: "Project Management Essentials", Description: "Understand the lifecycle of a project from start to finish."},
				{Title: "Content Marketing and Strategy", Description: "Create compelling content that drives engagement."},
				{Title: "Social Media Marketing", Description: "Build and engage an audience on social platforms."},
				{Title: "Email Marketing Mastery", Description: "Learn to build and manage effective email campaigns."},
				{Title: "Growth Hacking", Description: "Innovative strategies to grow your user base."},
			},
		},
	}

	data := map[string]interface{}{
		"title":      "Courses",
		"pageTitle":  "Explore Our Courses",
		"categories": mockData,
	}
	c.Type("html")
	return app.render(c, "courses.page.hbs", data)
}


// render executes a template and writes the output to the response.
func (app *application) render(c *fiber.Ctx, name string, data map[string]interface{}) error {
	// Retrieve the page template from the cache
	pageTmpl, ok := app.templateCache[name]
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "page template not found: "+name)
	}

	// Render the page template
	body, err := pageTmpl.Exec(data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to execute page template: "+err.Error())
	}
	
	// Retrieve the layout template from the cache
	layoutTmpl, ok := app.templateCache["base.layout.hbs"]
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "layout template not found: base.layout.hbs")
	}

	// Add the rendered page body to the data for the layout
	data["body"] = body

	// Render the layout template
	finalResult, err := layoutTmpl.Exec(data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to execute layout template: "+err.Error())
	}

	return c.SendString(finalResult)
}


// newTemplateCache parses all .hbs templates and caches them.
func newTemplateCache() (map[string]*raymond.Template, error) {
	cache := map[string]*raymond.Template{}

	templatePath := os.Getenv("TEMPLATE_PATH")
	if templatePath == "" {
		templatePath = "./templates" // Default for local dev
	}
	
	// Register all partials first
	partials, err := filepath.Glob(filepath.Join(templatePath, "partials", "*.hbs"))
	if err != nil {
		return nil, err
	}
	for _, partial := range partials {
		name := strings.TrimSuffix(filepath.Base(partial), ".partial.hbs")
		content, err := os.ReadFile(partial)
		if err != nil {
			return nil, err
		}
		raymond.RegisterPartial(name, string(content))
	}
	
	// Find all template files (layouts and pages)
	allTemplates, err := filepath.Glob(filepath.Join(templatePath, "**", "*.hbs"))
	if err != nil {
		return nil, err
	}

	for _, tmplFile := range allTemplates {
		// We already registered partials, so we can skip them here
		if strings.Contains(tmplFile, "/partials/") {
			continue
		}
		
		name := filepath.Base(tmplFile)
		
		content, err := os.ReadFile(tmplFile)
		if err != nil {
			return nil, err
		}

		// Parse the template
		tmpl, err := raymond.Parse(string(content))
		if err != nil {
			return nil, err
		}
		
		cache[name] = tmpl
	}

	return cache, nil
}
