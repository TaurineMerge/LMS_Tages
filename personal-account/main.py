"""Personal Account API - Education Platform."""
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.openapi.utils import get_openapi
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.logging import LoggingInstrumentor
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

from app.config import get_settings
from app.database import init_db_pool, close_db_pool
from app.exceptions import app_exception
from app.routers import students, certificates, visits, health, auth, pages

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


_TRACING_CONFIGURED = False
_LOGGING_INSTRUMENTED = False


def configure_tracing() -> None:
    """Configure OpenTelemetry tracing once per process."""
    global _TRACING_CONFIGURED, _LOGGING_INSTRUMENTED
    if _TRACING_CONFIGURED:
        return

    try:
        resource = Resource.create({"service.name": settings.OTEL_SERVICE_NAME})
        provider = TracerProvider(resource=resource)
        exporter = OTLPSpanExporter(
            endpoint=settings.OTEL_EXPORTER_OTLP_ENDPOINT,
            insecure=settings.OTEL_EXPORTER_OTLP_INSECURE,
        )
        provider.add_span_processor(BatchSpanProcessor(exporter))
        trace.set_tracer_provider(provider)
        if not _LOGGING_INSTRUMENTED:
            LoggingInstrumentor().instrument(set_logging_format=False)
            _LOGGING_INSTRUMENTED = True
        _TRACING_CONFIGURED = True
        logger.info(
            "OpenTelemetry tracing configured for %s",
            settings.OTEL_SERVICE_NAME,
        )
    except Exception as exc:  # pragma: no cover - observability optional
        logger.warning("Failed to configure OpenTelemetry tracing: %s", exc)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    # Startup
    logger.info("Starting up Personal Account API...")
    await init_db_pool()
    logger.info("Database pool initialized")
    
    yield
    
    # Shutdown
    logger.info("Shutting down Personal Account API...")
    await close_db_pool()
    logger.info("Database pool closed")


settings = get_settings()
configure_tracing()

# OAuth2 security scheme for Swagger
oauth2_scheme = {
    "type": "oauth2",
    "flows": {
        "authorizationCode": {
            "authorizationUrl": f"{settings.KEYCLOAK_PUBLIC_URL}/realms/{settings.KEYCLOAK_REALM}/protocol/openid-connect/auth",
            "tokenUrl": f"{settings.KEYCLOAK_PUBLIC_URL}/realms/{settings.KEYCLOAK_REALM}/protocol/openid-connect/token",
            "scopes": {
                "openid": "OpenID Connect scope",
                "profile": "User profile",
                "email": "User email"
            }
        }
    }
}

app = FastAPI(
    title="Personal Account API",
    description="API для личного кабинета системы онлайн образования",
    version="1.0.0",
    root_path="/account",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    swagger_ui_init_oauth={
        "clientId": settings.KEYCLOAK_CLIENT_ID,
        "scopes": "openid profile email"
    }
)

# Add security scheme to OpenAPI
def custom_openapi():
    if app.openapi_schema:
        return app.openapi_schema
    openapi_schema = get_openapi(
        title=app.title,
        version=app.version,
        openapi_version=app.openapi_version,
        description=app.description,
        routes=app.routes,
    )
    openapi_schema["components"]["securitySchemes"] = {
        "OAuth2": oauth2_scheme,
        "BearerAuth": {
            "type": "http",
            "scheme": "bearer",
            "bearerFormat": "JWT"
        }
    }
    app.openapi_schema = openapi_schema
    return app.openapi_schema

app.openapi = custom_openapi

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify actual origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.middleware("http")
async def add_trace_headers(request: Request, call_next):
    response = await call_next(request)
    span = trace.get_current_span()
    span_context = span.get_span_context()
    if span_context and span_context.is_valid:
        trace_id = format(span_context.trace_id, "032x")
        span_id = format(span_context.span_id, "016x")
        response.headers["X-Trace-Id"] = trace_id
        response.headers["X-Span-Id"] = span_id
    return response


@app.exception_handler(app_exception)
async def app_exception_handler(request: Request, exc: app_exception):
    """Handle custom application exceptions."""
    return JSONResponse(
        status_code=exc.status_code,
        content={"error": exc.message}
    )


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """Handle unexpected exceptions."""
    logger.exception(f"Unexpected error: {exc}")
    return JSONResponse(
        status_code=500,
        content={"error": "Internal server error"}
    )


# Include routers
app.include_router(pages.router)  # Frontend pages (no prefix)
app.include_router(health.router)
app.include_router(auth.router, prefix=settings.API_PREFIX)
app.include_router(students.router, prefix=settings.API_PREFIX)
app.include_router(certificates.router, prefix=settings.API_PREFIX)
app.include_router(visits.router, prefix=settings.API_PREFIX)

# Mount static files
app.mount("/static", StaticFiles(directory="static"), name="static")

FastAPIInstrumentor.instrument_app(app, tracer_provider=trace.get_tracer_provider())


# Root endpoint removed - pages router handles /