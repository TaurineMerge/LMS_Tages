"""Student repository."""
import json
from typing import Any
from uuid import UUID

from app.database import fetch_one, fetch_all, execute_returning
from app.repositories.base import BaseRepository


class StudentRepository(BaseRepository):
    """Repository for student operations."""
    
    def __init__(self):
        super().__init__("student_s")
    
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new student."""
        contacts_json = json.dumps(data.get("contacts") or {})
        
        query = """
            INSERT INTO personal_account.student_s 
            (name, surname, birth_date, avatar, contacts, email, phone)
            VALUES (%s, %s, %s, %s, %s::jsonb, %s, %s)
            RETURNING *
        """
        params = (
            data["name"],
            data["surname"],
            data.get("birth_date"),
            data.get("avatar"),
            contacts_json,
            data["email"],
            data.get("phone"),
        )
        return await execute_returning(query, params)
    
    async def update(self, student_id: UUID, data: dict[str, Any]) -> dict[str, Any] | None:
        """Update student by ID."""
        # Filter out None values to only update provided fields
        update_data = {k: v for k, v in data.items() if v is not None}
        
        if not update_data:
            return await self.get_by_id(student_id)
        
        # Build dynamic UPDATE query
        set_clauses = []
        params = []
        
        for key, value in update_data.items():
            if key == "contacts":
                set_clauses.append(f"{key} = %s::jsonb")
                params.append(json.dumps(value))
            else:
                set_clauses.append(f"{key} = %s")
                params.append(value)
        
        set_clauses.append("updated_at = CURRENT_TIMESTAMP")
        params.append(student_id)
        
        query = f"""
            UPDATE personal_account.student_s
            SET {', '.join(set_clauses)}
            WHERE id = %s
            RETURNING *
        """
        return await execute_returning(query, tuple(params))
    
    async def get_by_email(self, email: str) -> dict[str, Any] | None:
        """Get student by email."""
        query = "SELECT * FROM personal_account.student_s WHERE email = %s"
        return await fetch_one(query, (email,))
    
    async def get_paginated(
        self,
        page: int = 1,
        limit: int = 20
    ) -> tuple[list[dict[str, Any]], int]:
        """Get paginated list of students."""
        offset = (page - 1) * limit
        
        # Get total count
        count_query = "SELECT COUNT(*) as count FROM personal_account.student_s"
        count_result = await fetch_one(count_query)
        total = count_result["count"] if count_result else 0
        
        # Get paginated data
        query = """
            SELECT * FROM personal_account.student_s
            ORDER BY created_at DESC
            LIMIT %s OFFSET %s
        """
        students = await fetch_all(query, (limit, offset))
        
        return students, total
    
    async def email_exists(self, email: str, exclude_id: UUID | None = None) -> bool:
        """Check if email already exists."""
        if exclude_id:
            query = """
                SELECT 1 FROM personal_account.student_s 
                WHERE email = %s AND id != %s 
                LIMIT 1
            """
            result = await fetch_one(query, (email, exclude_id))
        else:
            query = "SELECT 1 FROM personal_account.student_s WHERE email = %s LIMIT 1"
            result = await fetch_one(query, (email,))
        return result is not None


# Singleton instance
student_repository = StudentRepository()
