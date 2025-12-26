// Package tracing инкапсулирует настройку и управление OpenTelemetry.
package tracing

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

// Tracer инкапсулирует провайдер трассировки OpenTelemetry для управления его жизненным циклом.
type Tracer struct {
	traceProvider *sdktrace.TracerProvider
}

// New инициализирует и настраивает глобальный провайдер трассировки OpenTelemetry.
// Он создает ресурс, настраивает экспортер OTLP/gRPC и устанавливает
// глобальный провайдер трассировки и пропагатор.
func New(cfg *config.OtelConfig) (*Tracer, error) {
	ctx := context.Background()

	// Создание ресурса OTel, который описывает сервис.
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// Настройка экспортера трассировок через OTLP/gRPC.
	traceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(), // Использование небезопасного соединения.
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()), // Блокировка до установления соединения.
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Создание провайдера трассировки с батчером и настроенным ресурсом.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	// Установка глобального провайдера трассировки и пропагатора контекста.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &Tracer{
		traceProvider: tp,
	}, nil
}

// Close корректно завершает работу провайдера трассировки, обеспечивая отправку всех
// оставшихся в буфере трассировок. Должен вызываться при завершении работы приложения.
func (t *Tracer) Close() {
	if err := t.traceProvider.Shutdown(context.Background()); err != nil {
		slog.Error("Error shutting down tracer provider", "error", err)
	}
}
