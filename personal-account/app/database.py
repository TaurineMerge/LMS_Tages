"""Database connection management using SQLAlchemy Core (async)."""

from __future__ import annotations

import json
import logging
from contextlib import asynccontextmanager
from typing import Any, AsyncGenerator, Mapping

from opentelemetry import trace
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor
from opentelemetry.trace import StatusCode
from sqlalchemy import text
from sqlalchemy.engine import Result
from sqlalchemy.ext.asyncio import AsyncConnection, AsyncEngine, create_async_engine

from app.config import get_settings

logger = logging.getLogger(__name__)

_engine: AsyncEngine | None = None
_sqlalchemy_instrumented = False
_tracer = trace.get_tracer(__name__)

# Maximum length for serialized values in span attributes
_MAX_ATTR_LEN = 1024


def _build_async_url(url: str) -> str:
    if url.startswith("postgresql+"):
        return url
    if url.startswith("postgresql://"):
        return url.replace("postgresql://", "postgresql+asyncpg://", 1)
    return url


async def init_db_pool() -> None:
    """Initialize SQLAlchemy async engine."""
    global _engine, _sqlalchemy_instrumented

    settings = get_settings()
    async_url = _build_async_url(settings.database_url)
    _engine = create_async_engine(
        async_url,
        pool_size=settings.DATABASE_POOL_MIN_SIZE,
        max_overflow=max(0, settings.DATABASE_POOL_MAX_SIZE - settings.DATABASE_POOL_MIN_SIZE),
        echo=settings.DEBUG,
    )
    if not _sqlalchemy_instrumented:
        SQLAlchemyInstrumentor().instrument(engine=_engine.sync_engine)
        _sqlalchemy_instrumented = True
    logger.info("SQLAlchemy async engine initialized")


async def close_db_pool() -> None:
    """Dispose SQLAlchemy engine."""
    global _engine
    if _engine:
        await _engine.dispose()
        _engine = None
        logger.info("SQLAlchemy async engine disposed")


@asynccontextmanager
async def get_connection() -> AsyncGenerator[AsyncConnection, None]:
    if _engine is None:
        raise RuntimeError("Database engine is not initialized")
    async with _engine.connect() as conn:
        yield conn


def _as_dict_rows(result: Result) -> list[dict[str, Any]]:
    return [dict(row) for row in result.mappings().all()]


async def fetch_one(query: str, params: Mapping[str, Any] | None = None) -> dict[str, Any] | None:
    stmt = text(query)
    async with get_connection() as conn:
        with _tracer.start_as_current_span("db.fetch_one") as span:
            _record_db_span_attributes(span, query, params)
            try:
                result = await conn.execute(stmt, params or {})
                row = result.mappings().first()
                row_dict = dict(row) if row else None
                _record_result(span, row_dict, single=True)
                span.set_status(StatusCode.OK)
                return row_dict
            except Exception as exc:
                _record_exception(span, exc)
                raise


async def fetch_all(query: str, params: Mapping[str, Any] | None = None) -> list[dict[str, Any]]:
    stmt = text(query)
    async with get_connection() as conn:
        with _tracer.start_as_current_span("db.fetch_all") as span:
            _record_db_span_attributes(span, query, params)
            try:
                result = await conn.execute(stmt, params or {})
                rows = _as_dict_rows(result)
                _record_result(span, rows, single=False)
                span.set_status(StatusCode.OK)
                return rows
            except Exception as exc:
                _record_exception(span, exc)
                raise


async def execute(query: str, params: Mapping[str, Any] | None = None) -> int:
    stmt = text(query)
    if _engine is None:
        raise RuntimeError("Database engine is not initialized")
    async with _engine.begin() as conn:
        with _tracer.start_as_current_span("db.execute") as span:
            _record_db_span_attributes(span, query, params)
            try:
                result = await conn.execute(stmt, params or {})
                rowcount = result.rowcount if result else 0
                span.set_attribute("db.result.rowcount", rowcount)
                span.set_status(StatusCode.OK)
                return rowcount
            except Exception as exc:
                _record_exception(span, exc)
                raise


async def execute_returning(query: str, params: Mapping[str, Any] | None = None) -> dict[str, Any] | None:
    stmt = text(query)
    if _engine is None:
        raise RuntimeError("Database engine is not initialized")
    async with _engine.begin() as conn:
        with _tracer.start_as_current_span("db.execute_returning") as span:
            _record_db_span_attributes(span, query, params)
            try:
                result = await conn.execute(stmt, params or {})
                row = result.mappings().first()
                row_dict = dict(row) if row else None
                _record_result(span, row_dict, single=True)
                span.set_status(StatusCode.OK)
                return row_dict
            except Exception as exc:
                _record_exception(span, exc)
                raise


def _record_db_span_attributes(
    span: trace.Span,
    query: str,
    params: Mapping[str, Any] | None,
) -> None:
    if not span.is_recording():
        return
    span.set_attribute("db.system", "postgresql")
    span.set_attribute("db.statement", _normalize_query(query))
    if params:
        if isinstance(params, Mapping):
            span.set_attribute(
                "db.params.keys",
                ",".join(sorted(str(key) for key in params.keys())),
            )
            # Serialize param values (limited for security)
            param_preview = _safe_serialize_params(params)
            span.set_attribute("db.params.values", param_preview)
        else:
            span.set_attribute("db.params.count", len(params))


def _safe_serialize_params(params: Mapping[str, Any]) -> str:
    """Serialize query parameters for span attributes."""
    try:
        # Limit to first 10 params
        preview = {k: _safe_value(v) for i, (k, v) in enumerate(params.items()) if i < 10}
        s = json.dumps(preview, default=str)
        return s[:_MAX_ATTR_LEN] + ("..." if len(s) > _MAX_ATTR_LEN else "")
    except Exception:
        return "<unserializable>"


def _safe_value(v: Any) -> Any:
    """Mask sensitive values and convert to serializable form."""
    if v is None:
        return None
    if isinstance(v, (str, int, float, bool)):
        # Potentially mask sensitive fields if needed
        return v
    return str(v)


def _record_result(span: trace.Span, result: Any, single: bool) -> None:
    """Record query result details in span."""
    if not span.is_recording():
        return
    if single:
        if result is None:
            span.set_attribute("db.result.found", False)
            span.set_attribute("db.result.count", 0)
        else:
            span.set_attribute("db.result.found", True)
            span.set_attribute("db.result.count", 1)
            # Preview result keys
            if isinstance(result, dict):
                span.set_attribute("db.result.columns", ",".join(result.keys()))
    else:
        count = len(result) if result else 0
        span.set_attribute("db.result.count", count)
        if result and count > 0 and isinstance(result[0], dict):
            span.set_attribute("db.result.columns", ",".join(result[0].keys()))


def _record_exception(span: trace.Span, exc: Exception) -> None:
    """Record exception in span."""
    span.set_status(StatusCode.ERROR, str(exc))
    span.record_exception(exc)
    span.set_attribute("error.type", type(exc).__name__)
    span.set_attribute("error.message", str(exc)[:_MAX_ATTR_LEN])


def _normalize_query(query: str) -> str:
    cleaned = " ".join(query.split())
    return cleaned[:2048]


async def fetch_one_value(query: str, params: Mapping[str, Any] | None = None) -> Any:
    """Fetch a single scalar value from the database."""
    stmt = text(query)
    async with get_connection() as conn:
        with _tracer.start_as_current_span("db.fetch_one_value") as span:
            _record_db_span_attributes(span, query, params)
            try:
                result = await conn.execute(stmt, params or {})
                row = result.first()
                value = row[0] if row else None
                span.set_attribute("db.result.value", str(value)[:_MAX_ATTR_LEN])
                span.set_status(StatusCode.OK)
                return value
            except Exception as exc:
                _record_exception(span, exc)
                raise
