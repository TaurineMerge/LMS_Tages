"""Health check endpoints."""
from fastapi import APIRouter

from app.database import fetch_one

router = APIRouter(tags=["Health"])


@router.get(
    "/health",
    summary="Health check",
    description="Проверка работоспособности сервиса"
)
async def health_check():
    """Basic health check endpoint."""
    return {"status": "healthy"}


@router.get(
    "/health/db",
    summary="Database health check",
    description="Проверка подключения к базе данных"
)
async def db_health_check():
    """Database connectivity health check."""
    try:
        result = await fetch_one("SELECT 1 as ok")
        if result and result.get("ok") == 1:
            return {"status": "healthy", "database": "connected"}
        return {"status": "unhealthy", "database": "query failed"}
    except Exception as e:
        return {"status": "unhealthy", "database": "disconnected", "error": str(e)}
