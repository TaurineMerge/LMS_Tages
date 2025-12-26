"""Certificate service layer."""

from typing import Any
from uuid import UUID

from app.exceptions import not_found_error
from app.repositories.certificate import certificate_repository
from app.repositories.student import student_repository
from app.schemas.certificate import certificate_create, certificate_response
from app.services.storage_service import StorageError
from app.telemetry import traced


class certificate_service:
    """Service for certificate business logic."""

    def __init__(self):
        self.repository = certificate_repository
        self.student_repository = student_repository

    @traced()
    async def get_certificates(
        self, student_id: UUID | None = None, course_id: UUID | None = None
    ) -> list[certificate_response]:
        """Get certificates with optional filters."""
        certificates = await self.repository.get_filtered(student_id, course_id)
        return [certificate_response(**c) for c in certificates]

    @traced()
    async def get_certificate(self, certificate_id: UUID) -> certificate_response:
        """Get certificate by ID."""
        certificate = await self.repository.get_by_id(certificate_id)

        if not certificate:
            raise not_found_error("Certificate", str(certificate_id))

        return certificate_response(**certificate)

    @traced()
    async def create_certificate(self, data: certificate_create) -> certificate_response:
        """Create a new certificate."""
        # Verify student exists
        if not await self.student_repository.exists(data.student_id):
            raise not_found_error("Student", str(data.student_id))

        # Note: Course validation would require cross-schema query
        # In production, you'd validate course_id exists in knowledge_base.course_b

        certificate = await self.repository.create(data.model_dump())

        if not certificate:
            raise Exception("Failed to create certificate")

        return certificate_response(**certificate)

    @traced()
    async def delete_certificate(self, certificate_id: UUID) -> bool:
        """Delete certificate by ID."""
        if not await self.repository.exists(certificate_id):
            raise not_found_error("Certificate", str(certificate_id))

        return await self.repository.delete(certificate_id)

    @traced()
    async def get_certificates_by_student_grouped(self, student_id: UUID) -> dict[str, list[dict[str, Any]]]:
        """Get certificates for a student grouped by course.

        Returns a dictionary where keys are course IDs and values are lists of certificates
        for that course, including download URLs.

        Args:
            student_id: UUID of the student

        Returns:
            Dictionary with course_id as keys and list of certificate data as values
        """
        from app.services.storage_service import get_storage_service

        storage_service = get_storage_service()
        certificates = await self.repository.get_by_student_with_course_info(student_id)

        # Group certificates by course
        grouped_certificates = {}

        for cert in certificates:
            course_id = cert.get("course_id")
            if course_id is None:
                # Skip certificates without course association
                continue

            course_id_str = str(course_id)
            if course_id_str not in grouped_certificates:
                grouped_certificates[course_id_str] = []

            # Add download URL if S3 key exists
            cert_data = dict(cert)
            if cert.get("pdf_s3_key"):  # Use pdf_s3_key instead of content
                try:
                    download_url = await storage_service.get_certificate_url(cert["id"])
                    cert_data["download_url"] = download_url
                except StorageError:
                    # If URL generation fails, continue without URL
                    pass

            grouped_certificates[course_id_str].append(cert_data)

        return grouped_certificates

    @traced()
    async def download_certificate(self, certificate_id: UUID) -> bytes:
        """Download certificate PDF content.

        Args:
            certificate_id: UUID of the certificate

        Returns:
            PDF content as bytes

        Raises:
            not_found_error: If certificate not found or has no S3 key
        """
        from app.services.storage_service import get_storage_service

        # Get certificate data
        certificate = await self.repository.get_by_id(certificate_id)
        if not certificate:
            raise not_found_error("Certificate", str(certificate_id))

        # Check if certificate has S3 key
        s3_key = certificate.get("content")
        if not s3_key:
            raise not_found_error("Certificate content", f"for certificate {certificate_id}")

        # Get PDF content from S3
        storage_service = get_storage_service()
        return await storage_service.get_certificate(s3_key)


# Singleton instance
certificate_service = certificate_service()
