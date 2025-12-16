"""Visit repository."""

from typing import Any
from uuid import UUID

from app.database import execute_returning, fetch_all, fetch_one
from app.db.queries import (
    VISIT_BY_ID,
    VISIT_BY_LESSON,
    VISIT_BY_STUDENT,
    VISIT_EXISTS,
    VISIT_FILTERED_TEMPLATE,
    VISIT_INSERT,
)
from app.repositories.base import base_repository
from app.telemetry import traced


class visit_repository(base_repository):
    """Repository for visit operations."""

    def __init__(self):
        super().__init__(
            "visit_students_for_lessons",
            default_order_by="id",
            orderable_columns={"id"},
        )

    @traced()
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new visit record."""
        params = {
            "student_id": data["student_id"],
            "lesson_id": data["lesson_id"],
        }
        return await execute_returning(VISIT_INSERT, params)

    @traced()
    async def get_by_student(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all visits for a student."""
        return await fetch_all(VISIT_BY_STUDENT, {"student_id": student_id})

    @traced()
    async def get_by_lesson(self, lesson_id: UUID) -> list[dict[str, Any]]:
        """Get all visits for a lesson."""
        return await fetch_all(VISIT_BY_LESSON, {"lesson_id": lesson_id})

    @traced()
    async def get_filtered(self, student_id: UUID | None = None, lesson_id: UUID | None = None) -> list[dict[str, Any]]:
        """Get visits with optional filters."""
        conditions: list[str] = []
        params: dict[str, Any] = {}

        if student_id:
            conditions.append("student_id = :student_id")
            params["student_id"] = student_id

        if lesson_id:
            conditions.append("lesson_id = :lesson_id")
            params["lesson_id"] = lesson_id

        where_clause = ""
        if conditions:
            where_clause = "WHERE " + " AND ".join(conditions)

        query = VISIT_FILTERED_TEMPLATE.format(where_clause=where_clause)
        return await fetch_all(query, params)

    @traced()
    async def visit_exists(self, student_id: UUID, lesson_id: UUID) -> bool:
        """Check if a visit record already exists."""
        params = {"student_id": student_id, "lesson_id": lesson_id}
        result = await fetch_one(VISIT_EXISTS, params)
        return result is not None

    @traced()
    async def get_by_id(self, entity_id: UUID) -> dict[str, Any] | None:
        """Get visit by ID. Override because this table doesn't have created_at."""
        return await fetch_one(VISIT_BY_ID, {"id": entity_id})


# Singleton instance
visit_repository = visit_repository()
