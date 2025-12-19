"""Student repository."""

import json
from typing import Any
from uuid import UUID

from app.database import execute_returning, fetch_all, fetch_one
from app.db import queries as q
from app.repositories.base import base_repository
from app.telemetry import traced


class student_repository(base_repository):
    """Repository for student operations."""

    def __init__(self):
        super().__init__(
            "student_s",
            default_order_by="created_at",
            orderable_columns={"created_at", "updated_at", "name", "surname", "email"},
        )

    _mutable_fields = {"name", "surname", "birth_date", "avatar", "contacts", "email", "phone"}

    @traced()
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new student."""
        contacts_json = json.dumps(data.get("contacts") or {})

        params = {
            "name": data["name"],
            "surname": data["surname"],
            "birth_date": data.get("birth_date"),
            "avatar": data.get("avatar"),
            "contacts": contacts_json,
            "email": data["email"],
            "phone": data.get("phone"),
        }
        return await execute_returning(q.STUDENT_INSERT, params)

    @traced()
    async def update(self, student_id: UUID, data: dict[str, Any], conn=None) -> dict[str, Any] | None:
        """Update student by ID."""
        # Filter out None values to only update provided fields
        update_data = {k: v for k, v in data.items() if v is not None}

        if not update_data:
            return await self.get_by_id(student_id)

        # Build dynamic UPDATE query
        set_clauses: list[str] = []
        params = {"id": student_id}

        for key, value in update_data.items():
            if key not in self._mutable_fields:
                continue
            if key == "contacts":
                params[key] = json.dumps(value)
                set_clauses.append(f"{key} = CAST(:{key} AS jsonb)")
            else:
                params[key] = value
                set_clauses.append(f"{key} = :{key}")

        if not set_clauses:
            return await self.get_by_id(student_id)

        set_clauses.append("updated_at = CURRENT_TIMESTAMP")
        query = q.STUDENT_UPDATE_TEMPLATE.format(set_clause=", ".join(set_clauses))
        return await execute_returning(query, params)

    @traced()
    async def get_by_email(self, email: str) -> dict[str, Any] | None:
        """Get student by email."""
        return await fetch_one(q.STUDENT_BY_EMAIL, {"email": email})

    @traced()
    async def get_paginated(self, page: int = 1, limit: int = 20) -> tuple[list[dict[str, Any]], int]:
        """Get paginated list of students."""
        offset = (page - 1) * limit

        # Get total count
        count_result = await fetch_one(q.STUDENT_COUNT, {})
        total = count_result["count"] if count_result else 0

        # Get paginated data
        students = await fetch_all(q.STUDENT_PAGINATED, {"limit": limit, "offset": offset})

        return students, total

    @traced()
    async def email_exists(self, email: str, exclude_id: UUID | None = None, conn=None) -> bool:
        """Check if email already exists."""
        if exclude_id:
            exclude_clause = "AND id != :exclude_id"
            params = {"email": email, "exclude_id": exclude_id}
        else:
            exclude_clause = ""
            params = {"email": email}
        query = q.STUDENT_EMAIL_EXISTS.format(exclude_clause=exclude_clause)
        result = await fetch_one(query, params)
        return result is not None


# Singleton instance
student_repository = student_repository()
