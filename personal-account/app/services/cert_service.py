"""Certificate generation service.

Handles certificate generation based on test attempts and student data.
Generates beautiful PDF certificates and stores them in S3.
"""

import logging
import math
from datetime import UTC, datetime
from io import BytesIO
from uuid import UUID

from reportlab.lib import colors
from reportlab.lib.pagesizes import landscape, letter
from reportlab.lib.units import inch
from reportlab.pdfgen import canvas
from reportlab.pdfgen.pathobject import PDFPathObject

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

            # Update test attempt with certificate_id
            await self.certificate_repo.update_test_attempt_certificate_id(test_attempt_id, certificate_id)

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
        - Elegant gradient background
        - Centered layout with decorative elements
        - Student name and course information
        - Score achieved with visual representation
        - Issue date and certificate number
        - Digital seal and signature

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

            # Create PDF with landscape orientation
            pdf_canvas = canvas.Canvas(buffer, pagesize=landscape(letter))
            width, height = landscape(letter)

            # Draw elegant background
            self._draw_background(pdf_canvas, width, height)

            # Draw decorative border
            self._draw_elegant_border(pdf_canvas, width, height)

            # Main title
            self._draw_title(pdf_canvas, width, height)

            # Certificate content
            self._draw_certificate_content(
                pdf_canvas, width, height, student_name, course_name, score, max_score, issue_date, student_email
            )

            # Decorative elements
            self._draw_decorative_elements(pdf_canvas, width, height)

            # Signature and seal
            self._draw_signature_and_seal(pdf_canvas, width, height)

            pdf_canvas.save()
            buffer.seek(0)
            return buffer.getvalue()

        except Exception as e:
            logger.exception("Failed to generate PDF certificate")
            raise CertificateGenerationError(f"Failed to generate PDF: {e}") from e

    @staticmethod
    def _draw_background(pdf_canvas, width: float, height: float):
        """Draw elegant gradient background.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
        """
        # Main background color
        pdf_canvas.setFillColor(colors.HexColor("#f8f9fa"))
        pdf_canvas.rect(0, 0, width, height, fill=1)

        # Gradient effect with diagonal lines
        pdf_canvas.setStrokeColor(colors.HexColor("#e9ecef"))
        pdf_canvas.setLineWidth(0.5)

        for i in range(0, int(width + height), 20):
            pdf_canvas.line(i - height, 0, i, height)

    @staticmethod
    def _draw_elegant_border(pdf_canvas, width: float, height: float):
        """Draw elegant decorative border.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
        """
        margin = 0.75 * inch

        # Outer border
        pdf_canvas.setStrokeColor(colors.HexColor("#2c3e50"))
        pdf_canvas.setLineWidth(3)
        pdf_canvas.roundRect(margin, margin, width - 2 * margin, height - 2 * margin, 20)

        # Inner decorative border
        pdf_canvas.setStrokeColor(colors.HexColor("#3498db"))
        pdf_canvas.setLineWidth(1)
        pdf_canvas.roundRect(margin + 15, margin + 15, width - 2 * (margin + 15), height - 2 * (margin + 15), 15)

        # Corner decorations
        pdf_canvas.setFillColor(colors.HexColor("#3498db"))
        corner_size = 10

        # Top-left corner
        pdf_canvas.circle(margin + corner_size, height - margin - corner_size, corner_size, fill=1)
        # Top-right corner
        pdf_canvas.circle(width - margin - corner_size, height - margin - corner_size, corner_size, fill=1)
        # Bottom-left corner
        pdf_canvas.circle(margin + corner_size, margin + corner_size, corner_size, fill=1)
        # Bottom-right corner
        pdf_canvas.circle(width - margin - corner_size, margin + corner_size, corner_size, fill=1)

    @staticmethod
    def _draw_title(pdf_canvas, width: float, height: float):
        """Draw certificate title.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
        """
        # Main title
        pdf_canvas.setFont("Helvetica-Bold", 42)
        pdf_canvas.setFillColor(colors.HexColor("#2c3e50"))
        title_text = "Certificate of Achievement"
        title_width = pdf_canvas.stringWidth(title_text, "Helvetica-Bold", 42)
        pdf_canvas.drawString((width - title_width) / 2, height - 100, title_text)

        # Subtitle
        pdf_canvas.setFont("Helvetica", 16)
        pdf_canvas.setFillColor(colors.HexColor("#7f8c8d"))
        subtitle_text = "Professional Development Program"
        subtitle_width = pdf_canvas.stringWidth(subtitle_text, "Helvetica", 16)
        pdf_canvas.drawString((width - subtitle_width) / 2, height - 130, subtitle_text)

        # Decorative line under title
        pdf_canvas.setStrokeColor(colors.HexColor("#3498db"))
        pdf_canvas.setLineWidth(2)
        line_y = height - 145
        line_length = 200
        pdf_canvas.line((width - line_length) / 2, line_y, (width + line_length) / 2, line_y)

    @staticmethod
    def _draw_certificate_content(
        pdf_canvas,
        width: float,
        height: float,
        student_name: str,
        course_name: str,
        score: int,
        max_score: int,
        issue_date: datetime,
        student_email: str,
    ):
        """Draw main certificate content.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
            student_name: Full name of the student
            course_name: Name of the course
            score: Achieved score
            max_score: Maximum possible score
            issue_date: Date when certificate was issued
            student_email: Email of the student
        """

        # Introductory text
        pdf_canvas.setFont("Helvetica", 14)
        pdf_canvas.setFillColor(colors.HexColor("#34495e"))
        intro_text = "This is to certify that"
        intro_width = pdf_canvas.stringWidth(intro_text, "Helvetica", 14)
        pdf_canvas.drawString((width - intro_width) / 2, height - 180, intro_text)

        # Student name (prominent)
        pdf_canvas.setFont("Helvetica-Bold", 32)
        pdf_canvas.setFillColor(colors.HexColor("#e74c3c"))
        name_width = pdf_canvas.stringWidth(student_name, "Helvetica-Bold", 32)
        pdf_canvas.drawString((width - name_width) / 2, height - 220, student_name)

        # Course completion text
        pdf_canvas.setFont("Helvetica", 14)
        pdf_canvas.setFillColor(colors.HexColor("#34495e"))
        completion_text = "has successfully completed the course"
        completion_width = pdf_canvas.stringWidth(completion_text, "Helvetica", 14)
        pdf_canvas.drawString((width - completion_width) / 2, height - 260, completion_text)

        # Course name
        pdf_canvas.setFont("Helvetica-Bold", 20)
        pdf_canvas.setFillColor(colors.HexColor("#27ae60"))
        course_width = pdf_canvas.stringWidth(course_name, "Helvetica-Bold", 20)
        pdf_canvas.drawString((width - course_width) / 2, height - 295, course_name)

        # Score information
        pdf_canvas.setFont("Helvetica", 12)
        pdf_canvas.setFillColor(colors.HexColor("#34495e"))
        score_percent = int((score / max_score) * 100) if max_score > 0 else 0

        score_text = f"Achieving a score of {score} out of {max_score} points ({score_percent}%)"
        score_width = pdf_canvas.stringWidth(score_text, "Helvetica", 12)
        pdf_canvas.drawString((width - score_width) / 2, height - 330, score_text)

        # Issue date
        pdf_canvas.setFont("Helvetica", 12)
        pdf_canvas.setFillColor(colors.HexColor("#7f8c8d"))
        date_text = f"Awarded on {issue_date.strftime('%B %d, %Y')}"
        date_width = pdf_canvas.stringWidth(date_text, "Helvetica", 12)
        pdf_canvas.drawString((width - date_width) / 2, height - 360, date_text)

        # Certificate ID
        pdf_canvas.setFont("Helvetica", 10)
        pdf_canvas.setFillColor(colors.HexColor("#95a5a6"))
        cert_id = f"Certificate ID: {issue_date.strftime('%Y%m%d')}-{student_email[:3].upper()}"
        cert_width = pdf_canvas.stringWidth(cert_id, "Helvetica", 10)
        pdf_canvas.drawString((width - cert_width) / 2, height - 385, cert_id)

    @staticmethod
    def _draw_decorative_elements(pdf_canvas, width: float, height: float):
        """Draw decorative elements.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
        """
        # Left side decoration
        pdf_canvas.setFillColor(colors.HexColor("#3498db"))
        pdf_canvas.setStrokeColor(colors.HexColor("#2980b9"))

        for i in range(5):
            y_pos = height / 2 + (i - 2) * 40
            pdf_canvas.circle(60, y_pos, 8, fill=1)
            pdf_canvas.circle(60, y_pos, 12, fill=0)

        # Right side decoration
        for i in range(5):
            y_pos = height / 2 + (i - 2) * 40
            pdf_canvas.circle(width - 60, y_pos, 8, fill=1)
            pdf_canvas.circle(width - 60, y_pos, 12, fill=0)

    @staticmethod
    def _draw_signature_and_seal(pdf_canvas, width: float, height: float):
        """Draw signature area and official seal.

        Args:
            pdf_canvas: reportlab canvas object
            width: Page width
            height: Page height
        """
        # Signature line
        pdf_canvas.setStrokeColor(colors.HexColor("#2c3e50"))
        pdf_canvas.setLineWidth(1)
        pdf_canvas.line(width - 250, 100, width - 100, 100)

        # Signature text
        pdf_canvas.setFont("Helvetica", 10)
        pdf_canvas.setFillColor(colors.HexColor("#34495e"))
        pdf_canvas.drawString(width - 220, 80, "Course Director")

        # Official seal (circular)
        seal_center_x = 120
        seal_center_y = 120
        seal_radius = 40

        # Seal background
        pdf_canvas.setFillColor(colors.HexColor("#e74c3c"))
        pdf_canvas.circle(seal_center_x, seal_center_y, seal_radius, fill=1)

        # Seal border
        pdf_canvas.setStrokeColor(colors.HexColor("#c0392b"))
        pdf_canvas.setLineWidth(3)
        pdf_canvas.circle(seal_center_x, seal_center_y, seal_radius, fill=0)

        # Seal text
        pdf_canvas.setFillColor(colors.white)
        pdf_canvas.setFont("Helvetica-Bold", 8)
        pdf_canvas.drawString(seal_center_x - 15, seal_center_y + 5, "OFFICIAL")
        pdf_canvas.drawString(seal_center_x - 12, seal_center_y - 5, "SEAL")

        # Seal star
        pdf_canvas.setFillColor(colors.HexColor("#f39c12"))
        star_points = []
        for i in range(10):
            angle = i * 36  # 36 degrees per point
            radius = 15 if i % 2 == 0 else 8
            x = seal_center_x + radius * math.cos(math.radians(angle))
            y = seal_center_y + radius * math.sin(math.radians(angle))
            star_points.extend([x, y])

        if len(star_points) >= 6:
            # Draw star polygon using PDF path
            path = PDFPathObject()
            path.moveTo(star_points[0], star_points[1])
            for i in range(2, len(star_points), 2):
                path.lineTo(star_points[i], star_points[i + 1])
            path.close()
            pdf_canvas.drawPath(path, fill=1, stroke=0)


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
