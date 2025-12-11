"""Base repository with common CRUD operations."""

from __future__ import annotations

from typing import Any, Iterable, Mapping
from uuid import UUID

from app.database import execute, fetch_all, fetch_one
from app.db import queries as q
from app.telemetry import traced


class base_repository:
    """Base repository class with common database operations."""

    allowed_order_fields: set[str] = {"created_at", "updated_at", "id"}

    def __init__(
        self,
        table_name: str,
        schema: str = "personal_account",
        default_order_by: str = "created_at",
        orderable_columns: Iterable[str] | None = None,
    ):
        self.table_name = table_name
        self.schema = schema
        self.full_table_name = f"{schema}.{table_name}"
        self.default_order_by = default_order_by
        base_columns = orderable_columns or self.allowed_order_fields
        self.orderable_columns = {column.lower() for column in base_columns}
        if self.default_order_by.lower() not in self.orderable_columns:
            self.orderable_columns.add(self.default_order_by.lower())
        self._table_identifier = self.full_table_name

    @traced()
    async def get_by_id(self, entity_id: UUID) -> dict[str, Any] | None:
        """Get entity by ID."""
        query = q.BASE_SELECT_BY_ID.format(table=self._table_identifier)
        return await fetch_one(query, {"id": entity_id})

    @traced()
    async def get_all(
        self, limit: int = 100, offset: int = 0, order_by: str | None = None, order_dir: str = "DESC"
    ) -> list[dict[str, Any]]:
        """Get all entities with pagination."""
        order_column = self._resolve_order_column(order_by)
        order_direction = self._resolve_order_direction(order_dir)
        query = q.BASE_SELECT_ALL.format(
            table=self._table_identifier,
            order_clause=f"{order_column} {order_direction}",
        )
        return await fetch_all(query, {"limit": limit, "offset": offset})

    @traced()
    async def count(self, where_clause: str = "", params: Mapping[str, Any] | None = None) -> int:
        """Count entities."""
        clause = f"WHERE {where_clause}" if where_clause else ""
        query = q.BASE_COUNT.format(table=self._table_identifier, where_clause=clause)
        result = await fetch_one(query, params or {})
        return result["count"] if result else 0

    @traced()
    async def delete(self, entity_id: UUID) -> bool:
        """Delete entity by ID."""
        query = q.BASE_DELETE.format(table=self._table_identifier)
        affected = await execute(query, {"id": entity_id})
        return affected > 0

    @traced()
    async def exists(self, entity_id: UUID) -> bool:
        """Check if entity exists."""
        query = q.BASE_EXISTS.format(table=self._table_identifier)
        result = await fetch_one(query, {"id": entity_id})
        return result is not None

    def _resolve_order_column(self, order_by: str | None) -> str:
        candidate = (order_by or self.default_order_by).lower()
        if candidate not in self.orderable_columns:
            return self.default_order_by.lower()
        return candidate

    @staticmethod
    def _resolve_order_direction(order_dir: str) -> str:
        direction = order_dir.upper()
        if direction not in {"ASC", "DESC"}:
            return "DESC"
        return direction
