"""Certificate repository."""

from typing import Any
from uuid import UUID

from app.database import execute_returning, fetch_all, fetch_one
from app.db import queries as q
from app.repositories.base import base_repository
from app.telemetry import traced


class certificate_repository(base_repository):
    """Repository for certificate operations."""

    def __init__(self):
        super().__init__("certificate_b")

    @traced()
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new certificate."""
        params = {
            "content": data.get("content"),
            "student_id": data["student_id"],
            "course_id": data["course_id"],
            "test_attempt_id": data.get("test_attempt_id"),
        }
        return await execute_returning(q.CERTIFICATE_INSERT, params)

    @traced()
    async def get_by_student(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a student."""
        return await fetch_all(q.CERTIFICATES_BY_STUDENT, {"student_id": student_id})

    @traced()
    async def get_by_course(self, course_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a course."""
        return await fetch_all(q.CERTIFICATES_BY_COURSE, {"course_id": course_id})

    @traced()
    async def get_filtered(self, student_id: UUID | None = None, course_id: UUID | None = None) -> list[dict[str, Any]]:
        """Get certificates with optional filters."""
        conditions = []
        params: dict[str, Any] = {}

        if student_id:
            conditions.append("student_id = :student_id")
            params["student_id"] = student_id

        if course_id:
            conditions.append("course_id = :course_id")
            params["course_id"] = course_id

        where_clause = "WHERE " + " AND ".join(conditions) if conditions else ""
        query = q.CERTIFICATES_FILTERED_TEMPLATE.format(where_clause=where_clause)
        return await fetch_all(query, params or {})

    @traced()
    async def get_by_number(self, certificate_number: int) -> dict[str, Any] | None:
        """Get certificate by its unique number."""
        return await fetch_one(q.CERTIFICATE_BY_NUMBER, {"certificate_number": certificate_number})


# Singleton instance
certificate_repository = certificate_repository()
