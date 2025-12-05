"""Personal Account API - Education Platform."""
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware

from app.config import get_settings
from app.database import init_db_pool, close_db_pool
from app.exceptions import AppException
from app.routers import students, certificates, visits, health

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


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

app = FastAPI(
    title="Personal Account API",
    description="API для личного кабинета системы онлайн образования",
    version="1.0.0",
    root_path="/account",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # In production, specify actual origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.exception_handler(AppException)
async def app_exception_handler(request: Request, exc: AppException):
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
app.include_router(health.router)
app.include_router(students.router, prefix=settings.API_PREFIX)
app.include_router(certificates.router, prefix=settings.API_PREFIX)
app.include_router(visits.router, prefix=settings.API_PREFIX)


@app.get("/")
async def root():
    """Root endpoint."""
    return {
        "message": "Personal Account API",
        "version": "1.0.0",
        "docs": "/docs"
    }