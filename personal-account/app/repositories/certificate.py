"""Certificate repository."""
from typing import Any
from uuid import UUID

from app.database import fetch_one, fetch_all, execute_returning
from app.repositories.base import BaseRepository


class CertificateRepository(BaseRepository):
    """Repository for certificate operations."""
    
    def __init__(self):
        super().__init__("certificate_b")
    
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new certificate."""
        query = """
            INSERT INTO personal_account.certificate_b 
            (content, student_id, course_id, test_attempt_id)
            VALUES (%s, %s, %s, %s)
            RETURNING *
        """
        params = (
            data.get("content"),
            data["student_id"],
            data["course_id"],
            data.get("test_attempt_id"),
        )
        return await execute_returning(query, params)
    
    async def get_by_student(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a student."""
        query = """
            SELECT * FROM personal_account.certificate_b
            WHERE student_id = %s
            ORDER BY created_at DESC
        """
        return await fetch_all(query, (student_id,))
    
    async def get_by_course(self, course_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a course."""
        query = """
            SELECT * FROM personal_account.certificate_b
            WHERE course_id = %s
            ORDER BY created_at DESC
        """
        return await fetch_all(query, (course_id,))
    
    async def get_filtered(
        self,
        student_id: UUID | None = None,
        course_id: UUID | None = None
    ) -> list[dict[str, Any]]:
        """Get certificates with optional filters."""
        conditions = []
        params = []
        
        if student_id:
            conditions.append("student_id = %s")
            params.append(student_id)
        
        if course_id:
            conditions.append("course_id = %s")
            params.append(course_id)
        
        where_clause = ""
        if conditions:
            where_clause = "WHERE " + " AND ".join(conditions)
        
        query = f"""
            SELECT * FROM personal_account.certificate_b
            {where_clause}
            ORDER BY created_at DESC
        """
        return await fetch_all(query, tuple(params) if params else None)
    
    async def get_by_number(self, certificate_number: int) -> dict[str, Any] | None:
        """Get certificate by its unique number."""
        query = """
            SELECT * FROM personal_account.certificate_b
            WHERE certificate_number = %s
        """
        return await fetch_one(query, (certificate_number,))


# Singleton instance
certificate_repository = CertificateRepository()
