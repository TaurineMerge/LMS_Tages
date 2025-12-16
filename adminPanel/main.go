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
	"os"
	"strings"
	"time"

	"adminPanel/config"
	"adminPanel/database"
	"adminPanel/handlers"
	"adminPanel/middleware"
	"adminPanel/repositories"
	"adminPanel/services"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
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
//
// Возвращает:
//   - TracerProvider: провайдер для создания трасс
//   - error: ошибка инициализации (если есть)
func setupTracerProvider(ctx context.Context) (*tracesdk.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://otel-collector:4317"
	}
	target := strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(target),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("admin-panel"),
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
//   - Инициализацию конфигурации
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
	settings := config.NewSettings()
	if err := middleware.InitAuth(); err != nil {
		log.Fatalf("⚠️  Failed to initialize auth: %v", err)
	}

	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	tp, err := setupTracerProvider(ctx)
	if err != nil {
		log.Printf("⚠️  Failed to initialize tracing: %v", err)
	} else {
		defer func() {
			if shutdownErr := tp.Shutdown(ctx); shutdownErr != nil {
				log.Printf("⚠️  Failed to shutdown tracer provider: %v", shutdownErr)
			}
		}()
	}

	app := fiber.New(fiber.Config{
		AppName:               "Admin Panel API",
		DisableStartupMessage: false,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(tracingMiddleware(otel.Tracer("admin-panel")))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(settings.GetCORSOrigins(), ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: settings.CORSAllowCredentials,
		ExposeHeaders:    "Content-Length",
	}))

	app.Use(middleware.ErrorHandlerMiddleware())

	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	app.Static("/doc", "./docs")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/doc/swagger.json",
		DeepLinking: true,
		Title:       "Admin Panel API",
		OAuth: &swagger.OAuthConfig{
			ClientId:     settings.ClientID,
			ClientSecret: settings.ClientSecret,
			AppName:      settings.AppName,
			Scopes:       []string{"openid", "profile", "email"},
		},
	}))

	api := app.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	categoryRepo := repositories.NewCategoryRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	categoryService := services.NewCategoryService(categoryRepo)
	courseService := services.NewCourseService(courseRepo, categoryRepo)
	lessonService := services.NewLessonService(lessonRepo, courseRepo)

	categoryHandler := handlers.NewCategoryHandler(categoryService)
	courseHandler := handlers.NewCourseHandler(courseService)
	lessonHandler := handlers.NewLessonHandler(lessonService)

	categoryHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	lessons := api.Group("/categories/:category_id/courses/:course_id/lessons")
	lessonHandler.RegisterRoutes(lessons)

	log.Printf("🚀 Server starting on %s", settings.APIAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/", settings.APIAddress)
	log.Printf("📖 Swagger JSON: http://localhost%s/swagger/doc.json", settings.APIAddress)
	log.Printf("🏥 Health check: http://localhost%s/health", settings.APIAddress)

	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
