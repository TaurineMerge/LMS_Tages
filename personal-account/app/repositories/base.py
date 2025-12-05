"""Base repository with common CRUD operations."""
from typing import Any
from uuid import UUID

from app.database import fetch_one, fetch_all, execute, execute_returning


class BaseRepository:
    """Base repository class with common database operations."""
    
    def __init__(self, table_name: str, schema: str = "personal_account"):
        self.table_name = table_name
        self.schema = schema
        self.full_table_name = f"{schema}.{table_name}"
    
    async def get_by_id(self, entity_id: UUID) -> dict[str, Any] | None:
        """Get entity by ID."""
        query = f"SELECT * FROM {self.full_table_name} WHERE id = %s"
        return await fetch_one(query, (entity_id,))
    
    async def get_all(
        self, 
        limit: int = 100, 
        offset: int = 0,
        order_by: str = "created_at",
        order_dir: str = "DESC"
    ) -> list[dict[str, Any]]:
        """Get all entities with pagination."""
        query = f"""
            SELECT * FROM {self.full_table_name}
            ORDER BY {order_by} {order_dir}
            LIMIT %s OFFSET %s
        """
        return await fetch_all(query, (limit, offset))
    
    async def count(self, where_clause: str = "", params: tuple = ()) -> int:
        """Count entities."""
        query = f"SELECT COUNT(*) as count FROM {self.full_table_name}"
        if where_clause:
            query += f" WHERE {where_clause}"
        result = await fetch_one(query, params)
        return result["count"] if result else 0
    
    async def delete(self, entity_id: UUID) -> bool:
        """Delete entity by ID."""
        query = f"DELETE FROM {self.full_table_name} WHERE id = %s"
        affected = await execute(query, (entity_id,))
        return affected > 0
    
    async def exists(self, entity_id: UUID) -> bool:
        """Check if entity exists."""
        query = f"SELECT 1 FROM {self.full_table_name} WHERE id = %s LIMIT 1"
        result = await fetch_one(query, (entity_id,))
        return result is not None
