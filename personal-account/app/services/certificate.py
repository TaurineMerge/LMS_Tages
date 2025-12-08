"""Certificate service layer."""
from uuid import UUID

from app.repositories.certificate import certificate_repository
from app.repositories.student import student_repository
from app.schemas.certificate import certificate_create, certificate_response
from app.exceptions import not_found_error
from app.telemetry import traced


class certificate_service:
    """Service for certificate business logic."""
    
    def __init__(self):
        self.repository = certificate_repository
        self.student_repository = student_repository
    
    @traced()
    async def get_certificates(
        self,
        student_id: UUID | None = None,
        course_id: UUID | None = None
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


# Singleton instance
certificate_service = certificate_service()
