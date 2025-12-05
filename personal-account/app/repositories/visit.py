"""Visit repository."""
from typing import Any
from uuid import UUID

from app.database import fetch_one, fetch_all, execute_returning
from app.repositories.base import BaseRepository


class VisitRepository(BaseRepository):
    """Repository for visit operations."""
    
    def __init__(self):
        super().__init__("visit_students_for_lessons")
    
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new visit record."""
        query = """
            INSERT INTO personal_account.visit_students_for_lessons 
            (student_id, lesson_id)
            VALUES (%s, %s)
            ON CONFLICT (student_id, lesson_id) DO NOTHING
            RETURNING *
        """
        params = (data["student_id"], data["lesson_id"])
        return await execute_returning(query, params)
    
    async def get_by_student(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all visits for a student."""
        query = """
            SELECT * FROM personal_account.visit_students_for_lessons
            WHERE student_id = %s
        """
        return await fetch_all(query, (student_id,))
    
    async def get_by_lesson(self, lesson_id: UUID) -> list[dict[str, Any]]:
        """Get all visits for a lesson."""
        query = """
            SELECT * FROM personal_account.visit_students_for_lessons
            WHERE lesson_id = %s
        """
        return await fetch_all(query, (lesson_id,))
    
    async def get_filtered(
        self,
        student_id: UUID | None = None,
        lesson_id: UUID | None = None
    ) -> list[dict[str, Any]]:
        """Get visits with optional filters."""
        conditions = []
        params = []
        
        if student_id:
            conditions.append("student_id = %s")
            params.append(student_id)
        
        if lesson_id:
            conditions.append("lesson_id = %s")
            params.append(lesson_id)
        
        where_clause = ""
        if conditions:
            where_clause = "WHERE " + " AND ".join(conditions)
        
        query = f"""
            SELECT * FROM personal_account.visit_students_for_lessons
            {where_clause}
        """
        return await fetch_all(query, tuple(params) if params else None)
    
    async def visit_exists(self, student_id: UUID, lesson_id: UUID) -> bool:
        """Check if a visit record already exists."""
        query = """
            SELECT 1 FROM personal_account.visit_students_for_lessons
            WHERE student_id = %s AND lesson_id = %s
            LIMIT 1
        """
        result = await fetch_one(query, (student_id, lesson_id))
        return result is not None
    
    async def get_by_id(self, entity_id: UUID) -> dict[str, Any] | None:
        """Get visit by ID. Override because this table doesn't have created_at."""
        query = "SELECT * FROM personal_account.visit_students_for_lessons WHERE id = %s"
        return await fetch_one(query, (entity_id,))


# Singleton instance
visit_repository = VisitRepository()
