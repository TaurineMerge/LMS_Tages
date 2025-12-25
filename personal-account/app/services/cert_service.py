"""Certificate generation service.

Handles certificate generation based on test attempts and student data.
Generates beautiful PDF certificates and stores them in S3.
"""

import logging
from datetime import UTC, datetime
from io import BytesIO
from uuid import UUID

from reportlab.lib import colors
from reportlab.lib.pagesizes import landscape, letter
from reportlab.lib.units import inch
from reportlab.pdfgen import canvas

from app.repositories.certificate import certificate_repository
from app.schemas.certificate import certificate_create
from app.services.keycloak import KeycloakService
from app.services.storage_service import StorageError, get_storage_service
from app.telemetry import traced

logger = logging.getLogger(__name__)


class CertificateGenerationError(Exception):
    """Exception raised during certificate generation."""

    pass


class CertificateService:
    """Service for generating certificates based on test attempts.

    Generates beautiful PDF certificates with:
    - Student full name (from Keycloak)
    - Course name
    - Score achieved
    - Issue date
    - Certificate number
    - Digital signature/seal

    Stores generated certificates in S3 and updates database with S3 key.
    """

    def __init__(
        self,
        storage_service=None,
        certificate_repository_inst=None,
        keycloak_service_inst=None,
    ):
        """Initialize certificate service.

        Args:
            storage_service: S3 storage service instance
            certificate_repository_inst: Certificate repository instance
            keycloak_service_inst: Keycloak service instance
        """
        self.storage_service = storage_service or get_storage_service()
        self.certificate_repo = certificate_repository_inst or certificate_repository
        self.keycloak_service = keycloak_service_inst or KeycloakService()

    @traced("certificate_service.generate_certificate", record_args=True, record_result=True)
    async def generate_certificate(
        self,
        student_id: UUID,
        course_id: UUID,
        course_name: str,
        test_attempt_id: UUID,
        score: int,
        max_score: int,
    ) -> tuple[UUID, str]:
        """Generate and store certificate for successful test attempt.

        Fetches student data from Keycloak, generates beautiful PDF certificate,
        stores it in S3, and updates database with S3 key.

        Args:
            student_id: UUID of the student
            course_id: UUID of the course
            course_name: Human-readable course name
            test_attempt_id: UUID of the test attempt
            score: Achieved score
            max_score: Maximum possible score

        Returns:
            Tuple of (certificate_id, s3_key)

        Raises:
            CertificateGenerationError: If certificate generation or storage fails
        """
        try:
            # Fetch student data from Keycloak
            student_data = await self._fetch_student_data(student_id)

            # Generate PDF certificate
            pdf_content = await self._generate_pdf(
                student_name=f"{student_data['name']} {student_data['surname']}",
                student_email=student_data["email"],
                course_name=course_name,
                score=score,
                max_score=max_score,
                issue_date=datetime.now(tz=UTC),
            )

            # Create certificate record in database
            cert_data = certificate_create(
                student_id=student_id,
                course_id=course_id,
                test_attempt_id=test_attempt_id,
                content=None,  # We're storing in S3, not in DB
            )
            certificate = await self.certificate_repo.create(cert_data.model_dump())

            if not certificate:
                raise CertificateGenerationError("Failed to create certificate record in database")

            certificate_id = certificate.get("id")

            # Store PDF in S3
            try:
                s3_key = await self.storage_service.store_certificate(
                    certificate_id=certificate_id,
                    student_id=student_id,
                    course_id=course_id,
                    pdf_content=pdf_content,
                )
            except StorageError as e:
                # Rollback database record if S3 storage fails
                await self.certificate_repo.delete(certificate_id)
                raise CertificateGenerationError(f"Failed to store certificate in S3: {e}") from e

            # Update certificate record with S3 key
            await self.certificate_repo.update_s3_key(certificate_id, s3_key)

            logger.info(
                "Certificate generated successfully: certificate_id=%s, s3_key=%s",
                certificate_id,
                s3_key,
            )

            return certificate_id, s3_key

        except CertificateGenerationError:
            raise
        except Exception as e:
            logger.exception("Unexpected error during certificate generation")
            raise CertificateGenerationError(f"Unexpected error: {e}") from e

    @traced("certificate_service._fetch_student_data", record_args=True, record_result=True)
    async def _fetch_student_data(self, student_id: UUID) -> dict:
        """Fetch student data from Keycloak.

        Args:
            student_id: UUID of the student (should match Keycloak user ID)

        Returns:
            Dict with student data (name, surname, email)

        Raises:
            CertificateGenerationError: If student data cannot be fetched
        """
        try:
            # Fetch user data from Keycloak using service
            user_info = self.keycloak_service.get_user_data(str(student_id))

            if not user_info:
                raise CertificateGenerationError(f"Student {student_id} not found in Keycloak")

            return {
                "name": user_info.get("name", ""),
                "surname": user_info.get("surname", ""),
                "email": user_info.get("email", ""),
            }
        except Exception as e:
            logger.exception("Failed to fetch student data from Keycloak")
            raise CertificateGenerationError(f"Failed to fetch student data: {e}") from e

    async def _generate_pdf(
        self,
        student_name: str,
        student_email: str,
        course_name: str,
        score: int,
        max_score: int,
        issue_date: datetime,
    ) -> bytes:
        """Generate beautiful PDF certificate.

        Creates a professional-looking certificate with:
        - Centered title and decorative elements
        - Student name and course information
        - Score achieved
        - Issue date
        - Certificate number (sequential)

        Args:
            student_name: Full name of the student
            student_email: Email of the student
            course_name: Name of the course
            score: Achieved score
            max_score: Maximum possible score
            issue_date: Date when certificate was issued

        Returns:
            PDF content as bytes

        Raises:
            CertificateGenerationError: If PDF generation fails
        """
        try:
            buffer = BytesIO()

            # Create PDF with landscape orientation for better certificate layout
            pdf_canvas = canvas.Canvas(buffer, pagesize=landscape(letter))
            width, height = landscape(letter)

            # Set up fonts
            pdf_canvas.setFont("Helvetica-Bold", 36)

            # Background - optional decorative border
            self._draw_border(pdf_canvas, width, height)

            # Title
            pdf_canvas.drawString(
                width / 2 - 100,
                height - 80,
                "Certificate of Completion",
            )

            pdf_canvas.setFont("Helvetica", 12)
            pdf_canvas.drawString(20, height - 120, "This certifies that")

            # Student name (highlighted)
            pdf_canvas.setFont("Helvetica-Bold", 20)
            pdf_canvas.setFillColor(colors.HexColor("#1F4788"))
            pdf_canvas.drawString(20, height - 160, student_name)

            # Course and score information
            pdf_canvas.setFont("Helvetica", 11)
            pdf_canvas.setFillColor(colors.black)

            y_position = height - 200
            pdf_canvas.drawString(
                20,
                y_position,
                "has successfully completed the course:",
            )

            pdf_canvas.setFont("Helvetica-Bold", 14)
            pdf_canvas.setFillColor(colors.HexColor("#1F4788"))
            pdf_canvas.drawString(20, y_position - 30, course_name)

            # Score information
            pdf_canvas.setFont("Helvetica", 11)
            pdf_canvas.setFillColor(colors.black)
            score_percent = int((score / max_score) * 100) if max_score > 0 else 0
            pdf_canvas.drawString(
                20,
                y_position - 70,
                f"with a score of {score} out of {max_score} points ({score_percent}%)",
            )

            # Issue date
            pdf_canvas.setFont("Helvetica", 10)
            pdf_canvas.setFillColor(colors.grey)
            issue_date_str = issue_date.strftime("%B %d, %Y")
            pdf_canvas.drawString(
                20,
                y_position - 110,
                f"Issued on: {issue_date_str}",
            )

            # Signature area (decorative)
            pdf_canvas.setFont("Helvetica", 10)
            pdf_canvas.setFillColor(colors.black)
            pdf_canvas.drawString(width - 200, 80, "________________________")
            pdf_canvas.drawString(width - 200, 60, "Course Administrator")

            # Footer
            pdf_canvas.setFont("Helvetica", 8)
            pdf_canvas.setFillColor(colors.grey)
            pdf_canvas.drawString(
                20,
                30,
                f"Certificate ID: {issue_date.strftime('%Y%m%d')}-{student_email[:3].upper()}",
            )

            pdf_canvas.save()
            buffer.seek(0)
            return buffer.getvalue()

        except Exception as e:
            logger.exception("Failed to generate PDF certificate")
            raise CertificateGenerationError(f"Failed to generate PDF: {e}") from e

    @staticmethod
    def _draw_border(pdf_canvas, width: float, height: float, margin: float = 0.5 * inch):
        """Draw decorative border around certificate.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
            margin: Margin from edge
        """
        pdf_canvas.setStrokeColor(colors.HexColor("#1F4788"))
        pdf_canvas.setLineWidth(2)
        pdf_canvas.rect(margin, margin, width - 2 * margin, height - 2 * margin)

        # Optional inner border
        pdf_canvas.setLineWidth(1)
        pdf_canvas.rect(
            margin + 10,
            margin + 10,
            width - 2 * (margin + 10),
            height - 2 * (margin + 10),
        )


# Global instance
_certificate_service_instance: CertificateService | None = None


def get_certificate_service() -> CertificateService:
    """Get or create certificate service instance.

    Returns:
        CertificateService instance
    """
    global _certificate_service_instance
    if _certificate_service_instance is None:
        _certificate_service_instance = CertificateService()
    return _certificate_service_instance
