"""Production-ready OpenTelemetry configuration.

Features:
- Configurable sampling strategies
- Span attribute limits
- Export batching and timeouts
- Resource detection
- Environment-based settings
"""
import logging
from typing import Optional
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider, SpanLimits
from opentelemetry.sdk.trace.export import (
    BatchSpanProcessor,
    ConsoleSpanExporter,
    SpanExporter,
)
from opentelemetry.sdk.trace.sampling import (
    TraceIdRatioBased,
    ParentBasedTraceIdRatio,
    ALWAYS_ON,
    ALWAYS_OFF,
)
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import (
    Resource,
    SERVICE_NAME,
    SERVICE_VERSION,
    DEPLOYMENT_ENVIRONMENT,
    SERVICE_NAMESPACE,
)
from opentelemetry.instrumentation.logging import LoggingInstrumentor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.httpx import HTTPXClientInstrumentor
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor

from app.config import get_settings

logger = logging.getLogger(__name__)
settings = get_settings()


class TelemetryConfig:
    """Centralized OpenTelemetry configuration."""
    
    # Sampling configuration
    SAMPLING_RATE = float(settings.OTEL_SAMPLING_RATE or "1.0")  # 1.0 = 100%, 0.1 = 10%
    
    # Span limits (prevent memory issues with huge spans)
    MAX_ATTRIBUTES_PER_SPAN = int(settings.OTEL_SPAN_ATTRIBUTE_COUNT_LIMIT or "128")
    MAX_EVENTS_PER_SPAN = int(settings.OTEL_SPAN_EVENT_COUNT_LIMIT or "128")
    MAX_LINKS_PER_SPAN = int(settings.OTEL_SPAN_LINK_COUNT_LIMIT or "128")
    MAX_ATTRIBUTE_LENGTH = int(settings.OTEL_ATTRIBUTE_VALUE_LENGTH_LIMIT or "4096")
    
    # Export configuration
    EXPORT_TIMEOUT_MILLIS = int(settings.OTEL_EXPORTER_OTLP_TIMEOUT or "30000")  # 30s
    EXPORT_MAX_BATCH_SIZE = int(settings.OTEL_BSP_MAX_EXPORT_BATCH_SIZE or "512")
    EXPORT_SCHEDULE_DELAY_MILLIS = int(settings.OTEL_BSP_SCHEDULE_DELAY or "5000")  # 5s
    EXPORT_MAX_QUEUE_SIZE = int(settings.OTEL_BSP_MAX_QUEUE_SIZE or "2048")
    
    # Feature flags
    ENABLE_CONSOLE_EXPORTER = settings.OTEL_ENABLE_CONSOLE_EXPORTER or False
    ENABLE_LOGGING_INSTRUMENTATION = settings.OTEL_ENABLE_LOGGING or True
    ENABLE_SQLALCHEMY_INSTRUMENTATION = settings.OTEL_ENABLE_SQLALCHEMY or True
    ENABLE_HTTPX_INSTRUMENTATION = settings.OTEL_ENABLE_HTTPX or True
    
    @staticmethod
    def get_sampler():
        """Get configured sampler based on sampling rate.
        
        Uses ParentBasedTraceIdRatio:
        - If parent span is sampled, child is sampled
        - Otherwise, uses probability-based sampling
        """
        if TelemetryConfig.SAMPLING_RATE >= 1.0:
            return ALWAYS_ON
        elif TelemetryConfig.SAMPLING_RATE <= 0.0:
            return ALWAYS_OFF
        else:
            return ParentBasedTraceIdRatio(TelemetryConfig.SAMPLING_RATE)
    
    @staticmethod
    def get_span_limits():
        """Get configured span limits."""
        return SpanLimits(
            max_attributes=TelemetryConfig.MAX_ATTRIBUTES_PER_SPAN,
            max_events=TelemetryConfig.MAX_EVENTS_PER_SPAN,
            max_links=TelemetryConfig.MAX_LINKS_PER_SPAN,
            max_attribute_length=TelemetryConfig.MAX_ATTRIBUTE_LENGTH,
        )
    
    @staticmethod
    def get_resource():
        """Get service resource with metadata."""
        return Resource.create({
            SERVICE_NAME: settings.OTEL_SERVICE_NAME,
            SERVICE_VERSION: settings.SERVICE_VERSION or "1.0.0",
            DEPLOYMENT_ENVIRONMENT: settings.ENVIRONMENT or "production",
            SERVICE_NAMESPACE: settings.SERVICE_NAMESPACE or "lms-tages",
            # Custom attributes
            "service.instance.id": settings.SERVICE_INSTANCE_ID or "unknown",
            "service.team": "backend",
        })
    
    @staticmethod
    def get_span_exporter() -> SpanExporter:
        """Get configured span exporter (OTLP)."""
        return OTLPSpanExporter(
            endpoint=settings.OTEL_EXPORTER_OTLP_ENDPOINT,
            insecure=settings.OTEL_EXPORTER_OTLP_INSECURE,
            timeout=TelemetryConfig.EXPORT_TIMEOUT_MILLIS // 1000,  # Convert to seconds
        )
    
    @staticmethod
    def get_batch_span_processor(exporter: SpanExporter) -> BatchSpanProcessor:
        """Get configured batch span processor."""
        return BatchSpanProcessor(
            exporter,
            max_queue_size=TelemetryConfig.EXPORT_MAX_QUEUE_SIZE,
            schedule_delay_millis=TelemetryConfig.EXPORT_SCHEDULE_DELAY_MILLIS,
            max_export_batch_size=TelemetryConfig.EXPORT_MAX_BATCH_SIZE,
            export_timeout_millis=TelemetryConfig.EXPORT_TIMEOUT_MILLIS,
        )


def configure_telemetry() -> None:
    """Initialize OpenTelemetry with production-ready configuration.
    
    This should be called once at application startup.
    """
    try:
        # Create resource
        resource = TelemetryConfig.get_resource()
        
        # Create tracer provider
        provider = TracerProvider(
            resource=resource,
            sampler=TelemetryConfig.get_sampler(),
            span_limits=TelemetryConfig.get_span_limits(),
        )
        
        # Add OTLP exporter
        otlp_exporter = TelemetryConfig.get_span_exporter()
        provider.add_span_processor(
            TelemetryConfig.get_batch_span_processor(otlp_exporter)
        )
        
        # Optionally add console exporter for debugging
        if TelemetryConfig.ENABLE_CONSOLE_EXPORTER:
            console_exporter = ConsoleSpanExporter()
            provider.add_span_processor(BatchSpanProcessor(console_exporter))
            logger.info("Console span exporter enabled (debug mode)")
        
        # Set global tracer provider
        trace.set_tracer_provider(provider)
        
        logger.info(
            "OpenTelemetry configured: service=%s, sampling_rate=%.2f, endpoint=%s",
            settings.OTEL_SERVICE_NAME,
            TelemetryConfig.SAMPLING_RATE,
            settings.OTEL_EXPORTER_OTLP_ENDPOINT,
        )
        
    except Exception as exc:
        logger.error("Failed to configure OpenTelemetry: %s", exc, exc_info=True)
        raise


def instrument_logging() -> None:
    """Instrument Python logging to attach trace context."""
    if not TelemetryConfig.ENABLE_LOGGING_INSTRUMENTATION:
        logger.info("Logging instrumentation disabled")
        return
    
    try:
        LoggingInstrumentor().instrument(
            set_logging_format=False,  # Keep custom format
            log_level=logging.INFO,
        )
        logger.info("Logging instrumentation enabled")
    except Exception as exc:
        logger.warning("Failed to instrument logging: %s", exc)


def instrument_fastapi(app) -> None:
    """Instrument FastAPI application."""
    try:
        FastAPIInstrumentor.instrument_app(
            app,
            tracer_provider=trace.get_tracer_provider(),
            excluded_urls=settings.OTEL_EXCLUDED_URLS or "",
        )
        logger.info("FastAPI instrumentation enabled")
    except Exception as exc:
        logger.warning("Failed to instrument FastAPI: %s", exc)


def instrument_httpx() -> None:
    """Instrument httpx for outbound HTTP request tracing."""
    if not TelemetryConfig.ENABLE_HTTPX_INSTRUMENTATION:
        logger.info("httpx instrumentation disabled")
        return
    
    try:
        HTTPXClientInstrumentor().instrument()
        logger.info("httpx instrumentation enabled")
    except Exception as exc:
        logger.warning("Failed to instrument httpx: %s", exc)


def instrument_sqlalchemy(engine=None) -> None:
    """Instrument SQLAlchemy for database query tracing.
    
    Args:
        engine: SQLAlchemy engine instance (optional, can auto-detect)
    """
    if not TelemetryConfig.ENABLE_SQLALCHEMY_INSTRUMENTATION:
        logger.info("SQLAlchemy instrumentation disabled")
        return
    
    try:
        if engine:
            SQLAlchemyInstrumentor().instrument(
                engine=engine,
                tracer_provider=trace.get_tracer_provider(),
            )
        else:
            SQLAlchemyInstrumentor().instrument(
                tracer_provider=trace.get_tracer_provider(),
            )
        logger.info("SQLAlchemy instrumentation enabled")
    except Exception as exc:
        logger.warning("Failed to instrument SQLAlchemy: %s", exc)


def shutdown_telemetry() -> None:
    """Gracefully shutdown telemetry providers.
    
    Call this on application shutdown to flush pending spans.
    """
    try:
        provider = trace.get_tracer_provider()
        if hasattr(provider, 'shutdown'):
            provider.shutdown()
            logger.info("Telemetry shutdown complete")
    except Exception as exc:
        logger.error("Error during telemetry shutdown: %s", exc)
