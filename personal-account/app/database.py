"""Database connection pool management using psycopg (v3)."""
import logging
from contextlib import asynccontextmanager
from typing import AsyncGenerator, Any

from psycopg_pool import AsyncConnectionPool
from psycopg.rows import dict_row

from app.config import get_settings

logger = logging.getLogger(__name__)

# Global connection pool
_pool: AsyncConnectionPool | None = None


async def init_db_pool() -> None:
    """Initialize the database connection pool."""
    global _pool
    
    settings = get_settings()
    
    _pool = AsyncConnectionPool(
        conninfo=settings.database_url,
        min_size=settings.DATABASE_POOL_MIN_SIZE,
        max_size=settings.DATABASE_POOL_MAX_SIZE,
        open=False,
    )
    await _pool.open()
    logger.info("Database connection pool initialized")


async def close_db_pool() -> None:
    """Close the database connection pool."""
    global _pool
    
    if _pool:
        await _pool.close()
        _pool = None
        logger.info("Database connection pool closed")


@asynccontextmanager
async def get_connection() -> AsyncGenerator:
    """Get a connection from the pool."""
    if _pool is None:
        raise RuntimeError("Database pool is not initialized")
    
    async with _pool.connection() as conn:
        yield conn


@asynccontextmanager
async def get_cursor() -> AsyncGenerator:
    """Get a cursor with dict row factory."""
    async with get_connection() as conn:
        async with conn.cursor(row_factory=dict_row) as cur:
            yield cur


async def fetch_one(query: str, params: tuple | dict | None = None) -> dict[str, Any] | None:
    """Execute query and fetch one row."""
    async with get_cursor() as cur:
        await cur.execute(query, params)
        return await cur.fetchone()


async def fetch_all(query: str, params: tuple | dict | None = None) -> list[dict[str, Any]]:
    """Execute query and fetch all rows."""
    async with get_cursor() as cur:
        await cur.execute(query, params)
        return await cur.fetchall()


async def execute(query: str, params: tuple | dict | None = None) -> int:
    """Execute query and return affected rows count."""
    async with get_connection() as conn:
        async with conn.cursor() as cur:
            await cur.execute(query, params)
            await conn.commit()
            return cur.rowcount


async def execute_returning(query: str, params: tuple | dict | None = None) -> dict[str, Any] | None:
    """Execute query with RETURNING clause and fetch the result."""
    async with get_connection() as conn:
        async with conn.cursor(row_factory=dict_row) as cur:
            await cur.execute(query, params)
            result = await cur.fetchone()
            await conn.commit()
            return result
