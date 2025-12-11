"""Personal Account API - Education Platform."""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.openapi.utils import get_openapi
from fastapi.responses import JSONResponse
from fastapi.staticfiles import StaticFiles
from opentelemetry import trace

from app.config import get_settings
from app.database import close_db_pool, init_db_pool
from app.exceptions import app_exception
from app.routers import auth, certificates, health, pages, students, visits
from app.telemetry_config import (
    configure_telemetry,
    instrument_fastapi,
    instrument_httpx,
    instrument_logging,
    shutdown_telemetry,
)


# Custom formatter that handles missing trace context
class TraceContextFormatter(logging.Formatter):
    """Formatter that safely handles missing OTEL trace context."""

    def format(self, record):
        # Add trace_id and span_id if not present (for logs outside trace context)
        if not hasattr(record, "otelTraceID"):
            record.otelTraceID = "0" * 32
        if not hasattr(record, "otelSpanID"):
            record.otelSpanID = "0" * 16

        # Use the shorter names for compatibility
        record.trace_id = getattr(record, "otelTraceID", "0" * 32)
        record.span_id = getattr(record, "otelSpanID", "0" * 16)

        return super().format(record)


# Configure logging with custom formatter
handler = logging.StreamHandler()
handler.setFormatter(
    TraceContextFormatter(
        fmt="%(asctime)s - %(name)s - %(levelname)s - [trace_id=%(trace_id)s span_id=%(span_id)s] - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
)

logging.basicConfig(level=logging.INFO, handlers=[handler])
logger = logging.getLogger(__name__)

settings = get_settings()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    # Startup
    logger.info("Starting up Personal Account API...")

    # Initialize telemetry first
    try:
        configure_telemetry()
        instrument_logging()
        instrument_httpx()
        logger.info("Telemetry initialized successfully")
    except Exception as exc:
        logger.error("Failed to initialize telemetry: %s", exc, exc_info=True)

    # Initialize database
    await init_db_pool()
    logger.info("Database pool initialized")

    yield

    # Shutdown
    logger.info("Shutting down Personal Account API...")
    await close_db_pool()
    logger.info("Database pool closed")

    # Shutdown telemetry (flush pending spans)
    shutdown_telemetry()


# OAuth2 security scheme for Swagger
oauth2_scheme = {
    "type": "oauth2",
    "flows": {
        "password": {
            "tokenUrl": f"{settings.KEYCLOAK_PUBLIC_URL}/realms/{settings.KEYCLOAK_REALM}/protocol/openid-connect/token",
            "scopes": {},
        }
    },
}

app = FastAPI(
    title="Personal Account API",
    description="""
API для личного кабинета системы онлайн образования

## Авторизация

### Способ 1: OAuth2 Password Flow (логин/пароль)
1. Нажмите кнопку **Authorize** 
2. Введите client_id: `student-client`
3. Введите логин и пароль:
   - Username: `student` Password: `student`

### Способ 2: Bearer Token  
1. Получите токен через POST `/api/v1/auth/token`
2. Скопируйте `access_token` из ответа
3. Нажмите **Authorize** → **BearerAuth** → вставьте токен

**Токены действительны 5 минут**
    """,
    version="1.0.0",
    lifespan=lifespan,
    root_path="/account",  # Базовый путь приложения за nginx
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
    swagger_ui_init_oauth={
        "clientId": settings.KEYCLOAK_CLIENT_ID,
        "clientSecret": settings.KEYCLOAK_CLIENT_SECRET,
        "usePkceWithAuthorizationCodeGrant": False,
    },
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
        "OAuth2PasswordBearer": oauth2_scheme,
        "BearerAuth": {
            "type": "http",
            "scheme": "bearer",
            "bearerFormat": "JWT",
            "description": "Введите JWT access_token (без префикса 'Bearer')",
        },
    }
    # Global security - endpoints can override
    openapi_schema["security"] = [{"OAuth2PasswordBearer": []}, {"BearerAuth": []}]
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
    # Record exception in span
    span = trace.get_current_span()
    if span and span.is_recording():
        span.record_exception(exc)
        span.set_attribute("error", True)
        span.set_attribute("error.type", "app_exception")
        span.set_attribute("error.status_code", exc.status_code)

    return JSONResponse(status_code=exc.status_code, content={"error": exc.message})


@app.exception_handler(Exception)
async def general_exception_handler(request: Request, exc: Exception):
    """Handle unexpected exceptions."""
    # Record exception in span
    span = trace.get_current_span()
    if span and span.is_recording():
        span.record_exception(exc)
        span.set_attribute("error", True)
        span.set_attribute("error.type", "unexpected_exception")

    logger.exception(f"Unexpected error: {exc}")
    return JSONResponse(status_code=500, content={"error": "Internal server error"})


# Include routers
app.include_router(pages.router)  # Frontend pages (no prefix)
app.include_router(health.router)
app.include_router(auth.router, prefix=settings.API_PREFIX)
app.include_router(students.router, prefix=settings.API_PREFIX)
app.include_router(certificates.router, prefix=settings.API_PREFIX)
app.include_router(visits.router, prefix=settings.API_PREFIX)

# Mount static files
app.mount("/static", StaticFiles(directory="static"), name="static")

# Instrument FastAPI with telemetry
instrument_fastapi(app)
