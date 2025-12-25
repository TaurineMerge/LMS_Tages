"""Certificate repository."""

from typing import Any
from uuid import UUID

from app.database import execute, execute_returning, fetch_all, fetch_one, fetch_one_value
from app.db import queries as q
from app.repositories.base import base_repository
from app.telemetry import traced


class certificate_repository(base_repository):
    """Repository for certificate operations."""

    def __init__(self):
        super().__init__("certificate_b")

    async def _get_next_certificate_number(self) -> int:
        """Generate next certificate number."""
        result = await fetch_one_value(q.GET_MAX_CERTIFICATE_NUMBER)
        return (result or 0) + 1

    @traced("creation of certificate", record_args=True, record_result=True)
    async def create(self, data: dict[str, Any]) -> dict[str, Any] | None:
        """Create a new certificate."""
        certificate_number = await self._get_next_certificate_number()
        params = {
            "student_id": data["student_id"],
            "certificate_number": certificate_number,
            "pdf_s3_key": data.get("pdf_s3_key"),
            "snapshot_s3_key": data.get("snapshot_s3_key"),
        }
        return await execute_returning(q.CERTIFICATE_INSERT, params)

    @traced("certificate.get_passing_attempts_without_certificates", record_args=True, record_result=True)
    async def get_passing_attempts_without_certificates(self) -> list[dict]:
        """Get passing test attempts that don't have certificates yet.

        Returns attempts where score >= passing_score and no certificate exists
        for the test_attempt_id.

        Returns:
            List of attempt dictionaries with student_id, course_id, id, score, max_score, course_name
        """

    @traced("certificate.get_passing_attempts_without_certificates_for_student", record_args=True, record_result=True)
    async def get_passing_attempts_without_certificates_for_student(self, student_id: UUID) -> list[dict]:
        """Get passing test attempts that don't have certificates yet for a specific student.

        Returns attempts where score >= passing_score and no certificate exists
        for the test_attempt_id.

        Args:
            student_id: UUID of the student

        Returns:
            List of attempt dictionaries with student_id, course_id, id, score, max_score, course_name
        """
        return await fetch_all(q.GET_PASSING_ATTEMPTS_WITHOUT_CERTIFICATES_FOR_STUDENT, {"student_id": student_id})

    async def get_by_student(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a student."""
        return await fetch_all(q.CERTIFICATES_BY_STUDENT, {"student_id": student_id})

    @traced()
    async def get_by_student_with_course_info(self, student_id: UUID) -> list[dict[str, Any]]:
        """Get all certificates for a student with course information."""
        return await fetch_all(q.CERTIFICATES_BY_STUDENT_WITH_COURSE, {"student_id": student_id})

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

    @traced()
    async def update_s3_key(self, certificate_id: UUID, s3_key: str) -> dict[str, Any] | None:
        """Update S3 key for certificate."""
        # Изменено: обновлять pdf_s3_key вместо content
        params = {
            "id": certificate_id,
            "pdf_s3_key": s3_key,  # Или snapshot_s3_key, в зависимости от логики
        }
        return await execute_returning(q.CERTIFICATE_UPDATE_S3_KEY, params)

    @traced()
    async def update_test_attempt_certificate_id(self, test_attempt_id: UUID, certificate_id: UUID) -> bool:
        """Update certificate_id in test attempt."""
        params = {
            "test_attempt_id": test_attempt_id,
            "certificate_id": certificate_id,
        }
        result = await execute(q.UPDATE_TEST_ATTEMPT_CERTIFICATE_ID, params)
        return result is not None


# Singleton instance
certificate_repository = certificate_repository()
