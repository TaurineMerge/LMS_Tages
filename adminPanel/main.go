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

// setupTracerProvider настраивает провайдер трассировки OpenTelemetry.
// Возвращает TracerProvider или nil если трассировка отключена.
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

// tracingMiddleware возвращает промежуточное ПО для трассировки HTTP-запросов.
// Создает span для каждого запроса и записывает метрики.
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

// main является точкой входа в приложение Admin Panel.
// Выполняет следующие шаги:
// 1. Загружает и валидирует конфигурацию из переменных окружения.
// 2. Инициализирует аутентификацию через Keycloak.
// 3. Подключается к базе данных PostgreSQL.
// 4. Настраивает трассировку OpenTelemetry (если включена).
// 5. Создает шаблонизатор Handlebars с вспомогательными функциями.
// 6. Инициализирует Fiber приложение с middleware (recover, logger, tracing, CORS, error handler).
// 7. Настраивает маршруты для health check, Swagger, статических файлов.
// 8. Создает репозитории, сервисы и обработчики для категорий, курсов, уроков и загрузки файлов.
// 9. Регистрирует API маршруты с аутентификацией.
// 10. Регистрирует веб-маршруты для админ-интерфейса.
// 11. Запускает HTTP-сервер на указанном адресе.
func main() {
	ctx := context.Background()

	settings := config.NewSettings()

	if err := settings.Validate(); err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	log.Printf("📋 Configuration loaded (debug=%v)", settings.Debug)

	if err := middleware.InitAuth(); err != nil {
		log.Fatalf("⚠️  Failed to initialize auth: %v", err)
	}

	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

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

	engine := handlebars.New("./templates", ".hbs")

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

	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

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

	categoryRepo := repositories.NewCategoryRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	courseService := services.NewCourseService(courseRepo, categoryRepo)
	lessonService := services.NewLessonService(lessonRepo, courseRepo)

	s3Service, err := services.NewS3Service(settings.Minio)
	if err != nil {
		log.Fatalf("❌ Failed to initialize S3 service: %v", err)
	}

	if err := s3Service.EnsureBucketExists(ctx); err != nil {
		log.Printf("⚠️  Failed to ensure S3 bucket exists: %v", err)
	} else {
		log.Printf("✅ S3 bucket '%s' is ready", settings.Minio.Bucket)
	}

	// Добавляем вспомогательную функцию для генерации URL изображений в шаблонах
	engine.AddFunc("s3ImageURL", func(imageKey string) string {
		if imageKey == "" {
			return ""
		}
		return s3Service.GetImageURL(imageKey)
	})

	categoryHandler := handlers.NewCategoryHandler(categoryService)
	courseHandler := handlers.NewCourseHandler(courseService)
	lessonHandler := handlers.NewLessonHandler(lessonService)
	uploadHandler := handlers.NewUploadHandler(s3Service)

	api := app.Group("/api/v1")

	upload := api.Group("/upload")
	uploadHandler.RegisterRoutes(upload)

	api.Use(middleware.AuthMiddleware())
	categoryHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	lessons := api.Group("/categories/:category_id/courses/:course_id/lessons")
	lessonHandler.RegisterRoutes(lessons)

	app.Static("/static", "./static")

	web := app.Group("")

	categoryWebHandler := webhandlers.NewCategoryWebHandler(categoryService)
	courseWebHandler := webhandlers.NewCourseWebHandler(courseService, categoryService, s3Service, settings.TestModule)
	lessonWebHandler := webhandlers.NewLessonWebHandler(lessonService, courseService, categoryService)
	homeWebHandler := webhandlers.NewHomeWebHandler(categoryService, courseService, lessonService)

	web.Get("/", homeWebHandler.RenderHome)
	web.Get("/categories", categoryWebHandler.RenderCategoriesEditor)
	web.Get("/categories/new", categoryWebHandler.RenderNewCategoryForm)
	web.Post("/categories/create", categoryWebHandler.CreateCategory)
	web.Get("/categories/:id", categoryWebHandler.RenderEditCategoryForm)
	web.Post("/categories/:id/update", categoryWebHandler.UpdateCategory)
	web.Post("/categories/:id/delete", categoryWebHandler.DeleteCategory)

	web.Get("/categories/:category_id/courses", courseWebHandler.RenderCoursesEditor)
	web.Get("/categories/:category_id/courses/new", courseWebHandler.RenderNewCourseForm)
	web.Post("/categories/:category_id/courses/create", courseWebHandler.CreateCourse)
	web.Get("/categories/:category_id/courses/:course_id", courseWebHandler.RenderEditCourseForm)
	web.Post("/categories/:category_id/courses/:course_id/update", courseWebHandler.UpdateCourse)
	web.Post("/categories/:category_id/courses/:course_id/delete", courseWebHandler.DeleteCourse)

	web.Get("/categories/:category_id/courses/:course_id/lessons", lessonWebHandler.RenderLessonsEditor)
	web.Get("/categories/:category_id/courses/:course_id/lessons/new", lessonWebHandler.RenderNewLessonForm)
	web.Post("/categories/:category_id/courses/:course_id/lessons/create", lessonWebHandler.CreateLesson)
	web.Get("/categories/:category_id/courses/:course_id/lessons/:lesson_id", lessonWebHandler.RenderEditLessonForm)
	web.Post("/categories/:category_id/courses/:course_id/lessons/:lesson_id/update", lessonWebHandler.UpdateLesson)
	web.Post("/categories/:category_id/courses/:course_id/lessons/:lesson_id/delete", lessonWebHandler.DeleteLesson)

	log.Printf("🚀 Server starting on %s", settings.Server.Address)
	log.Printf("📚 Swagger UI (via nginx): http://localhost/admin/swagger/")
	log.Printf("📖 Swagger JSON (via nginx): http://localhost/admin/doc/swagger.json")
	log.Printf("🏥 Health check (via nginx): http://localhost/health")
	log.Printf("📍 API (via nginx): http://localhost/admin/api/v1/")

	if err := app.Listen(settings.Server.Address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
