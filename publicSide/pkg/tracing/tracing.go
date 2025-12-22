// Package tracing provides utilities for initializing OpenTelemetry tracing.
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

type Tracer struct {
	traceProvider *sdktrace.TracerProvider
}

func New(cfg *config.OtelConfig) (*Tracer, error) {
	ctx := context.Background()

	// Create a new resource with service name and version attributes.
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// Create an OTLP gRPC exporter.
	// This assumes the collector is running on the specified endpoint and is accessible.
	// For local development, an insecure connection is typically used.
	traceExporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.CollectorEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create a new tracer provider with the resource and the Jaeger exporter.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	// Set the global tracer provider and propagator.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &Tracer{
		traceProvider: tp,
	}, nil
}

func (t *Tracer) Close() {
	if err := t.traceProvider.Shutdown(context.Background()); err != nil {
		slog.Error("Error shutting down tracer provider", "error", err)
	}
}
