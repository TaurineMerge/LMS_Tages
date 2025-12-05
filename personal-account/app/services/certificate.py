"""Certificate service layer."""
from uuid import UUID

from app.repositories.certificate import certificate_repository
from app.repositories.student import student_repository
from app.schemas.certificate import CertificateCreate, CertificateResponse
from app.exceptions import NotFoundError


class CertificateService:
    """Service for certificate business logic."""
    
    def __init__(self):
        self.repository = certificate_repository
        self.student_repository = student_repository
    
    async def get_certificates(
        self,
        student_id: UUID | None = None,
        course_id: UUID | None = None
    ) -> list[CertificateResponse]:
        """Get certificates with optional filters."""
        certificates = await self.repository.get_filtered(student_id, course_id)
        return [CertificateResponse(**c) for c in certificates]
    
    async def get_certificate(self, certificate_id: UUID) -> CertificateResponse:
        """Get certificate by ID."""
        certificate = await self.repository.get_by_id(certificate_id)
        
        if not certificate:
            raise NotFoundError("Certificate", str(certificate_id))
        
        return CertificateResponse(**certificate)
    
    async def create_certificate(self, data: CertificateCreate) -> CertificateResponse:
        """Create a new certificate."""
        # Verify student exists
        if not await self.student_repository.exists(data.student_id):
            raise NotFoundError("Student", str(data.student_id))
        
        # Note: Course validation would require cross-schema query
        # In production, you'd validate course_id exists in knowledge_base.course_b
        
        certificate = await self.repository.create(data.model_dump())
        
        if not certificate:
            raise Exception("Failed to create certificate")
        
        return CertificateResponse(**certificate)
    
    async def delete_certificate(self, certificate_id: UUID) -> bool:
        """Delete certificate by ID."""
        if not await self.repository.exists(certificate_id):
            raise NotFoundError("Certificate", str(certificate_id))
        
        return await self.repository.delete(certificate_id)


# Singleton instance
certificate_service = CertificateService()
