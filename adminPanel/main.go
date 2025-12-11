// main.go
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

// Простая OTEL middleware для Fiber
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

		// Логируем все атрибуты запроса
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

		// Логируем все заголовки запроса
		for k, v := range c.GetReqHeaders() {
			if len(v) > 0 {
				attrs = append(attrs, attribute.String("http.request.header."+k, v[0]))
			}
		}

		// Логируем тело запроса
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

		// Логируем время выполнения
		duration := time.Since(startTime)
		span.SetAttributes(attribute.Float64("http.request.duration_ms", float64(duration.Milliseconds())))

		// Логируем тело ответа
		responseBody := c.Response().Body()
		if len(responseBody) > 0 {
			const maxLoggedResponseBody = 2048
			if len(responseBody) > maxLoggedResponseBody {
				responseBody = responseBody[:maxLoggedResponseBody]
			}
			span.AddEvent("http.response.body", trace.WithAttributes(attribute.String("body", string(responseBody))))
		}

		// Логируем все заголовки ответа
		c.Response().Header.VisitAll(func(key, value []byte) {
			span.AddEvent("http.response.header."+string(key), trace.WithAttributes(attribute.String("value", string(value))))
		})

		// Логируем дополнительную информацию об ответе
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

		// Проставляем trace-id/span-id в ответ
		sc := span.SpanContext()
		if sc.HasTraceID() {
			c.Set("Trace-Id", sc.TraceID().String())
		}
		if sc.HasSpanID() {
			c.Set("Span-Id", sc.SpanID().String())
		}

		// Логируем ошибки
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

func main() {
	ctx := context.Background()
	// Инициализация конфигурации
	settings := config.NewSettings()

	// Инициализация аутентификации
	if err := middleware.InitAuth(); err != nil {
		log.Fatalf("⚠️  Failed to initialize auth: %v", err)
	}

	// Инициализация базы данных
	db, err := database.InitDB(settings)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Init tracing
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

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		AppName:               "Admin Panel API",
		DisableStartupMessage: false,
	})

	// Глобальные middleware
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

	// Общий обработчик ошибок
	app.Use(middleware.ErrorHandlerMiddleware())

	// Публичные маршруты (без префикса /admin)
	healthHandler := handlers.NewHealthHandler(db)
	app.Get("/health", healthHandler.HealthCheck)
	app.Get("/health/db", healthHandler.DBHealthCheck)

	app.Static("/doc", "./docs")

	// Затем Swagger UI
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

	// API маршруты с префиксом /admin/api/v1
	api := app.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	// Инициализация репозиториев
	categoryRepo := repositories.NewCategoryRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)

	// Инициализация сервисов
	categoryService := services.NewCategoryService(categoryRepo)
	courseService := services.NewCourseService(courseRepo, categoryRepo)
	lessonService := services.NewLessonService(lessonRepo, courseRepo)

	// Инициализация обработчиков
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	courseHandler := handlers.NewCourseHandler(courseService)
	lessonHandler := handlers.NewLessonHandler(lessonService)

	// Регистрация маршрутов
	categoryHandler.RegisterRoutes(api)
	courseHandler.RegisterRoutes(api)
	lessonHandler.RegisterRoutes(api)

	// Запуск сервера
	log.Printf("🚀 Server starting on %s", settings.APIAddress)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/", settings.APIAddress)
	log.Printf("📖 Swagger JSON: http://localhost%s/swagger/doc.json", settings.APIAddress)
	log.Printf("🏥 Health check: http://localhost%s/health", settings.APIAddress)

	if err := app.Listen(settings.APIAddress); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
