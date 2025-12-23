// Пакет main - точка входа для Admin Panel API
//
// Admin Panel - это веб-сервис для управления учебным контентом,
// предоставляющий REST API для работы с категориями, курсами и уроками.
//
// Сервис включает в себя:
//   - Аутентификацию через JWT-токены
//   - OpenTelemetry для трассировки запросов
//   - Swagger UI для документации API
//   - Middleware для CORS, логирования и обработки ошибок
//
// Пример использования:
//
//	# Запуск сервера
//	go run main.go
//
//	# Доступ к Swagger UI
//	http://localhost:4000/swagger/
//
//	# Health check
//	http://localhost:4000/health
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"adminPanel/config"
	"adminPanel/database"
	"adminPanel/handlers"
	webhandlers "adminPanel/handlers/web"
	"adminPanel/middleware"
	"adminPanel/repositories"
	"adminPanel/services"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/handlebars/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// setupTracerProvider инициализирует провайдер трассировки OpenTelemetry
//
// Функция настраивает экспорт трасс в OTLP-коллектор и возвращает
// TracerProvider для использования в приложении.
//
// Параметры:
//   - ctx: контекст выполнения
//   - cfg: конфигурация OpenTelemetry
//
// Возвращает:
//   - TracerProvider: провайдер для создания трасс
//   - error: ошибка инициализации (если есть)
func setupTracerProvider(ctx context.Context, cfg config.OTelConfig) (*tracesdk.TracerProvider, error) {
	if !cfg.Enabled {
		log.Println("ℹ️  OpenTelemetry tracing is disabled (OTEL_EXPORTER_OTLP_ENDPOINT not set)")
		return nil, nil
	}

	target := strings.TrimPrefix(strings.TrimPrefix(cfg.Endpoint, "http://"), "https://")
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(target),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Printf("✅ OpenTelemetry tracing initialized (endpoint=%s, service=%s)", cfg.Endpoint, cfg.ServiceName)
	return tp, nil
}

// tracingMiddleware создает middleware для трассировки HTTP-запросов
//
// Middleware добавляет в каждый запрос трассу с детальной информацией:
//   - Метод и путь запроса
//   - Заголовки и тело запроса/ответа
//   - Время выполнения
//   - Коды ответов и ошибки
//
// Параметры:
//   - tracer: Tracer для создания спанов
//
// Возвращает:
//   - fiber.Handler: middleware для использования в Fiber
func tracingMiddleware(tracer trace.Tracer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		carrier := propagation.HeaderCarrier{}
		for k, v := range c.GetReqHeaders() {
			if len(v) > 0 {
				carrier.Set(k, v[0])
			}
		}
		ctx := otel.GetTextMapPropagator().Extract(c.Context(), carrier)
		spanName := fmt.Sprintf("%s %s", c.Method(), c.Path())
		ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()
		c.SetUserContext(ctx)

		route := c.Route()
		status := c.Response().StatusCode()
		attrs := []attribute.KeyValue{
			semconv.HTTPMethodKey.String(c.Method()),
			semconv.HTTPRouteKey.String(route.Path),
			semconv.HTTPTargetKey.String(c.OriginalURL()),
			semconv.HTTPStatusCodeKey.Int(status),
			semconv.NetHostNameKey.String(c.Hostname()),
			semconv.HTTPUserAgentKey.String(c.Get("User-Agent")),
			attribute.String("http.request.start_time", startTime.Format(time.RFC3339)),
		}
		if ip := c.IP(); ip != "" {
			attrs = append(attrs, attribute.String("net.peer.ip", ip))
		}
		if q := c.Context().QueryArgs().String(); q != "" {
			attrs = append(attrs, attribute.String("http.query", q))
		}

		for k, v := range c.GetReqHeaders() {
			if len(v) > 0 {
				attrs = append(attrs, attribute.String("http.request.header."+k, v[0]))
			}
		}

		if len(c.Body()) > 0 {
			body := c.Body()
			const maxLoggedBody = 2048
			if len(body) > maxLoggedBody {
				body = body[:maxLoggedBody]
			}
			attrs = append(attrs, attribute.String("http.request.body", string(body)))
		}

		span.SetAttributes(attrs...)

		err := c.Next()

		duration := time.Since(startTime)
		span.SetAttributes(attribute.Float64("http.request.duration_ms", float64(duration.Milliseconds())))

		responseBody := c.Response().Body()
		if len(responseBody) > 0 {
			const maxLoggedResponseBody = 2048
			if len(responseBody) > maxLoggedResponseBody {
				responseBody = responseBody[:maxLoggedResponseBody]
			}
			span.AddEvent("http.response.body", trace.WithAttributes(attribute.String("body", string(responseBody))))
		}

		c.Response().Header.VisitAll(func(key, value []byte) {
			span.AddEvent("http.response.header."+string(key), trace.WithAttributes(attribute.String("value", string(value))))
		})

		span.SetAttributes(
			attribute.Int("http.response.size", len(responseBody)),
			attribute.String("http.response.time", time.Now().Format(time.RFC3339)),
		)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		if status >= 500 {
			span.SetStatus(codes.Error, http.StatusText(status))
		}

		sc := span.SpanContext()
		if sc.HasTraceID() {
			c.Set("Trace-Id", sc.TraceID().String())
		}
		if sc.HasSpanID() {
			c.Set("Span-Id", sc.SpanID().String())
		}

		if err != nil || status >= 500 {
			log.Printf("trace=%s span=%s method=%s path=%s status=%d err=%v duration=%s",
				sc.TraceID().String(), sc.SpanID().String(), c.Method(), c.Path(), status, err, duration)
		}

		if err != nil {
			return err
		}
		return nil
	}
}

// main - точка входа приложения Admin Panel
//
// Функция выполняет:
//   - Инициализацию конфигурации с валидацией
//   - Настройку аутентификации
//   - Подключение к базе данных
//   - Настройку трассировки
//   - Создание и настройку Fiber приложения
//   - Регистрацию маршрутов
//   - Запуск HTTP-сервера
//
// Используемые компоненты:
//   - Fiber: веб-фреймворк
//   - PostgreSQL: база данных
//   - Keycloak: аутентификация
//   - OpenTelemetry: трассировка
func main() {
	ctx := context.Background()

	// Загружаем конфигурацию
	settings := config.NewSettings()

	// Валидируем обязательные параметры
	if err := settings.Validate(); err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	log.Printf("📋 Configuration loaded (debug=%v)", settings.Debug)

	// Инициализируем аутентификацию
	if err := middleware.InitAuth(); err != nil {
		log.Fatalf("⚠️  Failed to initialize auth: %v", err)
	}

	// Подключаемся к базе данных
	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Настраиваем трассировку
	tp, err := setupTracerProvider(ctx, settings.OTel)
	if err != nil {
		log.Printf("⚠️  Failed to initialize tracing: %v", err)
	} else if tp != nil {
		defer func() {
			if shutdownErr := tp.Shutdown(ctx); shutdownErr != nil {
				log.Printf("⚠️  Failed to shutdown tracer provider: %v", shutdownErr)
			}
		}()
	}

	// Создаём Fiber приложение
	engine := handlebars.New("./templates", ".hbs")

	// Регистрируем хелпер eq для сравнения строк
	engine.AddFunc("eq", func(a, b string) bool {
		return a == b
	})

	app := fiber.New(fiber.Config{
		AppName:               settings.Server.AppName,
		DisableStartupMessage: false,
		Views:                 engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(tracingMiddleware(otel.Tracer(settings.OTel.ServiceName)))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(settings.GetCORSOrigins(), ","),
		AllowMethods:     settings.CORS.AllowMethods,
		AllowHeaders:     settings.CORS.AllowHeaders,
		AllowCredentials: settings.CORS.AllowCredentials,
		ExposeHeaders:    settings.CORS.ExposeHeaders,
	}))

	app.Use(middleware.ErrorHandlerMiddleware())

	// Health endpoints
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	// Documentation
	app.Static("/doc", "./docs")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/doc/swagger.json",
		DeepLinking: true,
		Title:       settings.Server.AppName,
		OAuth: &swagger.OAuthConfig{
			ClientId:     settings.Keycloak.ClientID,
			ClientSecret: settings.Keycloak.ClientSecret,
			AppName:      settings.Keycloak.AppName,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}))

	// Repositories
	categoryRepo := repositories.NewCategoryRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	// Services
	categoryService := services.NewCategoryService(categoryRepo)
	courseService := services.NewCourseService(courseRepo, categoryRepo)
	lessonService := services.NewLessonService(lessonRepo, courseRepo)

	// S3 Service
	s3Service, err := services.NewS3Service(settings.Minio)
	if err != nil {
		log.Fatalf("❌ Failed to initialize S3 service: %v", err)
	}

	// Ensure bucket exists
	if err := s3Service.EnsureBucketExists(ctx); err != nil {
		log.Printf("⚠️  Failed to ensure S3 bucket exists: %v", err)
	} else {
		log.Printf("✅ S3 bucket '%s' is ready", settings.Minio.Bucket)
	}

	// Handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	courseHandler := handlers.NewCourseHandler(courseService)
	lessonHandler := handlers.NewLessonHandler(lessonService)
	uploadHandler := handlers.NewUploadHandler(s3Service)

	// API routes
	api := app.Group("/api/v1")

	// Upload routes (БЕЗ AUTH для удобства загрузки из редактора)
	upload := api.Group("/upload")
	uploadHandler.RegisterRoutes(upload)

	// Protected API routes
	api.Use(middleware.AuthMiddleware())
	categoryHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	lessons := api.Group("/categories/:category_id/courses/:course_id/lessons")
	lessonHandler.RegisterRoutes(lessons)

	// Статические файлы (CSS, изображения и т.д.)
	app.Static("/static", "./static")

	// Web routes (без auth для админки)
	// Serve web pages under root (nginx adds /admin prefix), so routes will be available at /admin/...
	web := app.Group("")

	// Web handlers
	categoryWebHandler := webhandlers.NewCategoryWebHandler(categoryService)
	courseWebHandler := webhandlers.NewCourseWebHandler(courseService, categoryService)
	lessonWebHandler := webhandlers.NewLessonWebHandler(lessonService, courseService, categoryService)

	// Register web routes
	web.Get("/categories", categoryWebHandler.RenderCategoriesEditor)
	web.Get("/categories/new", categoryWebHandler.RenderNewCategoryForm)
	web.Post("/categories/create", categoryWebHandler.CreateCategory)
	web.Get("/categories/:id", categoryWebHandler.RenderEditCategoryForm)
	web.Post("/categories/:id/update", categoryWebHandler.UpdateCategory)
	web.Post("/categories/:id/delete", categoryWebHandler.DeleteCategory)

	// Course web routes
	web.Get("/categories/:category_id/courses", courseWebHandler.RenderCoursesEditor)
	web.Get("/categories/:category_id/courses/new", courseWebHandler.RenderNewCourseForm)
	web.Post("/categories/:category_id/courses/create", courseWebHandler.CreateCourse)
	web.Get("/categories/:category_id/courses/:course_id", courseWebHandler.RenderEditCourseForm)
	web.Post("/categories/:category_id/courses/:course_id/update", courseWebHandler.UpdateCourse)
	web.Post("/categories/:category_id/courses/:course_id/delete", courseWebHandler.DeleteCourse)

	// Lesson web routes
	web.Get("/categories/:category_id/courses/:course_id/lessons", lessonWebHandler.RenderLessonsEditor)
	web.Get("/categories/:category_id/courses/:course_id/lessons/new", lessonWebHandler.RenderNewLessonForm)
	web.Post("/categories/:category_id/courses/:course_id/lessons/create", lessonWebHandler.CreateLesson)
	web.Get("/categories/:category_id/courses/:course_id/lessons/:lesson_id", lessonWebHandler.RenderEditLessonForm)
	web.Post("/categories/:category_id/courses/:course_id/lessons/:lesson_id/update", lessonWebHandler.UpdateLesson)
	web.Post("/categories/:category_id/courses/:course_id/lessons/:lesson_id/delete", lessonWebHandler.DeleteLesson)

	// Start server
	log.Printf("🚀 Server starting on %s", settings.Server.Address)
	log.Printf("📚 Swagger UI (via nginx): http://localhost/admin/swagger/")
	log.Printf("📖 Swagger JSON (via nginx): http://localhost/admin/doc/swagger.json")
	log.Printf("🏥 Health check (via nginx): http://localhost/health")
	log.Printf("📍 API (via nginx): http://localhost/admin/api/v1/")

	if err := app.Listen(settings.Server.Address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
